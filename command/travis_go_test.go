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

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"gopkg.in/yaml.v2"
)

// GoCase defines a case of matrix evaluation for go projects.
type GoCase struct {
	Go  string `yaml:"go"`
	Env string `yaml:"env"`
}

func TestGoMatrixInclude(t *testing.T) {

	var err error
	temp := os.TempDir()
	target := path.Join(temp, "sample.yml")

	t.Logf("Creating a configuration file: %s", target)
	sample, err := yaml.Marshal(&Travis{
		Language: "go",
		Matrix: Matrix{
			Include: []interface{}{
				GoCase{
					Go:  "1.6",
					Env: "FOO=BAR",
				}, GoCase{
					Go:  "1.7",
					Env: "FOO=FUGA",
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

	res, err := travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Fatal("Generated arguments are wrong:", res)
	}
	if res[0].Version != "1.6" || res[1].Version != "1.7" {
		t.Error("Version is wrong:", res)
	}
	if res[0].Env != "FOO=BAR" || res[1].Env != "FOO=FUGA" {
		t.Error("Env has wrong values:", res)
	}

}

func TestGoMatrixExclude(t *testing.T) {

	var err error
	temp := os.TempDir()
	target := path.Join(temp, "sample.yml")

	t.Logf("Creating a configuration file: %s", target)
	sample, err := yaml.Marshal(&Travis{
		Language: "go",
		Go:       []string{"1.6", "1.7"},
		Env:      []string{"FOO=BAR", "FOO=FUGA"},
		Matrix: Matrix{
			Exclude: []interface{}{
				GoCase{
					Go:  "1.7",
					Env: "FOO=BAR",
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

	res, err := travis.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 3 {
		t.Fatal("Generated arguments are wrong:", res)
	}
	if res[0].Version != "1.6" || res[1].Version != "1.6" || res[2].Version != "1.7" {
		t.Error("Version is wrong:", res)
	}
	if res[0].Env != "FOO=BAR" || res[1].Env != "FOO=FUGA" || res[2].Env != "FOO=FUGA" {
		t.Error("Env has wrong values:", res)
	}

}

func TestGoArgumentSet(t *testing.T) {

	var v *Travis
	var res []Arguments
	var err error

	v = &Travis{
		Language: "go",
	}

	res, err = v.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0].Version != "any" {
		t.Error("Version is wrong:", res[0].Version)
	}
	if res[0].Env != "" {
		t.Error("Env has wrong values:", res[0].Env)
	}

	v = &Travis{
		Language: "go",
		Env:      []string{"FOO=BAR"},
	}
	res, err = v.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 1 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0].Version != "any" {
		t.Error("Version is wrong:", res)
	}
	if res[0].Env != "FOO=BAR" {
		t.Error("Env has wrong values:", res)
	}

	v = &Travis{
		Language: "go",
		Go:       []string{"1.6", "1.7"},
	}
	res, err = v.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 2 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0].Version != "1.6" || res[1].Version != "1.7" {
		t.Error("Version is wrong:", res)
	}
	if res[0].Env != "" || res[1].Env != "" {
		t.Error("Env has wrong values:", res)
	}

	v = &Travis{
		Language: "go",
		Go:       []string{"1.6", "1.7"},
		Env:      []string{"FOO=BAR", "FOO=FUGA"},
	}
	res, err = v.ArgumentSet()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Arguments:", res)
	if len(res) != 4 {
		t.Error("Generated arguments are wrong:", res)
	}
	if res[0].Version != "1.6" || res[1].Version != "1.6" || res[2].Version != "1.7" || res[3].Version != "1.7" {
		t.Error("Version is wrong:", res)
	}
	if res[0].Env != "FOO=BAR" {
		t.Error("Env has wrong values:", res)
	}
	if res[1].Env != "FOO=FUGA" {
		t.Error("Env has wrong values:", res)
	}
	if res[2].Env != "FOO=BAR" {
		t.Error("Env has wrong values:", res)
	}
	if res[3].Env != "FOO=FUGA" {
		t.Error("Env has wrong values:", res)
	}

}
