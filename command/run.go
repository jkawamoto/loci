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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	colorable "github.com/mattn/go-colorable"
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
	// Runtime version to which only versions matching will be run.
	Version string
	// Image tag.
	Tag string
	// Max processors to be used.
	Processors int
	// If true, logging information to be stored to files.
	OutputLog bool
	// If true, not using cache during buidling a docker image.
	NoCache bool
	// If true, omit printing color codes.
	NoColor bool
	// Printed on the header.
	Title string
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
		Filename:   c.Args().First(),
		Name:       c.String("name"),
		Version:    c.String("select"),
		Tag:        c.String("tag"),
		Processors: c.Int("max-processors"),
		OutputLog:  c.Bool("log"),
		NoCache:    c.Bool("no-cache"),
		NoColor:    c.Bool("no-color"),
		Title:      fmt.Sprintf("%v %v", c.App.Name, c.App.Version),
	}
	if err := run(&opt); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func run(opt *RunOpt) (err error) {

	// Prepare to be canceled.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGQUIT)
	go func() {
		<-sig
		cancel()
	}()

	// Prepare interface.
	display, ctx, err := NewDisplay(ctx, opt.Title, opt.Processors)
	if err != nil {
		return
	}
	defer display.Close()
	logger := display.Header.Logger

	var stdout io.Writer
	if opt.NoColor {
		stdout = colorable.NewNonColorable(os.Stdout)
		logger = colorable.NewNonColorable(logger)
		cli.ErrWriter = colorable.NewNonColorable(cli.ErrWriter)
	} else {
		stdout = colorable.NewColorableStdout()
	}

	// Load a Travis's script file.
	if opt.Filename == "" {
		opt.Filename = ".travis.yml"
	}
	fmt.Fprintln(logger, chalk.Cyan.Color("Loading .travis.yml"))

	travis, err := NewTravisFromFile(opt.Filename)
	if err != nil {
		return
	}

	// Get repository information.
	fmt.Fprintln(logger, chalk.Cyan.Color("Checking repository information"))
	origin, err := gitconfig.OriginURL()
	if err != nil {
		return
	}
	opt.Repository = getRepository(origin)

	// Set up the tag name of the container image.
	if opt.Tag == "" {
		opt.Tag = fmt.Sprintf("loci/%s", strings.ToLower(path.Base(opt.Repository)))
	}

	// Prepare docker images.
	fmt.Fprintln(logger, chalk.Cyan.Color("Preparing docker images for sandbox containers"))
	err = PrepareBaseImage(ctx, opt.BaseImage, logger)
	if err != nil {
		return
	}

	// Archive source files.
	fmt.Fprintln(logger, chalk.Cyan.Color("Archiving source code"))
	tempDir := filepath.Join(os.TempDir(), opt.Tag)
	if err = os.MkdirAll(tempDir, 0777); err != nil {
		return
	}
	defer os.RemoveAll(tempDir)
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	if err = Archive(ctx, pwd, filepath.Join(tempDir, SourceArchive)); err != nil {
		return
	}

	// Create Dockerfile.
	fmt.Fprintln(logger, chalk.Cyan.Color("Creating Dockerfile"))
	docker, err := Dockerfile(travis, opt.DockerfileOpt, SourceArchive)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "Dockerfile"), docker, 0644); err != nil {
		return
	}

	// Create entrypoint.sh.
	fmt.Fprintln(logger, chalk.Cyan.Color("Creating entrypoint.sh"))
	entry, err := Entrypoint(travis)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(tempDir, "entrypoint.sh"), entry, 0644); err != nil {
		return
	}

	argset, err := travis.ArgumentSet(logger)
	if err != nil {
		return
	}

	// Start testing with goroutines.
	fmt.Fprintln(logger, chalk.Cyan.Color("Building sandbox images and running tests"))
	var i int
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, opt.Processors)
	errs := NewErrorSet()
	for version, set := range argset {

		if opt.Version != "" && version != opt.Version {
			continue
		}

		wg.Add(1)
		go func(version string, set [][]string) (err error) {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wg.Done()
			}()

			// Build a container image.
			sec := display.AddSection(fmt.Sprintf("Building a image for %v", version))
			defer display.DeleteSection(sec)

			var output io.Writer
			writer := sec.Writer()
			defer writer.Close()
			output = writer
			if opt.NoColor {
				output = colorable.NewNonColorable(output)
			}

			if opt.OutputLog {
				var fp *os.File
				fp, err = os.OpenFile(fmt.Sprintf("loci-build-%v.log", version), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					errs.Add(version, err)
					return
				}
				defer fp.Close()
				output = io.MultiWriter(output, colorable.NewColorable(fp))
			}

			tag := fmt.Sprintf("%v/%v", opt.Tag, version)
			err = Build(ctx, tempDir, tag, version, opt.NoCache, output)
			if err == context.Canceled {
				errs.Add("", err)
				return
			} else if err != nil {
				msg := fmt.Sprintf(chalk.Red.Color("Faild to build a docker image for %v"), version)
				errs.Add(
					version,
					fmt.Errorf("%v\n%v\n%v\n", msg, err.Error(), sec.String()))
				fmt.Fprintln(logger, msg)
				return
			}
			fmt.Fprintln(logger, chalk.Green.Color(fmt.Sprintf("Built a image for %v", version)))

			for _, envs := range set {

				wg.Add(1)
				go func(envs []string) {
					semaphore <- struct{}{}
					defer func() {
						<-semaphore
						wg.Done()
					}()

					// Run tests in a sandbox.
					sec := display.AddSection(fmt.Sprintf("Running tests (%v: %v)", version, envs))
					defer display.DeleteSection(sec)

					var output io.Writer
					writer := sec.Writer()
					defer writer.Close()
					output = writer
					if opt.NoColor {
						output = colorable.NewNonColorable(output)
					}

					if opt.OutputLog {
						var fp *os.File
						fp, err = os.OpenFile(
							fmt.Sprintf("loci-%v.log", strings.Join(append([]string{version}, envs...), "-")),
							os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
						if err != nil {
							errs.Add(fmt.Sprintf("%v:%v", version, envs), err)
							return
						}
						defer fp.Close()
						output = io.MultiWriter(output, colorable.NewColorable(fp))
					}

					name := opt.Name
					if name != "" {
						i++
						name = fmt.Sprintf("%s-%d", name, i)
					}

					err = Start(ctx, tag, name, envs, output)
					if err == context.Canceled {
						errs.Add("", err)
					} else if err != nil {
						errs.Add(fmt.Sprintf("%v:%v", version, envs), fmt.Errorf("%s\n%s", chalk.Red.Color(err.Error()), sec.String()))
						fmt.Fprintln(logger, chalk.Red.Color(fmt.Sprintf("Failed tests (%v: %v) ", version, envs)))
					} else {
						fmt.Fprintln(logger, chalk.Green.Color(fmt.Sprintf("Passed tests (%v: %v) ", version, envs)))
					}
					return

				}(envs)

			}

			return

		}(version, set)

	}

	wg.Wait()
	err = display.Close()
	if err != nil {
		return
	}

	if errs.Size() == 0 {
		fmt.Fprintln(stdout, chalk.Green.Color("All tests have been passed."))
	} else {
		errList := errs.GetList()
		if errList[0] == context.Canceled {
			err = cli.NewExitError("canceled", 1)
		} else {
			err = cli.NewMultiError(errList...)
		}
	}
	return

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
