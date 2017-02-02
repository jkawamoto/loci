//
// command/travis.go
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
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"

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
	BeforeInstall ListOrString `yaml:"before_install,omitempty"`

	// List of commands used to install packages.
	Install ListOrString `yaml:"install,omitempty"`

	// List of commands run before main scripts.
	BeforeScript ListOrString `yaml:"before_script,omitempty"`

	// List of scripts.
	Script ListOrString `yaml:"script,omitempty"`

	// List of environment variables.
	Env Env `yaml:"env,omitempty"`

	// Configuration for matrix build.
	Matrix Matrix `yaml:"matrix,omitempty"`

	// List of python versions. (used only in python)
	Python []string `yaml:"python,omitempty"`

	// List of golang versions. (used only in go)
	Go []string `yaml:"go,omitempty"`
	// Go import path. (used only in go)
	GoImportPath string `yaml:"go_import_path,omitempty"`
	// Build args for go project. (used only in go)
	GoBuildArgs string `yaml:"gobuild_args,omitempty"`
}

// Env defines the full structure of a definition of environment variables.
type Env struct {
	Global []string `yaml:"global,omitempty"`
	Matrix []string `yaml:"matrix,omitempty"`
}

// Matrix defines the structure of matrix element in .travis.yml.
type Matrix struct {
	Include []interface{} `yaml:"include,omitempty"`
	Exclude []interface{} `yaml:"exclude,omitempty"`
}

// TestCaseSet defines a set of arguments for build matrix.
type TestCaseSet map[string][][]string

// NewTravis creates a Travis object from a byte array.
func NewTravis(buf []byte) (res *Travis, err error) {

	res = &Travis{}
	if err = yaml.Unmarshal(buf, res); err != nil {
		return
	}
	return

}

// NewTravisFromFile creates a Travis object from a file.
func NewTravisFromFile(filename string) (res *Travis, err error) {

	fp, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}
	return NewTravis(buf)

}

// ArgumentSet returns a set of arguments to run entrypoint based on a build
// matrix.
func (t *Travis) ArgumentSet() (res TestCaseSet, err error) {

	switch t.Language {
	case "python":
		res, err = t.argumentSetPython()
	case "go":
		res, err = t.argumentSetGo()
	default:
		res = make(TestCaseSet)
		res[""] = [][]string{}
	}

	return

}

// UnmarshalYAML defines a way to unmarshal variables of Env.
func (e *Env) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux interface{}
	if err = unmarshal(&aux); err != nil {
		return
	}

	switch raw := aux.(type) {
	case []interface{}:
		if len(raw) == 0 {
			return
		}
		value := make([]string, len(raw))
		for i, r := range raw {
			v, ok := r.(string)
			if !ok {
				return fmt.Errorf("An item in evn cannot be converted to a string: %v", aux)
			}
			value[i] = v
		}
		if len(strings.Split(strings.TrimSpace(value[0]), " ")) == 1 {
			e.Global = value
		} else {
			e.Matrix = value
		}

	case map[interface{}]interface{}:
		if err = mapstructure.Decode(raw, e); err != nil {
			return err
		}

	}

	return

}
