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
	global := parseEnv(strings.Join(t.Env.Global, " "))

	// Parse Matrix.Include.
	for _, v := range t.Matrix.Include {
		version, c, err := parseMatrixGo(v)
		if err != nil {
			return nil, err
		}
		res[version] = append(res[version], global.Copy().Merge(c))
	}

	// Parse general environment.
	if len(t.Go) == 0 && len(res) == 0 {
		t.Go = []string{"any"}
	}
	for _, version := range t.Go {

		if len(t.Env.Matrix) == 0 {
			// If there is no matrix configuration, use only global configuration.
			res[version] = append(res[version], global)
		} else {
			// Look up each matrix case, and merge sprcific configuration to the
			// global one.
			for _, m := range t.Env.Matrix {
				c := parseEnv(m)
				res[version] = append(res[version], global.Copy().Merge(c))
			}
		}

	}

	// Parse Matrix.Exclude.
	for _, v := range t.Matrix.Exclude {
		version, exclude, err := parseMatrixGo(v)
		if err != nil {
			return nil, err
		}
		if set, ok := res[version]; ok {
			var new []TestCase
			for _, c := range set {
				if !c.Match(exclude) {
					new = append(new, c)
				}
			}
			res[version] = new
		}
	}

	return

}

// parseMatrixGo parses an interface representing an entry of env.matrix;
// and returns the version and test case the interface specifies for golang.
func parseMatrixGo(v interface{}) (version string, c TestCase, err error) {

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
	c = parseEnv(variables)

	return

}
