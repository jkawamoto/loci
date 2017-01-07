//
// command/run.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
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
	"path"
	"path/filepath"
	"strings"

	gitconfig "github.com/tcnksm/go-gitconfig"
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

	// Load a Travis's script file.
	if opt.Filename == "" {
		opt.Filename = ".travis.yml"
	}
	travis, err := NewTravis(opt.Filename)
	if err != nil {
		return
	}

	// Get repository information.
	origin, err := gitconfig.OriginURL()
	if err != nil {
		return
	}
	opt.Repository = getRepository(origin)

	// Set up the tag name of the container image.
	if opt.Tag == "" {
		opt.Tag = fmt.Sprintf("loci/%s", path.Base(opt.Repository))
	}
	tempDir := filepath.Join(os.TempDir(), opt.Tag)
	if err = os.MkdirAll(tempDir, 0777); err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	// Archive source files.
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println(chalk.Bold.TextStyle("Creating archive of source codes."))
	if err = Archive(pwd, filepath.Join(tempDir, SourceArchive)); err != nil {
		return
	}

	// Create Dockerfile.
	fmt.Println(chalk.Bold.TextStyle("Creating Dockerfile"))
	docker, err := Dockerfile(travis, opt.DockerfileOpt, SourceArchive)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "Dockerfile"), docker, 0644); err != nil {
		return
	}
	if opt.Verbose {
		fmt.Println(string(docker))
	}

	// Create entrypoint.sh.
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

	// Build the container image.
	fmt.Println(chalk.Bold.TextStyle("Building a image."))
	err = Build(tempDir, opt.Tag)
	if err != nil {
		return
	}

	// Run tests in sandboxes.
	fmt.Println(chalk.Bold.TextStyle("Start CI."))
	argset, err := travis.ArgumentSet()
	if err != nil {
		return
	}
	for i, args := range argset {
		name := opt.Name
		if name != "" {
			name = fmt.Sprintf("%s-%d", name, i+1)
		}
		err := Start(opt.Tag, name, args.Version, args.Env)
		if err != nil {
			return err
		}
	}

	return nil
}

// getRepository returns the repository path from a given remote URL of
// origin repository. The repository path consists of a URL without
// sheme, user name, password, and .git suffix.
func getRepository(origin string) (res string) {

	switch {
	case strings.Contains(origin, "@"):
		res = strings.Replace(strings.Split(origin, "@")[1], ":", "/", 1)
	case strings.HasPrefix(origin, "http://"):
		res = origin[len("http://"):]
	case strings.HasPrefix(origin, "https://"):
		res = origin[len("https://"):]
	default:
		res = strings.Replace(origin, ":", "/", 1)
	}
	if strings.HasSuffix(res, ".git") {
		res = res[:len(res)-len(".git")]
	}

	return

}
