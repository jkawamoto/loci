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

// NewDockerfile creates a Dockerfile from an instance of Travis.
func NewDockerfile(travis *Travis, opt *DockerfileOpt, archive string) (res []byte, err error) {

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
	if opt.BaseImage == "" {
		opt.BaseImage = DefaultBaseImage
	}
	opt.PypiProxy = strings.TrimSuffix(opt.PypiProxy, "/")

	param := travisExt{
		Travis:        travis,
		DockerfileOpt: opt,
		Archive:       archive,
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
	if err = temp.ExecuteTemplate(&buf, "base", &param); err != nil {
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
