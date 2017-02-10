//
// command/travis_util_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import "testing"

func TestParseEnvStrings(t *testing.T) {

	var res []string
	res = parseEnv("FOO=bar")
	if len(res) != 1 || res[0] != "FOO=bar" {
		t.Error("parseEnv returns wrong envs:", res)
	}

	res = parseEnv("FOO=bar BAR=fuga")
	if len(res) != 2 || res[0] != "FOO=bar" || res[1] != "BAR=fuga" {
		t.Error("parseEnv returns wrong envs:", res)
	}

	res = parseEnv("FOO=\"bar fuga\"")
	if len(res) != 1 || res[0] != "FOO=bar fuga" {
		t.Error("parseEnv returns wrong envs:", res)
	}

	res = parseEnv("FOO=\"bar fuga\" BAR=\"foo fuga\"")
	if len(res) != 2 || res[0] != "FOO=bar fuga" || res[1] != "BAR=foo fuga" {
		t.Error("parseEnv returns wrong envs:", res)
	}

}
