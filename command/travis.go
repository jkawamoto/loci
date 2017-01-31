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
	"reflect"
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

	// RawBeforeInstall defines a temporary space to store before_install attribute for parseRawField.
	RawBeforeInstall interface{} `yaml:"before_install,omitempty"`
	// List of commands run before install steps.
	BeforeInstall []string `yaml:"_before_install,omitempty"`

	// RawInstall defines a temporary space to store install attribute for parseRawField.
	RawInstall interface{} `yaml:"install,omitempty"`
	// List of commands used to install packages.
	Install []string `yaml:"_install,omitempty"`

	// RawBeforeScript defines a temporary space to store before_script attribute for parseRawField.
	RawBeforeScript interface{} `yaml:"before_script,omitempty"`
	// List of commands run before main scripts.
	BeforeScript []string `yaml:"_before_script,omitempty"`

	// RawScript defines a temporary space to store script attribute for parseRawField.
	RawScript interface{} `yaml:"script,omitempty"`
	// List of scripts.
	Script []string `yaml:"_script,omitempty"`

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
	for _, name := range []string{"BeforeInstall", "Install", "BeforeScript", "Script"} {
		err = res.parseRawField(name)
		if err != nil {
			return
		}
	}
	err = res.parseEnv()
	return

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

func (t *Travis) parseRawField(name string) (err error) {

	r := reflect.Indirect(reflect.ValueOf(t))
	src := r.FieldByName(fmt.Sprintf("Raw%s", name))
	dest := r.FieldByName(name)

	switch raw := src.Interface().(type) {
	case string:
		dest.Set(reflect.ValueOf([]string{raw}))

	case []interface{}:
		list := make([]string, len(raw))
		for i, r := range raw {
			v, ok := r.(string)
			if !ok {
				return fmt.Errorf("An item in evn cannot be converted to a string: %v", src)
			}
			list[i] = v
		}
		dest.Set(reflect.ValueOf(list))

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

// TestCaseSet defines a set of arguments for build matrix.
type TestCaseSet map[string][][]string
