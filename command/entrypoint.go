//
// command/entrypoint.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"bytes"
	"text/template"
)

// EntrypointAsset defines a asset name for entrypoint.sh.
const EntrypointAsset = "asset/entrypoint.sh"

// Entrypoint creates an entrypoint.sh from an instance of Travis.
func Entrypoint(travis *Travis) (res []byte, err error) {

	data, err := Asset(EntrypointAsset)
	if err != nil {
		return
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	if err = temp.Execute(&buf, travis); err != nil {
		return
	}

	res = buf.Bytes()
	return

}
