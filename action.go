package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

const SourceArchive = "source.tar.gz"

func Run(c *cli.Context) error {

	filename := c.Args().First()
	if err := run(filename); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func run(filename string) (err error) {

	travis, err := NewTravis(filename)
	if err != nil {
		return
	}

	tempDir := CreateTempDir()
	fmt.Println(tempDir)
	// defer os.RemoveAll(tempDir)

	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	archive := filepath.Join(tempDir, SourceArchive)
	if err = Archive(pwd, archive, nil); err != nil {
		return
	}

	docker, err := NewDockerfile(travis, archive)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "Dockerfile"), docker, 0644); err != nil {
		return
	}
	// fmt.Println(string(docker))

	entry, err := Entrypoint(travis)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "entrypoint.sh"), entry, 0644); err != nil {
		return
	}
	// fmt.Println(string(entry))

	return

}

// CreateTempDir creates a temporaty directory.
func CreateTempDir() (res string) {

	for i := 0; i < 10; i++ {

		res = filepath.Join(os.TempDir(), "loci", time.Now().Format("20060102150405"))
		if err := os.Mkdir(res, 0777); err == nil {
			return
		}
	}
	return
}
