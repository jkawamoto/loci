package command

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const archiveFile = "test.tar.gz"

func TestArchive(t *testing.T) {

	temp := os.TempDir()
	target := path.Join(temp, archiveFile)
	t.Logf("Creating an archive file: %s", target)

	if err := Archive("..", target, []string{"*.tar"}); err != nil {
		t.Error(err.Error())
	}

	root, err := os.Getwd()
	if err != nil {
		t.Error(err.Error())
	}

	os.Chdir(temp)
	defer func() {
		os.Chdir(root)
	}()
	exec.Command("tar", "-zxvf", archiveFile)

	if err := filepath.Walk(path.Join(root, ".."), checkExistence(temp)); err != nil {
		t.Error(err.Error())
	}

}

func checkExistence(target string) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			if strings.HasSuffix(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		} else if strings.HasSuffix(path, ".tar") {
			return nil
		}

		_, check := os.Stat(path)
		return check

	}

}
