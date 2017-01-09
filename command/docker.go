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
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
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

	cd, err := os.Getwd()
	if err != nil {
		return
	}
	if err = os.Chdir(dir); err != nil {
		return
	}
	defer os.Chdir(cd)

	cmd := exec.CommandContext(ctx, "docker", "build", "-t", tag, ".")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return cmd.Run()

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
	host := container.HostConfig{}
	nw := network.NetworkingConfig{}
	container, err := cli.ContainerCreate(ctx, &config, &host, &nw, name)
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
