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
	"fmt"
	"text/template"
)

// EntrypointAsset defines a asset name for entrypoint.sh.
const EntrypointAsset = "assets/entrypoint"

// Entrypoint creates an entrypoint.sh from an instance of Travis.
func Entrypoint(travis *Travis) (res []byte, err error) {

	var data []byte
	name := fmt.Sprintf("%s-%s", EntrypointAsset, travis.Language)
	data, err = Asset(name)
	if err != nil {
		data, err = Asset(fmt.Sprintf("%s", EntrypointAsset))
		if err != nil {
			return
		}
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
