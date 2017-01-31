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
	BeforeInstall []string `yaml:"before_install,omitempty"`

	// TODO: The Install section can be a string not a list.
	// -> Use interface{} at first and case it to some other variables.
	// List of commands used to install packages.
	Install []string `yaml:"install,omitempty"`
	// List of commands run before main scripts.
	BeforeScript []string `yaml:"before_script,omitempty"`
	// TODO: The Script section can be a string instead of a list.
	// Use RasScript interface{} to recieve items then parse to and store here.
	// List of scripts.
	Script []string `yaml:"script,omitempty"`

	// RawEnv defines a temporary space to store env attribute for parseEnv.
	RawEnv interface{} `yaml:"env,omitempty"`
	// List of environment variables.
	Env Env `yaml:"_env,omitempty"`

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
	err = res.parseEnv()
	return

}

// ArgumentSet returns a set of arguments to run entrypoint based on a build
// matrix.
func (t *Travis) ArgumentSet() (res TestCaseSet, err error) {

	switch t.Language {
	// case "python":
	// 	res, err = t.argumentSetPython()
	case "go":
		res, err = t.argumentSetGo()
	default:
		// res = []Arguments{
		// 	Arguments{},
		// }
	}

	return

}

func (t *Travis) parseEnv() (err error) {

	switch raw := t.RawEnv.(type) {
	case []interface{}:
		if len(raw) == 0 {
			return
		}
		value := make([]string, len(raw))
		for i, r := range raw {
			v, ok := r.(string)
			if !ok {
				return fmt.Errorf("An item in evn cannot be converted to a string: %v", t.RawEnv)
			}
			value[i] = v
		}
		if len(strings.Split(strings.TrimSpace(value[0]), " ")) == 1 {
			t.Env.Global = value
		} else {
			t.Env.Matrix = value
		}

	case map[interface{}]interface{}:
		if err := mapstructure.Decode(raw, &t.Env); err != nil {
			return err
		}

	}

	return

}

// Arguments defines a set of arguments for build matrix.
// type Arguments struct {
// 	// Version of the runtime to be run.
// 	Version string
// 	// Evn variables; each variable invokes one container.
// 	Env [][]string
// }

type TestCaseSet map[string][][]string

// // String method returns a string format of an Arguments.
// func (a Arguments) String() string {
// 	return fmt.Sprintf("%s %s", a.Version, a.Env)
// }
