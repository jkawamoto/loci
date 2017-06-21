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
	"syscall"

	"golang.org/x/sync/errgroup"

	gitconfig "github.com/tcnksm/go-gitconfig"
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
	// Max processors to be used.
	Processors int
	// If true, print Dockerfile and entrypoint.sh.
	Verbose bool
	// If true, not using cache during buidling a docker image.
	NoCache bool
	// If true, omit printing color codes.
	NoColor bool
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
		Tag:        c.String("tag"),
		Processors: c.Int("max-processors"),
		Verbose:    c.Bool("verbose"),
		NoCache:    c.Bool("no-cache"),
		NoColor:    c.Bool("no-color"),
	}
	if err := run(&opt); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func run(opt *RunOpt) (err error) {

	var decorator *TextDecorator
	if opt.NoColor {
		decorator = NewNoopDecorator()
	} else {
		decorator = NewDecorator()
	}

	// Prepare to be canceled.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGQUIT)
	go func() {
		<-sig
		cancel()
	}()

	// Load a Travis's script file.
	if opt.Filename == "" {
		opt.Filename = ".travis.yml"
	}
	travis, err := NewTravisFromFile(opt.Filename)
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
		opt.Tag = fmt.Sprintf("loci/%s", strings.ToLower(path.Base(opt.Repository)))
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
	fmt.Println(decorator.Bold("Creating archive of source codes."))
	var dstout io.Writer
	if opt.Verbose {
		dstout = os.Stdout
	} else {
		dstout = ioutil.Discard
	}
	if err = Archive(ctx, pwd, filepath.Join(tempDir, SourceArchive), dstout, os.Stderr); err != nil {
		return
	}

	// Create Dockerfile.
	fmt.Println(decorator.Bold("Creating Dockerfile"))
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
	fmt.Println(decorator.Bold("Creating entrypoint."))
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

	argset, err := travis.ArgumentSet()
	if err != nil {
		return
	}

	// Start testing with goroutines.
	fmt.Println(decorator.Bold("Start testing."))
	wg, ctx := errgroup.WithContext(ctx)
	semaphore := make(chan struct{}, opt.Processors)
	display, err := NewDisplay()
	if err != nil {
		return
	}

	var i int
	for version, set := range argset {

		func(version string, set [][]string) {
			wg.Go(func() (err error) {
				semaphore <- struct{}{}
				defer func() {
					<-semaphore
				}()

				// Build a container image.
				fp, err := os.OpenFile(fmt.Sprintf("loci-build-v%v.log", version), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					return
				}
				defer fp.Close()

				sec := display.AddSection(fmt.Sprintf("Building a image for v%v", version))
				defer display.DeleteSection(sec)

				output := sec.Writer()
				defer output.Close()

				tag := fmt.Sprintf("%v/%v", opt.Tag, version)
				err = Build(ctx, tempDir, tag, version, opt.NoCache, io.MultiWriter(fp, output))
				if err != nil {
					return
				}

				for _, envs := range set {

					func(envs []string) {
						wg.Go(func() (err error) {
							semaphore <- struct{}{}
							defer func() {
								<-semaphore
							}()

							// Run tests in a sandbox.
							fp, err := os.OpenFile(
								fmt.Sprintf("loci-v%v.log", strings.Join(append([]string{version}, envs...), "-")),
								os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
							if err != nil {
								return
							}
							defer fp.Close()

							sec := display.AddSection(fmt.Sprintf("Running tests (v%v: %v)", version, envs))
							defer display.DeleteSection(sec)

							output := sec.Writer()
							defer output.Close()

							name := opt.Name
							if name != "" {
								i++
								name = fmt.Sprintf("%s-%d", name, i)
							}

							err = Start(ctx, tag, name, envs, io.MultiWriter(fp, output))
							if err != nil {
								return fmt.Errorf("%s\n%s", err, sec.String())
							}
							return

						})
					}(envs)

				}

				return

			})
		}(version, set)

	}

	err = wg.Wait()
	display.Close()
	if err == nil {
		fmt.Println(decorator.Bold("All tests have been passed."))
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
