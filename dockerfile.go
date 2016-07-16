package main

import (
	"bytes"
	"os"
	"text/template"
)

// DockerfileAsset defines a asset name for Dockerfile.
const DockerfileAsset = "asset/Dockerfile"

type travisExt struct {
	*Travis
	PWD string
}

// Dockerfile creates a Dockerfile from an instance of Travis.
func Dockerfile(travis *Travis) (res []byte, err error) {

	data, err := Asset(DockerfileAsset)
	if err != nil {
		return
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	param := travisExt{
		Travis: travis,
		PWD:    pwd,
	}

	buf := bytes.Buffer{}
	if err = temp.Execute(&buf, &param); err != nil {
		return
	}

	res = buf.Bytes()
	return

}
