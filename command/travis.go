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
	"io"
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

// TestCase is a set of environment variables and represented as a map of which
// a key is a name of one environment variable and the associated value is the
// value of the variable.
type TestCase map[string]string

// Slice returns a slice of strings representing this test case.
func (c TestCase) Slice() (res []string) {
	for key, value := range c {
		res = append(res, fmt.Sprintf("%v=%v", key, value))
	}
	return
}

// Copy returns a hard copy of this test case.
func (c TestCase) Copy() TestCase {
	res := make(TestCase)
	for k, v := range c {
		res[k] = v
	}
	return res
}

// Merge updates this TestCase so that it also has key and values defined in the
// given test case. If both test cases have a same key, the value associated
// with the key will be overwritten by the value in the given test case.
func (c TestCase) Merge(o TestCase) TestCase {
	for k, v := range o {
		c[k] = v
	}
	return c
}

// Match returns true if and only if the given TestCase has same configuration
// as this test case.
func (c TestCase) Match(o TestCase) bool {

	if len(c) != len(o) {
		return false
	}
	for k, v := range o {
		if c[k] != v {
			return false
		}
	}
	return true

}

// TestCaseSet defines a set of arguments for build matrix.
// The test case set is a map of which key is a version and the associated value
// is a list of test cases.
type TestCaseSet map[string][]TestCase

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
func (t *Travis) ArgumentSet(logger io.Writer) (res TestCaseSet, err error) {

	switch t.Language {
	case "python":
		res, err = t.argumentSetPython(logger)
	case "go":
		res, err = t.argumentSetGo()
	default:
		res = make(TestCaseSet)
		res[""] = []TestCase{}
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
		// If attribute env has one list instead of global and/or matrix attributes.
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
		// If each string has more than two variables, it means matrix configuration.
		if len(strings.Split(strings.TrimSpace(value[0]), " ")) == 1 {
			e.Global = value
		} else {
			e.Matrix = value
		}

	case map[interface{}]interface{}:
		e.Global = parseEnvMap(raw, "global")
		e.Matrix = parseEnvMap(raw, "matrix")

	}

	return

}

// parseEnvMap parses a map of which key and value are defined as interface{},
// and returns a list of strings the given map's value, which is associated with
// the given key, represents.
func parseEnvMap(m map[interface{}]interface{}, key string) (res []string) {

	if selected, exist := m[key]; exist {
		switch items := selected.(type) {
		case string:
			res = []string{items}

		case []interface{}:
			for _, v := range items {
				if s, ok := v.(string); ok {
					res = append(res, s)
				}
			}
		}
	}

	return

}

// parseEnv parses a string representing a set of environment variable
// definitions; and returns a TestCase.
func parseEnv(env string) (c TestCase) {
	c = make(TestCase)

	b := 0
	quoted := false
	for i, v := range env {
		if v == '"' {
			quoted = !quoted
		}
		if !quoted && v == ' ' {
			pair := strings.SplitN(strings.Replace(env[b:i], "\"", "", 2), "=", 2)
			if len(pair) == 2 {
				c[pair[0]] = pair[1]
			}
			b = i + 1
		}
	}
	pair := strings.SplitN(strings.Replace(env[b:], "\"", "", 2), "=", 2)
	if len(pair) == 2 {
		c[pair[0]] = pair[1]
	}
	return

}
