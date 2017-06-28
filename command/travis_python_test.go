//
// command/travis_python_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//
package command

import (
	"io/ioutil"
	"testing"
)

// PythonCase defines a case of matrix evaluation for python projects.
type PythonCase struct {
	Python string `yaml:"python"`
	Env    string `yaml:"env"`
}

func TestPythonMatrixInclude(t *testing.T) {

	var err error
	travis, err := storeAndLoadTravis(&Travis{
		Language: "python",
		Matrix: Matrix{
			Include: []interface{}{
				PythonCase{
					Python: "2.7",
					Env:    "FOO=bar",
				}, PythonCase{
					Python: "3.5",
					Env:    "FOO=fuga",
				},
			},
		},
	})
	if err != nil {
		t.Error(err.Error())
	}
	if len(travis.Matrix.Include) != 2 {
		t.Error("Size of items in matrix.include is wrong:", travis.Matrix.Include)
	}

	res, err := travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Fatal("Generated arguments are wrong:", res)
	}
	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 1 || set[0][0] != "FOO=bar" {
		t.Error("Env has wrong values:", res)
	}
	if set, ok := res["3.5"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 1 || set[0][0] != "FOO=fuga" {
		t.Error("Env has wrong values:", res)
	}

}

func TestPythonMatrixExclude(t *testing.T) {

	var err error
	travis, err := NewTravis([]byte(`language: "python"
python:
  - 2.7
  - 3.5
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
matrix:
  exclude:
    - python: 3.5
      env: FOO=bar BAR=foo
`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Matrix.Exclude) != 1 {
		t.Error("Size of items in matrix.exclude is wrong:", travis.Matrix.Exclude)
	}

	res, err := travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Fatal("Generated arguments are wrong:", res)
	}

	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 2 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 2 || set[0][0] != "FOO=foo" || set[0][1] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
		if len(set[1]) != 2 || set[1][0] != "FOO=bar" || set[1][1] != "BAR=foo" {
			t.Error("Env has wrong values:", res)
		}
	}
	if set, ok := res["3.5"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 2 || set[0][0] != "FOO=foo" || set[0][1] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
	}

}

// TestPythonArgumentSet tests ArgumentSet method returns correct argument sets.
func TestPythonArgumentSet(t *testing.T) {

	var (
		travis *Travis
		res    TestCaseSet
		err    error
	)

	travis, err = storeAndLoadTravis(&Travis{
		Language: "python",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 0 {
		t.Error("Env has wrong values:", res)
	}

	travis, err = NewTravis([]byte(`language: "python"
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
`))

	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 2 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 2 || set[0][0] != "FOO=foo" || set[0][1] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
		if len(set[1]) != 2 || set[1][0] != "FOO=bar" || set[1][1] != "BAR=foo" {
			t.Error("Env has wrong values:", res)
		}
	}

	travis, err = storeAndLoadTravis(&Travis{
		Language: "python",
		Python:   []string{"2.7", "3.5"},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 0 {
		t.Error("Env has wrong values:", res)
	}
	if set, ok := res["3.5"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 0 {
		t.Error("Env has wrong values:", res)
	}

	travis, err = NewTravis([]byte(`language: "python"
python:
  - 2.7
  - 3.5
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
`))
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 2 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 2 || set[0][0] != "FOO=foo" || set[0][1] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
		if len(set[1]) != 2 || set[1][0] != "FOO=bar" || set[1][1] != "BAR=foo" {
			t.Error("Env has wrong values:", res)
		}
	}
	if set, ok := res["3.5"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 2 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 2 || set[0][0] != "FOO=foo" || set[0][1] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
		if len(set[1]) != 2 || set[1][0] != "FOO=bar" || set[1][1] != "BAR=foo" {
			t.Error("Env has wrong values:", res)
		}
	}

}

func TestPythonArgumentSetWithFullDescriptions(t *testing.T) {

	travis, err := storeAndLoadTravis(&Travis{
		Language: "python",
		Python:   []string{"2.7", "3.5"},
		Env: Env{
			Global: []string{"GLOBAL=global"},
			Matrix: []string{"FOO=foo BAR=bar", "FOO=bar BAR=foo"},
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["2.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 2 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 3 || set[0][0] != "GLOBAL=global" || set[0][1] != "FOO=foo" || set[0][2] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
		if len(set[1]) != 3 || set[1][0] != "GLOBAL=global" || set[1][1] != "FOO=bar" || set[1][2] != "BAR=foo" {
			t.Error("Env has wrong values:", res)
		}
	}
	if set, ok := res["3.5"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 2 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 3 || set[0][0] != "GLOBAL=global" || set[0][1] != "FOO=foo" || set[0][2] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
		if len(set[1]) != 3 || set[1][0] != "GLOBAL=global" || set[1][1] != "FOO=bar" || set[1][2] != "BAR=foo" {
			t.Error("Env has wrong values:", res)
		}
	}

}

func TestPythonUnknownArgumentSet(t *testing.T) {

	var err error
	// The following configuration is copied from matplotlib.
	travis, err := NewTravis([]byte(`language: "python"
matrix:
  include:
    - python: "nightly"
      env: PRE=--pre
    - os: osx
      osx_image: xcode7.3
      language: generic
`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Matrix.Include) != 2 {
		t.Error("Size of items in matrix.include is wrong:", travis.Matrix.Exclude)
	}

	res, err := travis.ArgumentSet(ioutil.Discard)
	if err != nil {
		t.Error(err.Error())
	}

	if len(res) != 0 {
		t.Fatal("Generated arguments are wrong:", res)
	}

}
