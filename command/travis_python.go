//
// command/travis_python.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

// argumentSetPython returns a set of arguments to run entrypoint based on a build
// matrix for python projects.
func (t *Travis) argumentSetPython() [][]string {

	if len(t.Matrix.Include) != 0 {
		res := make([][]string, len(t.Matrix.Include))
		for i, v := range t.Matrix.Include {
			version, key, value := parsePythonCase(v)
			res[i] = []string{version, key, value}
		}
		return res
	}

	exclude := make(map[string]struct{})
	for _, v := range t.Matrix.Exclude {
		version, key, value := parsePythonCase(v)
		exclude[makeSetKey(version, key, value)] = struct{}{}
	}

	if len(t.Python) == 0 {

		if len(t.Env) == 0 {
			return [][]string{[]string{"2.7"}}
		}

		res := make([][]string, len(t.Env))
		for i, env := range t.Env {
			key, value := parseEnv(env)
			res[i] = []string{"2.7", key, value}
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
			key, value := parseEnv(env)
			if _, exist := exclude[makeSetKey(ver, key, value)]; exist {
				continue
			}
			res = append(res, []string{ver, key, value})
		}
	}
	return res

}

// parsePythonCase parses a given item v in the exclude list, and returns
// a tuple of version, key of an environment variable, and its value.
func parsePythonCase(v interface{}) (version string, key string, value string) {
	m, _ := v.(map[interface{}]interface{})
	version, _ = m["python"].(string)
	env, _ := m["env"].(string)
	key, value = parseEnv(env)
	return
}
