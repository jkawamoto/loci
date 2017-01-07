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
func Start(tag, name string, args []string) (err error) {

	var cmd *exec.Cmd
	if name == "" {
		cmd = exec.Command("docker", append([]string{"run", "-t", "--rm", tag}, args...)...)
	} else {
		cmd = exec.Command("docker", append([]string{"run", "-t", "--name", name, tag}, args...)...)
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
