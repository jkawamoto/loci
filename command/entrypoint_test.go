//
// command/entrypoint_test.go
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

func TestEntrypointPython(t *testing.T) {

	travis := Travis{
		Language:     "python",
		Install:      []string{"install_1", "install_2"},
		BeforeScript: []string{"before_1", "before_2"},
		Script:       []string{"script_1", "script_2"},
	}

	res, err := Entrypoint(&travis)
	if err != nil {
		t.Error("Entrypoint returns an error:", err.Error())
	}
	e := string(res)

	if !strings.Contains(e, travis.Install[0]) || !strings.Contains(e, travis.Install[1]) {
		t.Error("Entrypoint doesn't have correct install steps:", e)
	}

	if !strings.Contains(e, travis.BeforeScript[0]) || !strings.Contains(e, travis.BeforeScript[1]) {
		t.Error("Entrypoint doesn't have correct before script steps:", e)
	}

	if !strings.Contains(e, travis.Script[0]) || !strings.Contains(e, travis.Script[1]) {
		t.Error("Entrypoint doesn't have correct script steps:", e)
	}

}

func TestEntrypointGo(t *testing.T) {

	travis := Travis{
		Language:     "go",
		Install:      []string{"install_1", "install_2"},
		BeforeScript: []string{"before_1", "before_2"},
		Script:       []string{"script_1", "script_2"},
	}

	res, err := Entrypoint(&travis)
	if err != nil {
		t.Error("Entrypoint returns an error:", err.Error())
	}
	e := string(res)

	if !strings.Contains(e, travis.Install[0]) || !strings.Contains(e, travis.Install[1]) {
		t.Error("Entrypoint doesn't have correct install steps:", e)
	}

	if !strings.Contains(e, travis.BeforeScript[0]) || !strings.Contains(e, travis.BeforeScript[1]) {
		t.Error("Entrypoint doesn't have correct before script steps:", e)
	}

	if !strings.Contains(e, travis.Script[0]) || !strings.Contains(e, travis.Script[1]) {
		t.Error("Entrypoint doesn't have correct script steps:", e)
	}

}

func TestEntrypointGoByDefault(t *testing.T) {

	travis := Travis{
		Language: "go",
	}

	res, err := Entrypoint(&travis)
	if err != nil {
		t.Error("Entrypoint returns an error:", err.Error())
	}
	e := string(res)

	if !strings.Contains(e, "go get -t ./...") {
		t.Error("Entrypoint doesn't have the default install steps:", e)
	}

	if !strings.Contains(e, "go test -v ./...") {
		t.Error("Entrypoint doesn't have the default script steps:", e)
	}

}

func TestEntrypointGoWithGoBuildArgs(t *testing.T) {

	travis := Travis{
		Language:    "go",
		GoBuildArgs: "-a -b -c -d -e",
	}

	res, err := Entrypoint(&travis)
	if err != nil {
		t.Error("Entrypoint returns an error:", err.Error())
	}
	e := string(res)

	if !strings.Contains(e, fmt.Sprintf("go test %s", travis.GoBuildArgs)) {
		t.Error("Entrypoint doesn't have go build args in the script step:", e)
	}

}
