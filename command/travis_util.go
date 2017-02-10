//
// command/travis_util.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import "strings"

func parseEnv(env string) (envs []string) {

	b := 0
	quoted := false
	for i, v := range env {
		if v == '"' {
			quoted = !quoted
		}
		if !quoted && v == ' ' {
			envs = append(envs, strings.Replace(env[b:i], "\"", "", 2))
			b = i + 1
		}
	}
	return append(envs, strings.Replace(env[b:], "\"", "", 2))

}
