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

func TestTestCaseSlice(t *testing.T) {

	c := TestCase{
		"FOO": "foo",
		"BAR": "bar",
	}

	var foo, bar bool
	for _, s := range c.Slice() {
		switch s {
		case "FOO=foo":
			foo = true
		case "BAR=bar":
			bar = true
		default:
			t.Error("The slice converted from a test case has a wrong string:", s)
		}
	}
	if !foo || !bar {
		t.Errorf("The slice converted from a test case has missing variables: FOO=%v, BAR=%v", foo, bar)
	}

}

func TestTestCaseCopy(t *testing.T) {

	a := TestCase{
		"FOO": "foo",
		"BAR": "bar",
	}
	b := a.Copy()
	a["FOO"] = "piyo"

	if b["FOO"] != "foo" || b["BAR"] != "bar" {
		t.Error("Copied test case is not a hard copy", b)
	}

}

func TestTestCaseMerge(t *testing.T) {

	a := TestCase{
		"FOO": "foo",
		"BAR": "bar",
	}
	b := TestCase{
		"BAR":  "piyo",
		"FUGA": "fuga",
	}
	a.Merge(b)

	var foo, bar, fuga bool
	for key, value := range a {
		switch {
		case key == "FOO" && value == "foo":
			foo = true
		case key == "BAR" && value == "piyo":
			bar = true
		case key == "FUGA" && value == "fuga":
			fuga = true
		default:
			t.Errorf("Merged test case has a wrong pair: %v=%v", key, value)
		}
	}
	if !foo || !bar || !fuga {
		t.Errorf("Merged test case has missing pairs: FOO=%v, BAR=%v, FUGA=%v", foo, bar, fuga)
	}

}

func TestTestCaseMatch(t *testing.T) {

	a := TestCase{
		"FOO": "foo",
		"BAR": "bar",
	}
	if !a.Match(a.Copy()) {
		t.Error("Match returns false for a copied test case")
	}
	b := a.Copy()
	a["FOO"] = "piyo"
	if a.Match(b) {
		t.Error("Match returns true for a different test case")
	}

}

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
	return NewTravisFromFile(target)
}

func TestParseBeforeInstallWithNoValues(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeInstall) != 0 {
		t.Error("BeforeInstall is wrong:", travis.BeforeInstall)
	}
}

func TestParseBeforeInstallWithString(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
before_install: install 1`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeInstall) != 1 || travis.BeforeInstall[0] != "install 1" {
		t.Error("BeforeInstall is wrong:", travis.BeforeInstall)
	}
}

func TestParseBeforeInstallWithList(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
before_install:
  - install 1`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeInstall) != 1 || travis.BeforeInstall[0] != "install 1" {
		t.Error("BeforeInstall is wrong:", travis.BeforeInstall)
	}
}

func TestParseInstallWithNoValues(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Install) != 0 {
		t.Error("Install is wrong:", travis.Install)
	}
}

func TestParseInstallWithString(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
install: install 1`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Install) != 1 || travis.Install[0] != "install 1" {
		t.Error("Install is wrong:", travis.Install)
	}
}

func TestParseInstallWithList(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
install:
  - install 1`))

	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Install) != 1 || travis.Install[0] != "install 1" {
		t.Error("Install is wrong:", travis.Install)
	}
}

func TestParseBeforeScriptWithNoValues(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeScript) != 0 {
		t.Error("BeforeScript is wrong:", travis.BeforeScript)
	}
}

func TestParseBeforeScriptWithString(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
before_script: python setup.py test`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeScript) != 1 || travis.BeforeScript[0] != "python setup.py test" {
		t.Error("BeforeScript is wrong:", travis.BeforeScript)
	}
}

func TestParseBeforeScriptWithList(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
before_script:
  - python setup.py test`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeScript) != 1 || travis.BeforeScript[0] != "python setup.py test" {
		t.Error("BeforeScript is wrong:", travis.BeforeScript)
	}
}

func TestParseScriptWithNoValues(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Script) != 0 {
		t.Error("Script is wrong:", travis.Script)
	}
}

func TestParseScriptWithString(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
script: python setup.py test`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Script) != 1 || travis.Script[0] != "python setup.py test" {
		t.Error("Script is wrong:", travis.Script)
	}
}

func TestParseScriptWithList(t *testing.T) {
	travis, err := NewTravis([]byte(`language: ""
script:
  - python setup.py test`))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Script) != 1 || travis.Script[0] != "python setup.py test" {
		t.Error("Script is wrong:", travis.Script)
	}
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
	travis, err := NewTravis([]byte(`language: ""
env:
  - DB=postgres
  - SH=bash
  - PACKAGE_VERSION="1.0.*"
`))
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
	travis, err := NewTravis([]byte(`language: ""
env:
  - FOO=foo BAR=bar
  - FOO=bar BAR=foo
`))

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
	travis, err := NewTravis([]byte(`language: "go"
env:
  global:
    - DB=postgres
    - SH=bash
    - PACKAGE_VERSION="1.0.*"
`))
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
	travis, err := NewTravis([]byte(`language: go
env:
  matrix:
    - FOO=foo BAR=bar
    - FOO=bar BAR=foo
`))
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
	travis, err := NewTravis([]byte(`language: go
env:
  global:
    - DB=postgres
    - SH=bash
    - PACKAGE_VERSION="1.0.*"
  matrix:
    - FOO=foo BAR=bar
    - FOO=bar BAR=foo
`))
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

func TestParseSecretEnv(t *testing.T) {

	globals := []string{
		"DB=postgres",
		"SH=bash",
		"PACKAGE_VERSION=\"1.0.*\"",
	}
	travis, err := NewTravis([]byte(`language: go
env:
  global:
    - DB=postgres
    - SH=bash
    - PACKAGE_VERSION="1.0.*"
    - secret: xxxxxxxxxxxxxxx
`))
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(travis.Env.Global) != 3 {
		t.Fatal("The number of global variables is wrong:", travis.Env.Global)
	}
	for i, v := range globals {
		if travis.Env.Global[i] != v {
			t.Error("A global variable is not match:", travis.Env.Global)
		}
	}

}

func TestParseEnvStrings(t *testing.T) {

	var res TestCase
	res = parseEnv("FOO=bar")
	if len(res) != 1 || res["FOO"] != "bar" {
		t.Error("parseEnv returns wrong envs:", res)
	}

	res = parseEnv("FOO=bar BAR=fuga")
	if len(res) != 2 || res["FOO"] != "bar" || res["BAR"] != "fuga" {
		t.Error("parseEnv returns wrong envs:", res)
	}

	res = parseEnv(`FOO="bar fuga"`)
	if len(res) != 1 || res["FOO"] != "bar fuga" {
		t.Error("parseEnv returns wrong envs:", res)
	}

	res = parseEnv(`FOO="bar fuga" BAR="foo fuga"`)
	if len(res) != 2 || res["FOO"] != "bar fuga" || res["BAR"] != "foo fuga" {
		t.Error("parseEnv returns wrong envs:", res)
	}

}
