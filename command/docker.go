//
// command/docker.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"text/template"
)

// DockerfileAsset defines a asset name for Dockerfile.
const DockerfileAsset = "asset/Dockerfile"

type travisExt struct {
	*Travis
	Archive   string
	BaseImage string
}

// NewDockerfile creates a Dockerfile from an instance of Travis.
func NewDockerfile(travis *Travis, base, archive string) (res []byte, err error) {

	data, err := Asset(DockerfileAsset)
	if err != nil {
		return
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	if base == "" {
		base = "ubuntu:latest"
	}

	param := travisExt{
		Travis:    travis,
		Archive:   archive,
		BaseImage: base,
	}

	buf := bytes.Buffer{}
	if err = temp.Execute(&buf, &param); err != nil {
		return
	}

	res = buf.Bytes()
	return

}

// Build builds a docker image from a directory. The built image named tag.
// The directory must have Dockerfile.
func Build(dir, tag string) (err error) {

	cd, err := os.Getwd()
	if err != nil {
		return
	}
	if err = os.Chdir(dir); err != nil {
		return
	}
	defer os.Chdir(cd)

	cmd := exec.Command("docker", "build", "-t", tag, ".")
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
func Start(tag, name string) (err error) {

	var cmd *exec.Cmd
	if name == "" {
		cmd = exec.Command("docker", "run", "-t", "--rm", tag)
	} else {
		cmd = exec.Command("docker", "run", "-t", "--name", name, tag)
	}

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
