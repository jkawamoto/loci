//
// command/travis.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Travis defines the structure of .travis.yml.
type Travis struct {
	// Base language.
	Language string
	// List of addons.
	Addons struct {
		Apt struct {
			Packages []string
		} `yaml:"apt,omitempty"`
	} `yaml:"addons,omitempty"`
	// List of commands run before install steps.
	BeforeInstall []string `yaml:"before_install,omitempty"`
	// TODO: The Install section can be a string not a list.
	// List of commands used to install packages.
	Install []string `yaml:"install,omitempty"`
	// List of commands run before main scripts.
	BeforeScript []string `yaml:"before_script,omitempty"`
	// List of scripts.
	Script []string `yaml:"script,omitempty"`
	// List of environment variables.
	Env []string `yaml:"env,omitempty"`
	// Configuration for matrix build.
	Matrix Matrix `yaml:"matrix,omitempty"`

	// List of python versions. (used only in python)
	Python []string `yaml:"python,omitempty"`

	// List of golang versions. (used only in go)
	Go []string `yaml:"go,omitempty"`
	// Go import path. (used only in go)
	GoImportPath string `yaml:"go_import_path,omitempty"`
}

// Matrix defines the structure of matrix element in .travis.yml.
type Matrix struct {
	Include []interface{} `yaml:"include,omitempty"`
	Exclude []interface{} `yaml:"exclude,omitempty"`
}

// NewTravis loads a .travis.yaml file and creates a structure instance.
func NewTravis(filename string) (res *Travis, err error) {

	fp, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	res = &Travis{}
	if err = yaml.Unmarshal(buf, res); err != nil {
		return
	}
	return

}

// ArgumentSet returns a set of arguments to run entrypoint based on a build
// matrix.
func (t *Travis) ArgumentSet() [][]string {

	switch t.Language {
	case "python":
		return t.argumentSetPython()
	default:
		return [][]string{{""}}
	}

}

// parseEnv parses a string consisting of a name of an environment variable
// and its value by concatinating with =, and returns a tuple of the name and
// value.
func parseEnv(env string) (string, string) {
	s := strings.SplitN(env, "=", 2)
	return s[0], s[1]
}

// makeSetKey returns a string consisting of version, key, and value.
func makeSetKey(version, key, value string) string {
	return fmt.Sprintf("%s %s %s", version, key, value)
}
