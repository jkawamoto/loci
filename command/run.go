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

// SourceArchive defines a name of source archive file.
const SourceArchive = "source.tar.gz"

// RunOpt defines a option parameter for run function.
type RunOpt struct {
	// Same options as DockerfileOpt.
	*DockerfileOpt
	// Travis configuration file.
	Filename string
	// Container name.
	Name string
	// Image tag.
	Tag string
	// If true, print Dockerfile and entrypoint.sh.
	Verbose bool
}

// Run implements the action of this command.
func Run(c *cli.Context) error {

	opt := RunOpt{
		DockerfileOpt: &DockerfileOpt{
			BaseImage:  c.String("base"),
			AptProxy:   c.String("apt-proxy"),
			PypiProxy:  c.String("pypi-proxy"),
			HTTPProxy:  c.String("http-proxy"),
			HTTPSProxy: c.String("https-proxy"),
			NoProxy:    c.String("no-proxy"),
		},
		Filename: c.Args().First(),
		Name:     c.String("name"),
		Tag:      c.String("tag"),
		Verbose:  c.Bool("verbose"),
	}
	if err := run(&opt); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func run(opt *RunOpt) (err error) {

	if opt.Filename == "" {
		opt.Filename = ".travis.yml"
	}
	travis, err := NewTravis(opt.Filename)
	if err != nil {
		return
	}

	if opt.Tag == "" {
		opt.Tag = fmt.Sprintf("loci/%s", time.Now().Format("20060102150405"))
	}
	tempDir := filepath.Join(os.TempDir(), opt.Tag)
	if err = os.MkdirAll(tempDir, 0777); err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

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
	docker, err := NewDockerfile(travis, opt.DockerfileOpt, archive)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "Dockerfile"), docker, 0644); err != nil {
		return
	}
	if opt.Verbose {
		fmt.Println(string(docker))
	}

	fmt.Println(chalk.Bold.TextStyle("Creating entrypoint."))
	entry, err := Entrypoint(travis)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "entrypoint.sh"), entry, 0644); err != nil {
		return
	}
	if opt.Verbose {
		fmt.Println(string(entry))
	}

	fmt.Println(chalk.Bold.TextStyle("Building a image."))
	err = Build(tempDir, opt.Tag)
	if err != nil {
		return
	}
	fmt.Println(chalk.Bold.TextStyle("Start CI."))
	return Start(opt.Tag, opt.Name)

}
