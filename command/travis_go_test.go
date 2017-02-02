//
// command/travis_go_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//
package command

import "testing"

// GoCase defines a case of matrix evaluation for go projects.
type GoCase struct {
	Go  string `yaml:"go"`
	Env string `yaml:"env"`
}

func TestGoMatrixInclude(t *testing.T) {

	var err error
	travis, err := storeAndLoadTravis(&Travis{
		Language: "go",
		Matrix: Matrix{
			Include: []interface{}{
				GoCase{
					Go:  "1.6",
					Env: "FOO=bar",
				}, GoCase{
					Go:  "1.7",
					Env: "FOO=fuga",
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

	res, err := travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Fatal("Generated arguments are wrong:", res)
	}

	if set, ok := res["1.6"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 1 || set[0][0] != "FOO=bar" {
		t.Error("Env has wrong values:", res)
	}

	if set, ok := res["1.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 1 || set[0][0] != "FOO=fuga" {
		t.Error("Env has wrong values:", res)
	}

}

func TestGoMatrixExclude(t *testing.T) {

	var err error
	travis, err := NewTravis([]byte(`language: "go"
go:
  - 1.6
  - 1.7
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
matrix:
  exclude:
    - go: 1.7
      env: FOO=bar BAR=foo
`))

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Matrix.Exclude) != 1 {
		t.Error("Size of items in matrix.include is wrong:", travis.Matrix.Exclude)
	}

	res, err := travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Fatal("Generated arguments are wrong:", res)
	}

	if set, ok := res["1.6"]; !ok {
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
	if set, ok := res["1.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 {
		t.Error("Env has wrong values:", res)
	} else {
		if len(set[0]) != 2 || set[0][0] != "FOO=foo" || set[0][1] != "BAR=bar" {
			t.Error("Env has wrong values:", res)
		}
	}

}

func TestGoArgumentSet(t *testing.T) {

	var (
		travis *Travis
		res    TestCaseSet
		err    error
	)

	travis, err = storeAndLoadTravis(&Travis{
		Language: "go",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["any"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 0 {
		t.Error("Env has wrong values:", set)
	}

	travis, err = NewTravis([]byte(`language: "go"
env:
  - FOO=bar
`))
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["any"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 1 || set[0][0] != "FOO=bar" {
		t.Error("Env has wrong values:", res)
	}

	travis, err = NewTravis([]byte(`language: "go"
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
`))
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["any"]; !ok {
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
		Language: "go",
		Go:       []string{"1.6", "1.7"},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["1.6"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 0 {
		t.Error("Env has wrong values:", res)
	}
	if set, ok := res["1.7"]; !ok {
		t.Error("Version is wrong:", res)
	} else if len(set) != 1 || len(set[0]) != 0 {
		t.Error("Env has wrong values:", res)
	}

	travis, err = NewTravis([]byte(`language: "go"
go:
  - 1.6
  - 1.7
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
`))
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err = travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["1.6"]; !ok {
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
	if set, ok := res["1.7"]; !ok {
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

func TestGoArgumentSetWithFullDescriptions(t *testing.T) {

	travis, err := storeAndLoadTravis(&Travis{
		Language: "go",
		Go:       []string{"1.6", "1.7"},
		Env: Env{
			Global: []string{"GLOBAL=global"},
			Matrix: []string{"FOO=foo BAR=bar", "FOO=bar BAR=foo"},
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if set, ok := res["1.6"]; !ok {
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
	if set, ok := res["1.7"]; !ok {
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
