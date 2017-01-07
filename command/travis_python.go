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

import "fmt"

// argumentSetPython returns a set of arguments to run entrypoint based on a build
// matrix for python projects.
func (t *Travis) argumentSetPython() (res []Arguments, err error) {

	if len(t.Matrix.Include) != 0 {
		res = make([]Arguments, len(t.Matrix.Include))
		for i, v := range t.Matrix.Include {
			args, err := newPythonArguments(v)
			if err != nil {
				return nil, err
			}
			res[i] = args
		}
		return
	}

	exclude := make(map[string]struct{})
	for _, v := range t.Matrix.Exclude {
		args, err := newPythonArguments(v)
		if err != nil {
			return nil, err
		}
		exclude[args.String()] = struct{}{}
	}

	if len(t.Python) == 0 {

		if len(t.Env) == 0 {
			res = []Arguments{
				Arguments{
					Version: "2.7",
				},
			}
			return
		}

		res = make([]Arguments, len(t.Env))
		for i, env := range t.Env {
			res[i] = Arguments{
				Version: "2.7",
				Env:     env,
			}
		}
		return

	}

	if len(t.Env) == 0 {
		res = make([]Arguments, len(t.Python))
		for i, ver := range t.Python {
			res[i] = Arguments{
				Version: ver,
			}
		}
		return
	}

	for _, ver := range t.Python {
		for _, env := range t.Env {
			args := Arguments{
				Version: ver,
				Env:     env,
			}
			if _, exist := exclude[args.String()]; !exist {
				res = append(res, args)
			}
		}
	}

	return

}

// newPythonArguments parses a given item v in an include/exclude list.
// v must be castable to map[interface{}]interface{}.
func newPythonArguments(v interface{}) (res Arguments, err error) {

	m, ok := v.(map[interface{}]interface{})
	if !ok {
		err = fmt.Errorf("Given item is broken.")
		return
	}

	res.Version, ok = m["python"].(string)
	if !ok {
		err = fmt.Errorf("Python version of the given item is broken.")
		return
	}

	res.Env, ok = m["env"].(string)
	if !ok {
		err = fmt.Errorf("Env of the given item is broken.")
		return
	}

	return

}
