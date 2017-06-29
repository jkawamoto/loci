//
// command/archive.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
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
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

// pathListupFunc defines a function which lists up paths and put them to a
// given channel. This function is used with parallelListup.
type pathListupFunc func(context.Context, chan<- string) error

// Archive makes a tar.gz file consists of files maintained a git repository.
func Archive(ctx context.Context, dir string, filename string) (err error) {

	writeFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer writeFile.Close()

	zipWriter, err := gzip.NewWriterLevel(writeFile, gzip.BestCompression)
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
	wg, ctx := errgroup.WithContext(ctx)
	ch := make(chan string)
	wg.Go(func() error {
		defer close(ch)

		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			return listupGitRepository(ctx, ch)
		})
		eg.Go(func() error {
			return listupGitSources(ctx, ch)
		})
		return eg.Wait()

	})

	wg.Go(func() error {
		return tarballing(tarWriter, ch)
	})

	return wg.Wait()

}

// tarballing is a go-routine which write a file given via ch to a tar writer.
func tarballing(writer *tar.Writer, ch <-chan string) (err error) {

	var info os.FileInfo
	var header *tar.Header
	for path := range ch {

		// For Windows: Replace path delimiters.
		path = filepath.ToSlash(path)

		// Write a file header.
		info, err = os.Stat(path)
		if err != nil {
			err = fmt.Errorf("Cannot find %s (%s)", path, err.Error())
			break
		}

		header, err = tar.FileInfoHeader(info, path)
		if err != nil {
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

	return

}

// listupGitRepository lists up git repository and puts founded paths to a given ch.
func listupGitRepository(ctx context.Context, ch chan<- string) error {

	return filepath.Walk(".git", func(path string, info os.FileInfo, err error) error {

		// Check the given context is still alive.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// If an error is passed, propagate it.
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ch <- path
		}
		return nil

	})

}

// listupGitSources lists up git sources and puts finding paths to a given ch.
func listupGitSources(ctx context.Context, ch chan<- string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "ls-files")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	var info os.FileInfo
	s := bufio.NewScanner(stdout)
	for s.Scan() {

		path := s.Text()
		if info, err = os.Stat(path); err == nil && !info.IsDir() {
			ch <- path
		}

	}

	err = s.Err()
	if err != nil {
		return
	}

	return cmd.Wait()

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
