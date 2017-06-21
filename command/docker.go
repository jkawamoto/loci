//
// command/docker.go
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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	client "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerfileAsset defines a asset name for Dockerfile.
const DockerfileAsset = "assets/Dockerfile"

// DockerfileOpt defines option variables used in Dockerfile templates.
type DockerfileOpt struct {
	// Customize FROM clause.
	BaseImage string
	// Git repository.
	Repository string
	// URL for an Apt proxy.
	AptProxy string
	// URL for a pypi proxy.
	PypiProxy string
	// URL for a HTTP proxy.
	HTTPProxy string
	// URL for a HTTPS proxy.
	HTTPSProxy string
	// Comma separated URL lists.
	NoProxy string
}

type travisExt struct {
	*Travis
	*DockerfileOpt
	Archive string
}

// buildLog defines the JSON format of logs from building docker images.
type buildLog struct {
	Stream      string
	Error       string
	ErrorDetail struct {
		Code    int
		Message string
	}
}

// Dockerfile creates a Dockerfile from an instance of Travis.
func Dockerfile(travis *Travis, opt *DockerfileOpt, archive string) (res []byte, err error) {

	var data []byte

	// Loading the base template.
	data, err = Asset(DockerfileAsset)
	if err != nil {
		return
	}
	base, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	// Loading a child template.
	name := fmt.Sprintf("%s-%s", DockerfileAsset, travis.Language)
	data, err = Asset(name)
	if err != nil {
		data, err = Asset(fmt.Sprintf("%s-python", DockerfileAsset))
		if err != nil {
			return
		}
	}
	temp, err := base.Parse(string(data))
	if err != nil {
		return
	}

	// Checking optional parameters.
	opt.PypiProxy = strings.TrimSuffix(opt.PypiProxy, "/")

	// Creating Dockerfile.
	param := travisExt{
		Travis:        travis,
		DockerfileOpt: opt,
		Archive:       archive,
	}
	buf := bytes.Buffer{}
	if err = temp.ExecuteTemplate(&buf, "base", &param); err != nil {
		return
	}
	res = buf.Bytes()

	return

}

// Build builds a docker image from a directory. The built image named tag.
// The directory must have Dockerfile.
func Build(ctx context.Context, dir, tag, version string, noCache bool, output io.Writer) (err error) {

	// Create a docker client.
	cli, err := client.NewClient(client.DefaultDockerHost, "", nil, nil)
	if err != nil {
		return
	}
	defer cli.Close()

	// Create a pipe.
	reader, writer := io.Pipe()

	// Send the build context.
	go func() {
		defer writer.Close()
		archiveContext(ctx, dir, writer)
	}()

	// Start to build an image.
	res, err := cli.ImageBuild(ctx, reader, types.ImageBuildOptions{
		Tags:   []string{tag},
		Remove: true,
		BuildArgs: map[string]*string{
			"VERSION": &version,
		},
		NoCache: noCache,
	})
	if err != nil {
		return
	}
	defer res.Body.Close()

	// Wait untile the copy ends or the context will be canceled.
	done := make(chan struct{})
	go func() {
		defer close(done)

		s := bufio.NewScanner(res.Body)
		for s.Scan() {
			var log buildLog
			if json.Unmarshal(s.Bytes(), &log) == nil {
				if log.Error != "" {
					err = fmt.Errorf(log.Error)
					return
				}
				io.WriteString(output, log.Stream)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return
	}

}

// Start runs a container to run tests with a given context.
// This function creates a container from the image of the given tag name,
// and the created container has the given name. If the given name is empty,
// the container will have a random temporary name and be deleted after
// after all steps end. env is a list of environment variables to be passed
// to the created container.
func Start(ctx context.Context, tag, name string, env []string, output io.Writer) (err error) {

	// Create a docker client.
	cli, err := client.NewClient(client.DefaultDockerHost, "", nil, nil)
	if err != nil {
		return
	}
	defer cli.Close()

	// Create a docker container.
	config := container.Config{
		Image: tag,
		Env:   env,
	}
	c, err := cli.ContainerCreate(ctx, &config, nil, nil, name)
	if err != nil {
		return
	}
	if name == "" {
		// If any container name isn't given, remove the container.
		// Note that, the context ctx may be canceled before removing the container,
		// and use another context here.
		defer cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})
	}

	// Attach stdout and stderr of the container.
	stream, err := cli.ContainerAttach(ctx, c.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return
	}
	defer stream.Close()
	go stdcopy.StdCopy(output, output, stream.Reader)

	// Start the container.
	options := types.ContainerStartOptions{}
	if err = cli.ContainerStart(ctx, c.ID, options); err != nil {
		return
	}

	// Wait until the container ends.
	exit, errCh := cli.ContainerWait(ctx, c.ID, container.WaitConditionNotRunning)
	select {
	case <-ctx.Done():
		// Kill the running container when the context is canceled.
		// The context ctx has been canceled already, use another context here.
		cli.ContainerKill(context.Background(), c.ID, "")
		return ctx.Err()
	case err = <-errCh:
		// Kill the running container when ContainerWait returns an error.
		// The context ctx has been canceled already, use another context here.
		cli.ContainerKill(context.Background(), c.ID, "")
		return
	case status := <-exit:
		if status.StatusCode != 0 {
			err = fmt.Errorf("Testing container returns an error: %v", status.StatusCode)
		}
		return
	}

}

// archiveContext makes a tar.gz stream consists of files.
func archiveContext(ctx context.Context, root string, writer io.Writer) (err error) {

	// Create a buffered writer.
	bufWriter := bufio.NewWriter(writer)
	defer bufWriter.Flush()

	// Create a zipped writer on the bufferd writer.
	zipWriter, err := gzip.NewWriterLevel(bufWriter, gzip.BestCompression)
	if err != nil {
		return
	}
	defer zipWriter.Close()

	// Create a tarball writer on the zipped writer.
	tarWriter := tar.NewWriter(zipWriter)
	defer tarWriter.Close()

	// Create a tarball.
	sources, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}
	for _, info := range sources {

		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			// Write a file header.
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			tarWriter.WriteHeader(header)

			// Write the body.
			if err = copyFile(filepath.Join(root, info.Name()), tarWriter); err != nil {
				return err
			}
		}
	}

	return

}
