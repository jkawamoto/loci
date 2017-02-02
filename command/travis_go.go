//
// command/travis_go.go
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
)

// argumentSetGo returns a set of arguments to run entrypoint based on a build
// matrix for Go projects.
func (t *Travis) argumentSetGo() (res TestCaseSet, err error) {

	res = make(TestCaseSet)

	// Parse Matrix.Include.
	for _, v := range t.Matrix.Include {
		version, env, err := parseMatrixGo(v)
		if err != nil {
			return nil, err
		}
		if res[version] == nil {
			res[version] = [][]string{env}
		} else {
			res[version] = append(res[version], env)
		}

	}

	if len(t.Go) == 0 && len(res) == 0 {
		t.Go = []string{"any"}
	}
	for _, version := range t.Go {

		if len(t.Env.Matrix) == 0 {
			res[version] = [][]string{t.Env.Global}

		} else {
			res[version] = [][]string{}
			for _, m := range t.Env.Matrix {
				envs := t.Env.Global
				for _, pair := range strings.Split(strings.TrimSpace(m), " ") {
					envs = append(envs, pair)
				}
				res[version] = append(res[version], envs)
			}

		}

	}

	// Parse Matrix.Exclude.
	for _, v := range t.Matrix.Exclude {
		version, env, err := parseMatrixGo(v)
		if err != nil {
			return nil, err
		}
		excludes := make(map[string]struct{})
		for _, e := range env {
			excludes[e] = struct{}{}
		}
		if set, ok := res[version]; ok {
			res[version] = [][]string{}
			for _, item := range set {
				skip := false
				for _, pair := range item {
					if _, exist := excludes[pair]; exist {
						skip = true
						break
					}
				}
				if !skip {
					res[version] = append(res[version], item)
				}
			}
		}
	}

	return

}

func parseMatrixGo(v interface{}) (version string, env []string, err error) {

	m, ok := v.(map[interface{}]interface{})
	if !ok {
		err = fmt.Errorf("Given item is broken.")
		return
	}

	version = fmt.Sprint(m["go"])

	variables, ok := m["env"].(string)
	if !ok {
		err = fmt.Errorf("Env of the given item is broken.")
		return
	}
	env = strings.Split(strings.TrimSpace(variables), " ")

	return

}
