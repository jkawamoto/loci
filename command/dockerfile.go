//
// command/dockerfile.go
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
	Archive string
}

// NewDockerfile creates a Dockerfile from an instance of Travis.
func NewDockerfile(travis *Travis, archive string) (res []byte, err error) {

	data, err := Asset(DockerfileAsset)
	if err != nil {
		return
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	param := travisExt{
		Travis:  travis,
		Archive: archive,
	}

	buf := bytes.Buffer{}
	if err = temp.Execute(&buf, &param); err != nil {
		return
	}

	res = buf.Bytes()
	return

}

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

func Start(tag string) (err error) {

	cmd := exec.Command("docker", "run", "-t", "--rm", tag)
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
