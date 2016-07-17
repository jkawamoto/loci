package command

import (
	"bytes"
	"text/template"
)

// DockerfileAsset defines a asset name for Dockerfile.
const DockerfileAsset = "asset/Dockerfile"

type travisExt struct {
	*Travis
	Archive string
}

// NewDockerfile creates a Dockerfile from an instance of Travis.
func NewDockerfile(travis *Travis, archive string) (res []byte, err error) {

	data, err := Asset(DockerfileAsset)
	if err != nil {
		return
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	param := travisExt{
		Travis:  travis,
		Archive: archive,
	}

	buf := bytes.Buffer{}
	if err = temp.Execute(&buf, &param); err != nil {
		return
	}

	res = buf.Bytes()
	return

}
