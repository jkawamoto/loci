//
// command/run.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ttacon/chalk"
	"github.com/urfave/cli"
)

const SourceArchive = "source.tar.gz"

func Run(c *cli.Context) error {

	filename := c.Args().First()
	if filename == "" {
		filename = ".travis.yml"
	}
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

	tag := fmt.Sprintf("loci/%s", time.Now().Format("20060102150405"))
	tempDir := filepath.Join(os.TempDir(), tag)
	if err = os.Mkdir(tempDir, 0777); err != nil {
		return
	}
	// defer os.RemoveAll(tempDir)

	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	archive := filepath.Join(tempDir, SourceArchive)
	fmt.Println(chalk.Bold.TextStyle("Creating archive of source codes."))
	if err = Archive(pwd, archive); err != nil {
		return
	}

	fmt.Println(chalk.Bold.TextStyle("Creating Dockerfile"))
	docker, err := NewDockerfile(travis, archive)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "Dockerfile"), docker, 0644); err != nil {
		return
	}

	fmt.Println(chalk.Bold.TextStyle("Creating entrypoint."))
	entry, err := Entrypoint(travis)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "entrypoint.sh"), entry, 0644); err != nil {
		return
	}

	fmt.Println(chalk.Bold.TextStyle("Building a image."))
	err = Build(tempDir, tag)
	if err != nil {
		return
	}
	fmt.Println(chalk.Bold.TextStyle("Start CI."))
	return Start(tag)

}
