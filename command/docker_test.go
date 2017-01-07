//
// command/docker_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"fmt"
	"strings"
	"testing"
)

func TestDockerfile(t *testing.T) {

	travis := Travis{}
	opt := DockerfileOpt{
		BaseImage:  "testimage",
		AptProxy:   "http://apt.proxy.test/",
		PypiProxy:  "http://pypi.proxy.test/",
		HTTPProxy:  "http://proxy.test/",
		HTTPSProxy: "https://proxy.test/",
		NoProxy:    "localhost",
	}

	res, err := Dockerfile(&travis, &opt, SourceArchive)
	if err != nil {
		t.Error("Dockerfile returns an error:", err.Error())
	}
	dockerfile := string(res)

	if !strings.Contains(dockerfile, fmt.Sprintln("FROM", opt.BaseImage)) {
		t.Error("The base image of the Dockerfile isn't correct:", dockerfile)
	}

	if !strings.Contains(dockerfile, fmt.Sprintln("ENV HTTP_PROXY", opt.HTTPProxy)) {
		t.Error("Dockerfile doesn't set the correct http proxy:", dockerfile)
	}

	if !strings.Contains(dockerfile, fmt.Sprintln("ENV HTTPS_PROXY", opt.HTTPSProxy)) {
		t.Error("Dockerfile doesn't set the correct https proxy:", dockerfile)
	}

	if !strings.Contains(dockerfile, fmt.Sprintln("ENV NO_PROXY", opt.NoProxy)) {
		t.Error("Dockerfile doesn't set the correct no proxy env:", dockerfile)
	}

	if !strings.Contains(dockerfile, fmt.Sprintf("Acquire::http { Proxy \\\"%s\\\"; };", opt.AptProxy)) {
		t.Error("Dockerfile doesn't specify the correct apt proxy:", dockerfile)
	}

	if !strings.Contains(dockerfile, fmt.Sprintf("RUN PYPI_PROXY=%s", opt.PypiProxy)) {
		t.Error("Dockerfile doesn't specify the correct pypi proxy:", dockerfile)
	}

}

func TestDockerfileWithoutOptions(t *testing.T) {

	travis := Travis{}
	opt := DockerfileOpt{
		BaseImage: "testimage",
	}

	res, err := Dockerfile(&travis, &opt, SourceArchive)
	if err != nil {
		t.Error("Dockerfile returns an error:", err.Error())
	}
	dockerfile := string(res)

	if strings.Contains(dockerfile, "ENV HTTP_PROXY") {
		t.Error("Dockerfile set a http proxy:", dockerfile)
	}

	if strings.Contains(dockerfile, "ENV HTTPS_PROXY") {
		t.Error("Dockerfile set a https proxy:", dockerfile)
	}

	if strings.Contains(dockerfile, "ENV NO_PROXY") {
		t.Error("Dockerfile set no_proxy env:", dockerfile)
	}

	if strings.Contains(dockerfile, "Acquire::http") {
		t.Error("Dockerfile set an apt proxy:", dockerfile)
	}

	if strings.Contains(dockerfile, "RUN PYPI_PROXY=") {
		t.Error("Dockerfile set a pypi proxy:", dockerfile)
	}

}

func TestDockerfilePython(t *testing.T) {

	var travis Travis
	travis.Language = "python"
	travis.Addons.Apt.Packages = []string{"package1", "package2"}
	travis.BeforeInstall = []string{"abc", "def"}

	opt := DockerfileOpt{
		BaseImage: "ubuntu:latest",
	}
	archive := "source.tar.gz"

	res, err := Dockerfile(&travis, &opt, archive)
	if err != nil {
		t.Error("Dockerfile returns an error:", err.Error())
	}
	dockerfile := string(res)

	if !strings.Contains(dockerfile, "RUN pip install --upgrade pip") {
		t.Error("Dockerfile doesn't install pip packages:", dockerfile)
	}

	if !strings.Contains(dockerfile, "RUN apt-get install -y package1") {
		t.Error("Dockerfile doesn't install required packages:", dockerfile)
	}
	if !strings.Contains(dockerfile, "RUN apt-get install -y package2") {
		t.Error("Dockerfile doesn't install required packages:", dockerfile)
	}

	if !strings.Contains(dockerfile, "ADD source.tar.gz /data") {
		t.Error("Dockerfile doesn't add correct source files:", dockerfile)
	}

	if !strings.Contains(dockerfile, "RUN abc") || !strings.Contains(dockerfile, "RUN def") {
		t.Error("Dockerfile doesn't execute commands in before install:", dockerfile)
	}

}
