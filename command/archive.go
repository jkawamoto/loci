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
		return err
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

	// Litsing up and write to a tarball.
	ch := make(chan string)
	doneLGS := make(chan error)
	doneTB := make(chan error)

	go listupGitSources(ch, doneLGS)
	go tarballing(tarWriter, ch, doneTB)

	err = <-doneLGS
	close(ch)
	if err != nil {
		return
	}
	return <-doneTB

}

// listupGitSources is a go-routine which lists up git souces and puts
// finding paths to a given ch. After listing up all sources, put nil
// to done. If an error occurs, put the error to done.
func listupGitSources(ch chan<- string, done chan<- error) {

	cmd := exec.Command("git", "ls-files")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		done <- err
		return
	}
	defer stdout.Close()

	doneRL := make(chan error)
	defer close(doneRL)

	go readLine(stdout, ch, doneRL)
	if err = cmd.Run(); err != nil {
		done <- err
		return
	}

	err = <-doneRL
	done <- err
	return

}

// tarballing is a go-routine which write a file given via ch to a tar writer.
// It puts nil to done when it ends. If an error occurs, it puts the error to
// done.
func tarballing(writer *tar.Writer, ch <-chan string, done chan<- error) {

	var info os.FileInfo
	var header *tar.Header
	var err error

	for {
		path, ok := <-ch
		if !ok {
			break
		}

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

// readLine is a go-routine which reads a line from rd and puts it to ch.
// It sends nil to done when reads all line or somr error when an error
// occurs.
func readLine(rd io.Reader, ch chan<- string, done chan<- error) {

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		ch <- scanner.Text()
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
