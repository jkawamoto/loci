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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
func Build(ctx context.Context, dir, tag string) (err error) {

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
		archiveContext(ctx, dir, writer)
		writer.Close()
	}()

	// Start to build an image.
	res, err := cli.ImageBuild(ctx, reader, types.ImageBuildOptions{
		Tags: []string{tag},
	})
	if err != nil {
		return
	}
	defer res.Body.Close()

	// Wait untile the copy ends or the context will be canceled.
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stderr, res.Body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return
	}

}

// Start runs a container to run tests.
func Start(ctx context.Context, tag, name string, args ...string) (err error) {

	// Create a docker client.
	cli, err := client.NewClient(client.DefaultDockerHost, "", nil, nil)
	if err != nil {
		return
	}
	defer cli.Close()

	// Create a docker container.
	config := container.Config{
		Image: tag,
		Cmd:   args,
	}
	container, err := cli.ContainerCreate(ctx, &config, nil, nil, name)
	if err != nil {
		return
	}
	if name != "" {
		// If any container name isn't given, remove the container.
		defer cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
	}

	// Attach stdout of the container.
	stdout, err := cli.ContainerAttach(ctx, container.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
	})
	if err != nil {
		return
	}
	defer stdout.Close()
	go io.Copy(os.Stdout, stdout.Reader)

	// Attach stderr of the container.
	stderr, err := cli.ContainerAttach(ctx, container.ID, types.ContainerAttachOptions{
		Stream: true,
		Stderr: true,
	})
	if err != nil {
		return
	}
	defer stderr.Close()
	go io.Copy(os.Stderr, stderr.Reader)

	// Start the container.
	options := types.ContainerStartOptions{}
	if err = cli.ContainerStart(ctx, container.ID, options); err != nil {
		return
	}

	// Wait until the container ends.
	exit, err := cli.ContainerWait(ctx, container.ID)
	if err != nil {
		return
	} else if exit != 0 {
		return fmt.Errorf("Testing container returns an error:", exit)
	}

	return

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
	// TODO: Use ioutil.ReadDir
	sources, err := filepath.Glob(filepath.Join(root, "*"))
	if err != nil {
		return
	}
	for _, path := range sources {

		// Write a file header.
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Cannot find %s (%s)", path, err.Error())
			break
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tarWriter.WriteHeader(header)

		// Write the body.
		if err = copyFile(path, tarWriter); err != nil {
			break
		}

	}

	return

}
