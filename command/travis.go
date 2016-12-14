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
	// List of commands used to install packages.
	Install []string `yaml:"install,omitempty"`
	// List of commands run before main scripts.
	BeforeScript []string `yaml:"before_script,omitempty"`
	// List of scripts.
	Script []string `yaml:"script,omitempty"`
	// List of environment variables.
	Env []string `yaml:"env,omitempty"`
	// List of python versions. (used only in python)
	Python []string `yaml:"python,omitempty"`
	// Configuration for matrix build.
	Matrix Matrix `yaml:"matrix,omitempty"`
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

	if t.Language == "python" {

		if len(t.Matrix.Include) != 0 {

			res := make([][]string, len(t.Matrix.Include))
			for i, v := range t.Matrix.Include {
				m, _ := v.(map[interface{}]interface{})
				version, _ := m["python"].(string)
				env, _ := m["env"].(string)
				res[i] = append([]string{version}, strings.SplitN(env, "=", 2)...)
			}
			return res

		}

		if len(t.Python) == 0 {

			if len(t.Env) == 0 {
				return [][]string{[]string{"2.7"}}
			}

			res := make([][]string, len(t.Env))
			for i, env := range t.Env {
				res[i] = append([]string{"2.7"}, strings.SplitN(env, "=", 2)...)
			}
			return res

		}

		if len(t.Env) == 0 {
			res := make([][]string, len(t.Python))
			for i, ver := range t.Python {
				res[i] = []string{ver}
			}
			return res
		}

		res := [][]string{}
		for _, ver := range t.Python {
			for _, env := range t.Env {
				res = append(res, append([]string{ver}, strings.SplitN(env, "=", 2)...))
			}
		}
		return res

	}

	return [][]string{{""}}

}
