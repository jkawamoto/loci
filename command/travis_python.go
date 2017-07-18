//
// command/travis_python.go
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
	"strings"

	"github.com/ttacon/chalk"
)

// const (
// 	// PythonNightlyVersion defines a python version used as the nightly version.
// 	PythonNightlyVersion = "3.7-dev"
// )

var (
	// ErrUnknownPythonVersion is returned when the given python version is not
	// supported.
	ErrUnknownPythonVersion = fmt.Errorf("Given Python version is not supported")
)

// argumentSetPython returns a set of arguments to run entrypoint based on a build
// matrix for python projects.
func (t *Travis) argumentSetPython(logger io.Writer) (res TestCaseSet, err error) {

	res = make(TestCaseSet)
	global := parseEnv(strings.Join(t.Env.Global, " "))

	// Parse Matrix.Include.
	for _, v := range t.Matrix.Include {
		version, c, err := parseMatrixPython(v)
		if err == ErrUnknownPythonVersion {
			fmt.Fprintf(logger, chalk.Yellow.Color("Python version %v is not supported\n"), version)
			continue
		} else if err != nil {
			return nil, err
		}
		res[version] = append(res[version], global.Copy().Merge(c))
	}

	if len(t.Python) == 0 && len(res) == 0 {
		t.Python = []string{"2.7"}
	}
	for _, version := range t.Python {

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
		version, exclude, err := parseMatrixPython(v)
		if err == ErrUnknownPythonVersion {
			fmt.Fprintf(logger, chalk.Yellow.Color("Python version %v is not supported\n"), version)
			continue
		} else if err != nil {
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

// parseMatrixPython parses a given item v in an include/exclude list.
// v must be castable to map[interface{}]interface{}.
func parseMatrixPython(v interface{}) (version string, c TestCase, err error) {

	m, ok := v.(map[interface{}]interface{})
	if !ok {
		err = fmt.Errorf("Given item is broken.")
		return
	}

	if _, exist := m["python"]; !exist {
		version = "empty"
		err = ErrUnknownPythonVersion
		return
	}
	version = fmt.Sprint(m["python"])
	if version == "nightly" {
		// version = PythonNightlyVersion
		err = ErrUnknownPythonVersion
		return
	}
	variables, ok := m["env"].(string)
	if !ok {
		err = fmt.Errorf("Env of the given item is broken.")
		return
	}
	c = parseEnv(variables)

	return

}
