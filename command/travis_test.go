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
	travis, err := storeAndLoadTravis(&Travis{
		RawBeforeInstall: "install 1",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeInstall) != 1 || travis.BeforeInstall[0] != "install 1" {
		t.Error("BeforeInstall is wrong:", travis.BeforeInstall)
	}
}

func TestParseBeforeInstallWithList(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{
		RawBeforeInstall: []string{"install 1"},
	})
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
	travis, err := storeAndLoadTravis(&Travis{
		RawInstall: "install 1",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Install) != 1 || travis.Install[0] != "install 1" {
		t.Error("Install is wrong:", travis.Install)
	}
}

func TestParseInstallWithList(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{
		RawInstall: []string{"install 1"},
	})
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
	travis, err := storeAndLoadTravis(&Travis{
		RawBeforeScript: "python setup.py test",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.BeforeScript) != 1 || travis.BeforeScript[0] != "python setup.py test" {
		t.Error("BeforeScript is wrong:", travis.BeforeScript)
	}
}

func TestParseBeforeScriptWithList(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{
		RawBeforeScript: []string{"python setup.py test"},
	})
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
	travis, err := storeAndLoadTravis(&Travis{
		RawScript: "python setup.py test",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(travis.Script) != 1 || travis.Script[0] != "python setup.py test" {
		t.Error("Script is wrong:", travis.Script)
	}
}

func TestParseScriptWithList(t *testing.T) {
	travis, err := storeAndLoadTravis(&Travis{
		RawScript: []string{"python setup.py test"},
	})
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
