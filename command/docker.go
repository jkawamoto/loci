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
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/tcnksm/go-gitconfig"
)

// DockerfileAsset defines a asset name for Dockerfile.
const DockerfileAsset = "assets/Dockerfile"

// DefaultBaseImage is the default base image.
const DefaultBaseImage = "ubuntu:latest"

type travisExt struct {
	*Travis
	Archive    string
	BaseImage  string
	Repository string
}

// NewDockerfile creates a Dockerfile from an instance of Travis.
func NewDockerfile(travis *Travis, base, archive string) (res []byte, err error) {

	var data []byte
	name := fmt.Sprintf("%s-%s", DockerfileAsset, travis.Language)
	data, err = Asset(name)
	if err != nil {
		data, err = Asset(DockerfileAsset)
		if err != nil {
			return
		}
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	if base == "" {
		base = DefaultBaseImage
	}

	param := travisExt{
		Travis:    travis,
		Archive:   archive,
		BaseImage: base,
	}

	origin, err := gitconfig.OriginURL()
	if err != nil {
		return
	}
	switch {
	case strings.HasPrefix(origin, "http://"):
		param.Repository = origin[len("http://"):]
	case strings.HasPrefix(origin, "https://"):
		param.Repository = origin[len("https://"):]
	case strings.Contains(origin, "@"):
		param.Repository = strings.Replace(strings.Split(origin, "@")[1], ":", "/", 1)
	default:
		param.Repository = strings.Replace(origin, ":", "/", 1)
	}
	if strings.HasSuffix(param.Repository, ".git") {
		param.Repository = param.Repository[:len(param.Repository)-len(".git")]
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
