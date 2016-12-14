//
// command/travis_test.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//
package command

import "testing"

// TestArgumentSet tests ArgumentSet method returns correct argument sets.
func TestArgumentSet(t *testing.T) {

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
