//
// command/archive_test.go
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
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestArchive(t *testing.T) {

	temp := os.TempDir()
	target := path.Join(temp, "test.tar.gz")
	t.Logf("Creating an archive file: %s", target)

	if err := Archive("..", target); err != nil {
		t.Error(err.Error())
		return
	}
	if _, err := os.Stat(target); err != nil {
		t.Error(err.Error())
		return
	}
	// defer os.Remove(target)

	fp, err := os.Open(target)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer fp.Close()

	zip, err := gzip.NewReader(fp)
	if err != nil {
		t.Error(err.Error())
		return
	}

	reader := tar.NewReader(zip)
	for {

		info, err := reader.Next()
		if err != io.EOF {
			break
		} else if err != nil {
			t.Error(err.Error())
			return
		}

		original, err := os.Stat(filepath.Join("..", info.Name))
		if err != nil {
			t.Error(err.Error())
			return
		}
		if info.Size != original.Size() {
			t.Errorf("%s seems broken. (%d != %d)", info.Name, info.Size, original.Size())
			return
		}

	}

}
