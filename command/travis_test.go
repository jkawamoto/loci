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
)

func TestParseEnvWithNoValues(t *testing.T) {

	travis := &Travis{}
	travis.parseEnv()
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
	travis := &Travis{
		RawEnv: globals,
	}
	travis.parseEnv()
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
	travis := &Travis{
		RawEnv: matrix,
	}
	travis.parseEnv()
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
	travis := &Travis{
		RawEnv: globalEnv{
			Global: globals,
		},
	}
	travis.parseEnv()
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
	travis := &Travis{
		RawEnv: matrixEnv{
			Matrix: matrix,
		},
	}
	travis.parseEnv()
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
	travis := &Travis{
		RawEnv: fullEnv{
			Global: globals,
			Matrix: matrix,
		},
	}
	travis.parseEnv()
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

func TestNewTravis(t *testing.T) {

	var err error
	test := `language: go
go:
- 1.7.4
env:
- TEST_ENV=true
- TEST_ENV=false
install:
- make get-deps
script:
- echo $TEST_ENV
- "if [[ $TEST_ENV = true ]]; then make test; fi"
`

	temp := os.TempDir()
	target := path.Join(temp, ".travis.yml")
	t.Logf("Creating a Travis configuration file: %s", target)

	if err = ioutil.WriteFile(target, []byte(test), 0644); err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(target)

	travis, err := NewTravis(target)
	if err != nil {
		t.Fatal(err.Error())
	}
	if travis.RawEnv == nil {
		t.Error("Mapping Env attribute is failed:", travis.Env)
	}

}
