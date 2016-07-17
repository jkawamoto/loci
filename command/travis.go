package command

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Travis defines the structure of .travis.yml.
type Travis struct {
	Language string
	Addons   struct {
		Apt struct {
			Packages []string
		} `yaml:"apt,omitempty"`
	}
	Install      []string `yaml:"install,omitempty"`
	BeforeScript []string `yaml:"before_script,omitempty"`
	Script       []string `yaml:"script,omitempty"`
}

// NewTravis loads a .travis.yaml file and creates a structure instance.
func NewTravis(filename string) (res *Travis, err error) {

	fp, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	res = &Travis{}
	if err = yaml.Unmarshal(buf, res); err != nil {
		return
	}
	return

}
