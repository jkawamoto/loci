//
// command/travis_python_test.go
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
	"path"
	"testing"

	"gopkg.in/yaml.v2"
)

// PythonCase defines a case of matrix evaluation for python projects.
type PythonCase struct {
	Python string `yaml:"python"`
	Env    string `yaml:"env"`
}

func TestPythonMatrixInclude(t *testing.T) {

	var err error
	temp := os.TempDir()
	target := path.Join(temp, "sample.yml")

	t.Logf("Creating a configuration file: %s", target)
	sample, err := yaml.Marshal(&Travis{
		Language: "python",
		Matrix: Matrix{
			Include: []interface{}{
				PythonCase{
					Python: "2.7",
					Env:    "FOO=BAR",
				}, PythonCase{
					Python: "3.5",
					Env:    "FOO=FUGA",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = ioutil.WriteFile(target, sample, 0644); err != nil {
		t.Fatal(err.Error())
	}

	travis, err := NewTravis(target)
	if err != nil {
		t.Error(err.Error())
	}
	if len(travis.Matrix.Include) != 2 {
		t.Error("Size of items in matrix.include is wrong:", travis.Matrix.Include)
	}

	res := travis.ArgumentSet()
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Fatal("Generated arguments are wrong:", res)
	}
	if res[0][0] != "2.7" || res[1][0] != "3.5" {
		t.Error("Version is wrong:", res)
	}
	if len(res[0]) != 3 || res[0][1] != "FOO" || res[0][2] != "BAR" {
		t.Error("Env has wrong values:", res)
	}
	if len(res[1]) != 3 || res[1][1] != "FOO" || res[1][2] != "FUGA" {
		t.Error("Env has wrong values:", res)
	}

}

func TestPythonMatrixExclude(t *testing.T) {

	var err error
	temp := os.TempDir()
	target := path.Join(temp, "sample.yml")

	t.Logf("Creating a configuration file: %s", target)
	sample, err := yaml.Marshal(&Travis{
		Language: "python",
		Python:   []string{"2.7", "3.5"},
		Env:      []string{"FOO=BAR", "FOO=FUGA"},
		Matrix: Matrix{
			Exclude: []interface{}{
				PythonCase{
					Python: "3.5",
					Env:    "FOO=BAR",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = ioutil.WriteFile(target, sample, 0644); err != nil {
		t.Fatal(err.Error())
	}

	travis, err := NewTravis(target)
	if err != nil {
		t.Error(err.Error())
	}
	if len(travis.Matrix.Exclude) != 1 {
		t.Error("Size of items in matrix.include is wrong:", travis.Matrix.Exclude)
	}

	res := travis.ArgumentSet()
	t.Log("Arguments:", res)
	if len(res) != 3 {
		t.Fatal("Generated arguments are wrong:", res)
	}
	if res[0][0] != "2.7" || res[1][0] != "2.7" || res[2][0] != "3.5" {
		t.Error("Version is wrong:", res)
	}
	if len(res[0]) != 3 || res[0][1] != "FOO" || res[0][2] != "BAR" {
		t.Error("Env has wrong values:", res)
	}
	if len(res[1]) != 3 || res[1][1] != "FOO" || res[1][2] != "FUGA" {
		t.Error("Env has wrong values:", res)
	}
	if len(res[2]) != 3 || res[2][1] != "FOO" || res[2][2] != "FUGA" {
		t.Error("Env has wrong values:", res)
	}

}

// TestPythonArgumentSet tests ArgumentSet method returns correct argument sets.
func TestPythonArgumentSet(t *testing.T) {

	var v *Travis
	var res [][]string

	v = &Travis{
		Language: "python",
	}
	res = v.ArgumentSet()
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0][0] != "2.7" {
		t.Error("Version is wrong:", res[0][0])
	}
	if len(res[0]) != 1 {
		t.Error("Env has wrong values:", res[0])
	}

	v = &Travis{
		Language: "python",
		Env:      []string{"FOO=BAR"},
	}
	res = v.ArgumentSet()
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0][0] != "2.7" {
		t.Error("Version is wrong:", res)
	}
	if len(res[0]) != 3 || res[0][1] != "FOO" || res[0][2] != "BAR" {
		t.Error("Env has wrong values:", res)
	}

	v = &Travis{
		Language: "python",
		Python:   []string{"2.7", "3.5"},
	}
	res = v.ArgumentSet()
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0][0] != "2.7" || res[1][0] != "3.5" {
		t.Error("Version is wrong:", res)
	}
	if len(res[0]) != 1 || len(res[1]) != 1 {
		t.Error("Env has wrong values:", res)
	}

	v = &Travis{
		Language: "python",
		Python:   []string{"2.7", "3.5"},
		Env:      []string{"FOO=BAR", "FOO=FUGA"},
	}
	res = v.ArgumentSet()
	t.Log("Arguments:", res)
	if len(res) != 4 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0][0] != "2.7" || res[1][0] != "2.7" || res[2][0] != "3.5" || res[3][0] != "3.5" {
		t.Error("Version is wrong:", res)
	}
	if len(res[0]) != 3 || res[0][1] != "FOO" || res[0][2] != "BAR" {
		t.Error("Env has wrong values:", res)
	}
	if len(res[1]) != 3 || res[1][1] != "FOO" || res[1][2] != "FUGA" {
		t.Error("Env has wrong values:", res)
	}
	if len(res[2]) != 3 || res[2][1] != "FOO" || res[2][2] != "BAR" {
		t.Error("Env has wrong values:", res)
	}
	if len(res[3]) != 3 || res[3][1] != "FOO" || res[3][2] != "FUGA" {
		t.Error("Env has wrong values:", res)
	}

}
