//
// command/travis_test.go
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

	yaml "gopkg.in/yaml.v2"
)

func storeAndLoadTravis(src *Travis) (res *Travis, err error) {
	temp := os.TempDir()
	target := path.Join(temp, "sample.yml")
	data, err := yaml.Marshal(src)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(target, data, 0644); err != nil {
		return
	}
	defer os.Remove(target)
	return NewTravis(target)
}

func TestParseEnvWithNoValues(t *testing.T) {

	travis, err := storeAndLoadTravis(&Travis{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Env.Global) != 0 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	if len(travis.Env.Matrix) != 0 {
		t.Fatal("The number of matrix variables is wrong:", travis.Env.Matrix)
	}

}

func TestParseEnvWithGlobalsList(t *testing.T) {

	globals := []string{
		"DB=postgres",
		"SH=bash",
		"PACKAGE_VERSION=\"1.0.*\"",
	}
	travis, err := storeAndLoadTravis(&Travis{
		RawEnv: globals,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Env.Global) != 3 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	if len(travis.Env.Matrix) != 0 {
		t.Fatal("The number of matrix variables is wrong:", travis.Env.Matrix)
	}
	for i, v := range globals {
		if travis.Env.Global[i] != v {
			t.Error("A global variable is not match:", travis.Env.Global)
		}
	}

}

func TestParseEnvWithMultipleVariablesList(t *testing.T) {

	matrix := []string{
		"FOO=foo BAR=bar",
		"FOO=bar BAR=foo",
	}
	travis, err := storeAndLoadTravis(&Travis{
		RawEnv: matrix,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Env.Global) != 0 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	if len(travis.Env.Matrix) != 2 {
		t.Fatal("The number of matrix variables is wrong:", travis.Env.Matrix)
	}
	for i, v := range matrix {
		if travis.Env.Matrix[i] != v {
			t.Error("Matrix variables are not match:", travis.Env.Matrix)
		}
	}

}

func TestParseEnvWithSpecificGlobals(t *testing.T) {

	globals := []string{
		"DB=postgres",
		"SH=bash",
		"PACKAGE_VERSION=\"1.0.*\"",
	}
	travis, err := storeAndLoadTravis(&Travis{
		RawEnv: struct {
			Global []string
		}{
			Global: globals,
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Env.Global) != 3 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	if len(travis.Env.Matrix) != 0 {
		t.Fatal("The number of matrix variables is wrong:", travis.Env.Matrix)
	}
	for i, v := range globals {
		if travis.Env.Global[i] != v {
			t.Error("A global variable is not match:", travis.Env.Global)
		}
	}

}

func TestParseEnvWithSpecificMatrixVariables(t *testing.T) {

	matrix := []string{
		"FOO=foo BAR=bar",
		"FOO=bar BAR=foo",
	}
	travis, err := storeAndLoadTravis(&Travis{
		RawEnv: struct {
			Matrix []string
		}{
			Matrix: matrix,
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Env.Global) != 0 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	if len(travis.Env.Matrix) != 2 {
		t.Fatal("The number of matrix variables is wrong:", travis.Env.Matrix)
	}
	for i, v := range matrix {
		if travis.Env.Matrix[i] != v {
			t.Error("Matrix variables are not match:", travis.Env.Matrix)
		}
	}

}

func TestParseEnv(t *testing.T) {

	globals := []string{
		"DB=postgres",
		"SH=bash",
		"PACKAGE_VERSION=\"1.0.*\"",
	}
	matrix := []string{
		"FOO=foo BAR=bar",
		"FOO=bar BAR=foo",
	}
	travis, err := storeAndLoadTravis(&Travis{
		RawEnv: struct {
			Global []string
			Matrix []string
		}{
			Global: globals,
			Matrix: matrix,
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Env.Global) != 3 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	if len(travis.Env.Matrix) != 2 {
		t.Fatal("The number of matrix variables is wrong:", travis.Env.Matrix)
	}
	for i, v := range globals {
		if travis.Env.Global[i] != v {
			t.Error("A global variable is not match:", travis.Env.Global)
		}
	}
	for i, v := range matrix {
		if travis.Env.Matrix[i] != v {
			t.Error("Matrix variables are not match:", travis.Env.Matrix)
		}
	}

}
