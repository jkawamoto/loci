//
// command/entrypoint.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"bytes"
	"fmt"
	"text/template"
)

// EntrypointAsset defines a asset name for entrypoint.sh.
const EntrypointAsset = "assets/entrypoint"

// Entrypoint creates an entrypoint.sh from an instance of Travis.
func Entrypoint(travis *Travis) (res []byte, err error) {

	var (
		data []byte
		temp *template.Template
	)

	// Loading the base template.
	data, err = Asset(EntrypointAsset)
	if err != nil {
		return
	}
	base, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	// Loading a child template.
	name := fmt.Sprintf("%s-%s", EntrypointAsset, travis.Language)
	data, err = Asset(name)
	if err == nil {
		temp, err = base.Parse(string(data))
		if err != nil {
			return
		}
	} else {
		temp = base
	}

	// Creating an entrypont.
	buf := bytes.Buffer{}
	if err = temp.ExecuteTemplate(&buf, "base", travis); err != nil {
		return
	}
	res = buf.Bytes()

	return

}
