//
// command/archive.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

// Archive makes a tar.gz file consists of files maintained a git repository.
func Archive(dir string, filename string) (err error) {

	writeFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer writeFile.Close()

	writer := bufio.NewWriter(writeFile)
	defer writer.Flush()

	zipWriter, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
	if err != nil {
		return
	}
	defer zipWriter.Close()

	tarWriter := tar.NewWriter(zipWriter)
	defer tarWriter.Close()

	// Change dir and run.
	cd, err := os.Getwd()
	if err != nil {
		return
	}
	if err = os.Chdir(dir); err != nil {
		return
	}
	defer os.Chdir(cd)

	// Listing up and write to a tarball.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan string)
	doneLGS := make(chan error)
	doneLGR := make(chan error)
	doneTB := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		doneLGS <- listupGitSources(ctx, ch)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		doneLGR <- listupGitRepository(ctx, ch)
	}()
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Require to close channel ch to end this goroutine.
	go tarballing(tarWriter, ch, doneTB)
	err = <-doneTB
	if err != nil {
		cancel()
		return
	}

	err1 := <-doneLGS
	err2 := <-doneLGR
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}

	return
}

// tarballing is a go-routine which write a file given via ch to a tar writer.
// It puts nil to done when it ends. If an error occurs, it puts the error to
// done.
func tarballing(writer *tar.Writer, ch <-chan string, done chan<- error) {

	var info os.FileInfo
	var header *tar.Header
	var err error

	for path := range ch {

		// For Windows: Replace path delimiters.
		path = filepath.ToSlash(path)
		fmt.Println(path)

		// Write a file header.
		info, err = os.Stat(path)
		if err != nil {
			fmt.Printf("Cannot find %s (%s)", path, err.Error())
			break
		}

		header, err = tar.FileInfoHeader(info, path)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		if strings.HasPrefix(path, "../") {
			header.Name = path[3:]
		} else {
			header.Name = path
		}
		writer.WriteHeader(header)

		// Write the body.
		if err = copyFile(path, writer); err != nil {
			break
		}
	}
	done <- err

}

// listupGitRepository lists up git repository and puts founded paths to a given ch.
func listupGitRepository(ctx context.Context, ch chan<- string) error {

	return filepath.Walk(".git", func(path string, info os.FileInfo, err error) error {

		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ch <- path
			return nil
		}

	})

}

// listupGitSources lists up git sources and puts finding paths to a given ch.
func listupGitSources(ctx context.Context, ch chan<- string) (err error) {

	cmd := exec.Command("git", "ls-files")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer stdout.Close()

	doneRL := make(chan error)
	defer close(doneRL)

	go readLine(ctx, stdout, ch, doneRL)
	if err = cmd.Run(); err != nil {
		<-doneRL
		return
	}

	return <-doneRL

}

// readLine is a go-routine which reads a line from rd and puts it to ch.
// It sends nil to done when reads all line or some error when an error
// occurs.
func readLine(ctx context.Context, rd io.Reader, ch chan<- string, done chan<- error) {

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			done <- ctx.Err()
			return
		default:
			ch <- scanner.Text()
		}
	}
	// done <- scanner.Err()
	done <- nil

}

// copyFile opens a given file and put its body to a given writer.
func copyFile(path string, writer io.Writer) (err error) {

	// Prepare to write a file body.
	fp, err := os.Open(path)
	if err != nil {
		return
	}
	defer fp.Close()

	_, err = io.Copy(writer, fp)
	return

}
