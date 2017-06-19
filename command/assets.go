// Code generated by go-bindata.
// sources:
// assets/Dockerfile
// assets/Dockerfile-go
// assets/Dockerfile-python
// assets/entrypoint
// assets/entrypoint-go
// assets/entrypoint-python
// DO NOT EDIT!

package command

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsDockerfile = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x94\xe1\x8f\x9b\x36\x18\xc6\xbf\xf3\x57\xbc\x22\xb7\x6a\xab\x86\xc9\xf5\xc3\x26\x5d\x97\xd3\x68\x93\x6e\xac\x0b\x20\x8e\x75\x3b\x35\x55\xe5\x98\x37\xe0\x05\x6c\xcf\x36\x49\x33\xc4\xfe\xf6\x09\xc2\x8e\xb4\xd7\xf5\x43\x94\xbc\xf6\xef\x7d\xfc\xbc\xd6\x13\xcf\x9c\x19\x2c\x25\xdb\xa3\xde\xf1\x0a\x9d\xbe\x7c\x29\xd5\x49\xf3\xa2\xb4\xf0\x35\xfb\x06\x9e\xcd\xaf\xbf\xf3\x9e\xcd\xaf\xbf\x87\x5f\x1a\xa1\x90\xc3\x6b\x7a\xa4\xb5\xb4\x72\x60\xb3\x92\x1b\x30\x72\x67\x8f\x54\x23\x70\x03\x1a\x2b\xa4\x06\x73\x68\x44\x8e\x1a\x6c\x89\xb0\x0e\x33\xf8\x95\x33\x14\x06\xc9\xd0\x54\x5a\xab\x6e\x7c\x5f\x2a\x14\x46\x36\x9a\x21\x91\xba\xf0\xab\x33\x62\xfc\x9a\x5b\x6f\x2c\x88\x2a\x95\x33\x73\xda\x36\xc7\x1d\x17\x08\xee\x96\x1a\x74\xbb\xce\x79\x95\xc6\x6b\x68\x5b\xf2\x82\x1a\x0c\x6b\x5a\x60\xd7\x39\xeb\x20\x8c\xb2\x20\x8c\x56\xe9\xa7\x56\xe1\x87\xfd\xf8\x8b\xfc\x39\xec\xfc\x58\xd4\x94\x57\x84\xc9\xfa\xd6\x71\x56\xd1\x1b\xc8\x56\xe9\x1a\x0e\xf6\x7a\x3e\x1f\xca\xe5\xea\x45\x18\x44\xef\x5f\xa5\x71\x94\xad\xa2\x25\x08\x29\xb8\xb0\xa8\x29\xb3\xfc\x80\x4e\xdb\x1e\xb9\x2d\x81\xfc\x9c\x65\x49\xa2\xe5\x87\x53\xd7\x0d\x6d\x7d\xfd\x3e\x49\xe3\x3f\xee\x7b\x6f\x5d\xe7\xb4\x2d\x8a\x7c\xf8\x9e\x1a\xee\x3e\xed\xb8\xfb\x62\x4b\x24\x2f\xf9\x28\xfe\x2c\xfc\x40\x07\xca\xfe\x87\xa7\xbf\x45\x80\xac\x94\xe0\x06\xec\xaf\x86\x6b\xbc\xb9\xe9\x2f\x1e\x5a\x18\x08\xd8\xb8\x83\xc6\xc6\x7d\x0e\xdd\x73\x17\x6e\x6f\xc1\x47\xcb\x7c\xaa\x6c\xff\x21\x4c\x8a\x1d\xc9\xfd\xf9\xb5\xea\xe9\xc7\x27\x25\x27\xc5\x2f\x8f\x4a\xee\x93\xf0\xec\x6d\x31\xe8\xc2\x93\x27\xb0\x71\x00\x00\xc2\x24\x89\xd3\x6c\x71\xd5\x4e\xc8\xec\xa9\xef\x5f\x10\xf5\x3e\xe7\x1a\x3c\x05\xff\xf8\x44\x71\xe5\x4f\x3b\xe7\x01\xde\x16\x95\xdc\xd2\xea\xdd\x46\x70\x91\xe3\x07\xaf\xd1\xd5\xe2\x6a\x52\xf3\xb5\x94\xd6\x57\x27\xc5\x37\xc2\xea\xc6\x58\xcc\xbd\x52\x1a\xbb\xb8\x6a\xcf\x67\x7f\x75\xf3\xb4\x1b\x26\x1c\xf5\x15\x57\xc3\x7c\xd3\x54\x33\x08\x85\xb1\xb4\xaa\xe0\xa5\xac\x6b\x29\x20\x48\x32\x50\x94\xed\x69\x81\x86\x0c\x03\x52\x65\xbd\x02\x2d\x34\x2a\xa7\x16\x27\x8b\xd3\x7a\xa1\x69\x8e\xe0\x9d\x1e\xef\xf1\x51\xdc\x3b\x0d\x4b\x8d\xe5\x95\x81\x82\x5b\x60\x8d\xae\xe0\xd8\x23\xa6\xc9\x25\x34\xe2\x6f\xae\x2e\xcc\x54\x54\x14\x0d\x2d\x10\x8c\x42\xc6\x77\x9c\x5d\x78\x6a\xdb\x6d\x25\xd9\x1e\xdc\x71\xc9\x85\xcb\x44\x4c\x1a\xba\x11\x96\xd7\x68\xc0\x4a\xd8\x22\x34\xfd\x9f\x93\x0b\xb0\x68\xac\x21\x4e\x90\xfe\x04\x6f\x56\xe9\x5d\x18\x47\x93\xe2\x01\xb5\xe1\x52\xfc\xaf\x22\x0e\x79\xca\x3f\x32\xa3\xa9\x28\x10\x48\xa6\xe9\x81\x1b\x12\xe4\xb9\x14\xa6\x8f\x23\x49\x46\x68\xcc\xc9\x67\xae\xe4\x51\x96\x47\x1b\xe7\xc7\xe1\xec\x22\x58\x2e\x7b\x2e\xd0\xac\xe4\x07\xec\x3a\xf0\x73\x6a\xa9\xf3\x7b\x9c\xbe\x5e\x86\xe9\x58\x3d\x48\xf4\x34\x0a\xab\x4f\x4a\x72\x61\x89\x29\x61\xc8\x88\xb3\x8a\xb2\xf4\x3e\x89\xc3\x28\x83\xb7\xfd\x5b\x52\xba\xdf\x82\x7b\x8e\xcf\x47\xb8\xfb\xee\x41\xeb\xdf\x00\x00\x00\xff\xff\x5a\x6a\xb6\x11\x21\x05\x00\x00")

func assetsDockerfileBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerfile,
		"assets/Dockerfile",
	)
}

func assetsDockerfile() (*asset, error) {
	bytes, err := assetsDockerfileBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/Dockerfile", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockerfileGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x91\x51\x6b\xdb\x3e\x14\xc5\xdf\xf5\x29\x2e\x4e\x28\x2d\x7f\x6c\xb5\x7d\xf8\x0f\x0a\x79\x08\x4b\xc8\xbc\x36\x49\xc9\xb2\xee\x61\x1b\x45\x95\xaf\xed\x4b\x6c\x49\x48\x72\xb2\xa0\xfa\xbb\x8f\xd8\xa1\x09\x6c\x83\xbe\x49\x57\xe7\x77\xef\xb9\x47\x03\x36\x80\x89\x96\x1b\xb4\x39\x55\x18\x17\x9a\x1d\x2a\x1f\xb5\xd9\x5b\x2a\x4a\x0f\x97\xf2\x0a\x6e\xaf\x6f\xfe\x8f\x6f\xaf\x6f\x3e\xc0\xe7\x46\x19\x24\xb8\x17\x3b\x51\x6b\xdf\x6b\xd7\x25\x39\x70\x3a\xf7\x3b\x61\x11\xc8\x81\xc5\x0a\x85\xc3\x0c\x1a\x95\xa1\x05\x5f\x22\xcc\xd3\x35\x3c\x90\x44\xe5\x30\xe9\xa0\xd2\x7b\x73\xc7\xb9\x36\xa8\x9c\x6e\xac\xc4\x44\xdb\x82\x57\xbd\xc4\xf1\x9a\x7c\x7c\xbc\x24\xa6\x34\x6c\xc0\x42\xc8\x30\x27\x85\x10\x19\x21\x37\xa2\xc0\xa8\x6d\xd9\x74\xf1\x04\xb3\xe5\xe3\x78\xfd\x09\x78\x26\xbc\x60\xab\xaf\x0b\x10\xc6\xc7\x05\x7a\x20\xe5\xbc\xa8\x2a\x88\xf7\x50\xe8\x4a\xa8\x82\xb1\x01\xa4\xc7\x62\x41\x75\x8d\x9d\xbc\xde\x64\x64\x21\x36\x30\xec\x3b\xf1\x17\x52\x70\x71\x01\x3f\x18\x00\x80\x6c\x6c\x05\xb1\x7b\x80\x58\x9f\x0b\x78\xc7\x77\x5b\xb8\x3b\xce\xad\xd8\x25\x05\xf9\xb2\x79\x69\x1c\x5a\xa9\x95\x47\xe5\x13\xa9\x6b\xee\xad\xd8\x92\x8b\x25\xf5\x04\xaf\x85\xf3\x68\x8f\xf8\x69\x4a\x59\xeb\x0c\xfe\xfb\xf5\xe7\x08\x16\x02\xaa\xac\x6d\xd9\x59\x00\x5b\xb4\x8e\xb4\x3a\x04\x70\x5a\x48\x80\x33\x28\x29\x27\x09\xc7\x77\xd0\x39\xcc\x74\xb7\xe3\xa1\xe7\x28\x3a\x6f\x7e\x37\x3c\x1c\xa3\x93\x85\xef\x30\x7c\x9a\xae\xbe\xa4\xcb\x05\x8c\x20\x12\x6a\x1f\xc1\x4f\x78\x7d\x05\xdc\x8a\x0a\xa2\xe1\xe5\x2c\x9d\xcf\xa7\xcf\xb3\xe5\xf3\x51\x35\x7a\x93\x77\x3e\xaf\xa2\xbf\x39\xed\xbf\xf6\x60\x34\x04\xca\x21\x59\x77\x69\x24\x33\x9d\xd6\x46\x5b\xff\x28\x7c\xd9\xb6\x6c\x3c\x99\x40\x08\xc9\xd8\xca\x92\xb6\xd8\xb6\x6f\x29\x38\x2b\x79\x08\xff\xa0\x42\xc0\xca\xe1\x3b\xf0\x15\x1a\xed\xc8\x6b\xbb\xef\xa9\xce\xe4\xb7\xe5\xea\x7e\x92\xae\xde\xa5\xfd\x1d\x00\x00\xff\xff\xfa\x4c\x80\x99\x23\x03\x00\x00")

func assetsDockerfileGoBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerfileGo,
		"assets/Dockerfile-go",
	)
}

func assetsDockerfileGo() (*asset, error) {
	bytes, err := assetsDockerfileGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/Dockerfile-go", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockerfilePython = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\x51\x4f\xdb\x3c\x14\x7d\xcf\xaf\xb8\x4a\x2b\x44\x1f\x1c\x93\x7e\x12\x9f\x84\xc4\x03\x1b\x48\x74\x83\x52\x41\x41\x42\xdb\x40\x4e\x72\x9b\x5c\xd5\xb5\x3d\xdb\x49\x17\xda\xfe\xf7\x29\x69\x68\x51\x35\xf2\x92\xf8\x9e\x73\x6e\xce\x39\x72\x2f\xe8\xc1\xa5\x4e\xe7\x68\x67\x24\x91\x99\xda\x17\x5a\x05\xcd\xf4\xab\x36\xb5\xa5\xbc\xf0\x70\x9c\x0e\x60\x78\x12\x9f\xb2\xe1\x49\xfc\x3f\x7c\x2b\x95\x41\x82\xef\x62\x29\x16\xda\xeb\x96\x3b\x2d\xc8\x81\xd3\x33\xbf\x14\x16\x81\x1c\x58\x94\x28\x1c\x66\x50\xaa\x0c\x2d\xf8\x02\xe1\x76\x34\x85\x1b\x4a\x51\x39\x8c\x5a\x51\xe1\xbd\x39\xe3\x5c\x1b\x54\x4e\x97\x36\xc5\x48\xdb\x9c\xcb\x2d\xc5\xf1\x05\x79\xd6\x1d\x22\x53\x98\xa0\x17\xac\x56\x19\xce\x48\x21\x84\x46\xa4\x73\x91\x63\xb8\xd9\x04\x3d\x18\x29\xe7\x85\x94\x20\x8c\x87\x0e\x70\x51\x70\xff\x38\x6e\x26\x2c\x47\x0f\xd4\x31\x58\x0d\xdb\x80\xdd\x8b\x19\x32\x20\x2c\x89\x21\xfc\x0c\xe0\xdf\xcf\x9b\xa4\x24\xce\x59\x86\x15\x48\x4a\x2c\x8a\x4c\x92\xc2\xd3\x76\xf0\xa9\x48\x52\x92\xbc\x0d\xdf\x45\xee\xb7\x24\x8f\xff\xed\x8e\x4e\x36\x9f\x41\x0f\x1e\x4d\x26\x3c\x82\x21\xd3\xda\x6d\xdc\xec\xac\xb2\xd2\xe4\x56\x64\x5b\x74\x1f\xd2\xd4\xa8\xaa\x96\x9d\x96\x56\x02\xbb\x69\x6b\x74\x67\x9c\x5b\xb1\x8c\x72\xf2\x45\x99\x94\x0e\x6d\xaa\x95\x47\xe5\xa3\x54\x2f\x78\x5d\x97\x25\x6f\x85\xac\x5b\x8f\x96\x2f\x84\xf3\x68\x79\x42\xea\x10\x82\x35\x24\xc2\x15\xc1\x6a\x85\x2a\xdb\x6c\x82\x0f\xbd\x57\x68\x1d\x69\x75\xd0\x3b\x38\x83\x29\xcd\x28\x7d\x6f\xb7\xa3\x05\x57\xe3\x27\x98\x3c\x4f\xaf\xef\xc6\xaf\x5f\x1e\x47\x37\x97\xaf\x17\xf7\xa3\x8b\xe1\xeb\xdd\x64\xfa\x00\x21\xfb\x03\xf1\x09\xb0\x39\xc4\xb7\x61\x1b\x68\x72\x31\xbd\x3e\x0f\xb9\xd5\xda\xf3\xa8\xf5\xd4\xb8\x3b\xeb\x37\xf3\x10\x8e\x8e\xba\xba\xb1\x12\x12\xc2\xfe\x71\xcb\x00\x52\xe4\x81\x0d\x3e\xc5\x2b\xb2\xbe\x14\x72\x1b\xf0\x90\x3a\x79\x7e\xba\xba\x7f\x18\xdd\x8d\xcf\xf7\xeb\xba\xfa\x25\xac\xa1\xb9\xbf\xa1\xe3\xc0\x79\x1e\xc2\x1a\x72\x8b\x06\xd8\x15\x84\x2f\xfd\x4e\x76\xdc\x5f\xff\x78\x61\xbf\x06\x0d\xea\x05\x49\x60\xf1\x60\xbf\xfd\x60\xa3\x03\x56\x41\x7f\xf7\xcb\x43\x9e\xd4\xa9\x90\x1f\xf0\x5d\xfd\x7f\x03\x00\x00\xff\xff\x87\xdc\x50\xff\xa2\x03\x00\x00")

func assetsDockerfilePythonBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerfilePython,
		"assets/Dockerfile-python",
	)
}

func assetsDockerfilePython() (*asset, error) {
	bytes, err := assetsDockerfilePythonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/Dockerfile-python", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsEntrypoint = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x92\x4d\x6f\x9c\x30\x18\x84\xef\xfe\x15\x13\xc8\x21\x91\x36\x90\xe4\xd0\x4a\xa9\x2a\x55\xad\x54\x69\xfb\x71\x69\x72\x8b\xa2\xc6\xc0\x00\x56\x59\xdb\xb2\x4d\x36\x2b\xc4\x7f\xaf\xc0\x34\xab\x36\xd9\x43\x7a\x82\xf7\xf3\x79\x67\x60\x18\x2a\xd6\x4a\x13\x49\x21\x3d\x93\x71\x14\xe9\x51\x5e\x28\x9d\x17\xd2\xb7\x22\x15\x29\xa8\x83\xdb\x59\xa3\x74\x98\xc3\x4f\xc6\xee\x9c\x6a\xda\x80\x93\xf2\x14\x97\xe7\x17\x6f\xce\x2e\xcf\x2f\xde\xe2\x4b\xaf\x2d\x15\xbe\xca\xad\xdc\x98\x60\xe6\xde\x9b\x56\x79\x78\x53\x87\xad\x74\x84\xf2\x70\xec\x28\x3d\x2b\xf4\xba\xa2\x43\x68\x89\xef\xeb\x1b\x7c\x53\x25\xb5\x67\x36\x0f\xb5\x21\xd8\xab\x3c\x37\x96\xda\x9b\xde\x95\xcc\x8c\x6b\xf2\x2e\xb6\xf8\x7c\xa3\xc2\xd9\x12\x64\xb6\xb5\x22\x15\x22\xc5\x8f\x5e\xa3\x60\x6d\x26\x8a\xf6\x41\x76\x1d\x7c\xa0\xf5\xab\x7f\xc3\xa5\xc9\x97\x4e\xd9\xf0\x27\x29\x75\xb5\x64\x44\x1a\x73\x50\x1a\x61\x3a\xde\xb8\x8a\x2e\xc3\xba\x86\xd4\x3b\x58\xe9\xe4\x86\x81\xce\x63\x12\xd4\xa8\x07\xea\x15\xf8\xc8\xb2\x0f\x9c\xd5\xec\x3b\x44\x3a\xb3\x29\xab\x0c\x9f\x8d\x03\x1f\xe5\xc6\x76\x5c\x21\x18\x54\x2c\xfa\x26\x02\x22\x77\x05\xd7\x2f\xc4\xad\x0a\x2d\xee\x27\xf7\xef\x21\xfd\xb4\x54\xa4\xfb\xb5\xf3\xad\x0d\x03\x24\x7c\xcb\xae\x83\x2c\x4b\x7a\x9f\x09\xcf\x80\x33\x8a\x61\x28\x3a\x53\xfe\x42\x42\xfd\x90\x20\x1b\x47\x31\x0c\xd4\xd5\x38\x0a\x55\xe3\xf6\x16\xc7\x29\x8e\xde\xe3\x1c\x77\x77\xef\xa6\xd5\x5a\x60\x3e\x1f\xc7\x1f\x44\xad\xf6\xd3\xd6\xd1\x4a\xc7\xbf\x37\xec\xcb\xd1\xc5\x9f\x8b\xb7\xb1\x8b\x65\x6b\x90\x7c\x8c\xfe\xae\x17\xd7\xaf\x27\x33\xaf\x12\x31\x0c\x4e\xea\x86\xc8\x62\x7d\x29\x3f\x4d\x0d\x43\x36\x8e\x89\x38\x99\x9f\xa7\x4f\xc0\xe7\xe0\x17\x88\x07\x51\xff\x0d\x59\xd4\xc5\x2f\xf3\x82\xb8\xeb\xf8\xf3\x1c\xd0\x16\xab\xaf\xa7\x3e\xc7\x1d\xe2\xbc\x9a\x10\x5f\x7e\x07\x00\x00\xff\xff\x25\xf1\x25\xb8\xea\x03\x00\x00")

func assetsEntrypointBytes() ([]byte, error) {
	return bindataRead(
		_assetsEntrypoint,
		"assets/entrypoint",
	)
}

func assetsEntrypoint() (*asset, error) {
	bytes, err := assetsEntrypointBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/entrypoint", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsEntrypointGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x53\x4d\x6f\xdc\x36\x14\xbc\xf3\x57\x4c\x25\x03\xb5\x03\x8b\x4a\x52\xa0\x05\x5c\xe4\xe0\xb4\xe9\x76\xdb\x3a\x2e\xbc\x6e\x2e\x81\x91\x70\xa5\x27\x8a\x30\x45\x0a\x24\xb5\xb6\x21\xeb\xbf\x17\x94\xb4\x1f\xd8\xf8\x90\xcb\xae\x34\x1c\xbe\x37\xf3\xde\x28\xfd\x21\x5f\x2b\x93\xaf\x85\xaf\x59\xca\x52\x90\x09\xee\xa9\xb5\xca\x84\x4c\xda\x11\xf9\xcd\xb6\x4f\x4e\xc9\x3a\xe0\xb4\x38\xc3\xdb\xd7\x6f\x7e\xce\xde\xbe\x7e\xf3\x0b\xfe\xea\x4c\x4b\x0a\x7f\x8b\x07\xd1\xd8\x30\x71\x6f\x6b\xe5\xe1\x6d\x15\x1e\x84\x23\x28\x0f\x47\x9a\x84\xa7\x12\x9d\x29\xc9\x21\xd4\x84\xab\xe5\x2d\xfe\x51\x05\x19\x4f\x7c\xbc\x54\x87\xd0\x5e\xe4\xb9\x6d\xc9\x78\xdb\xb9\x82\xb8\x75\x32\xd7\x13\xc5\xe7\x8d\x0a\xd9\xfc\xc2\xdb\xba\x65\x29\x63\x29\x6e\x3a\x83\x35\x55\x36\x76\x31\x3e\x08\xad\xe1\x03\xb5\xfe\xfc\xf8\x75\x26\xf9\xc2\xa9\x36\x6c\x41\x61\xca\x19\x61\xe9\x84\x41\x19\x84\x28\xde\xba\x92\x1c\x9f\x8c\x8c\xb3\xc0\x38\x8c\x6d\x81\xca\x16\x9d\x27\x0f\x6b\xe0\x3a\x63\x94\x91\x08\xe4\x83\x47\x65\x1d\x4b\xb1\xb0\x10\x6d\xab\x55\x21\x82\xb2\xc6\x73\x2c\x2b\x08\xf3\x84\x56\x38\xd1\x50\x20\xe7\x11\x07\x23\xd5\x86\xcc\x39\xe8\x91\x8a\x2e\xd0\x38\x95\x3d\x83\xa5\xa3\x07\x12\x25\xc7\x1f\xd6\x81\x1e\x45\xd3\x6a\x3a\x47\xb0\x28\x69\xdd\xc9\x49\xe8\x24\xe8\x3c\xca\x98\x80\x07\x15\x6a\x7c\x8d\x8b\xfc\x0a\xe1\x63\x51\x96\xee\xcb\x8e\x9e\x25\x05\x08\xf8\x9a\xb4\x86\x28\x0a\xf2\x9e\xb3\xbe\x2f\xa9\x52\x86\x90\x90\xd9\x24\xc3\xc0\xfe\xbd\xbc\xfd\xf3\x5d\x72\xb2\xb8\x8e\x0f\x31\x1c\x17\x27\xf1\x29\x61\x25\x15\x3a\xaa\xcf\x1e\xb1\x58\x5e\x5d\x7d\xf8\xb2\xb8\xfe\xf2\xe9\xc3\xcd\x6a\x79\xfd\xf1\xdd\xc9\xa9\x54\x4d\x43\xc8\xf4\x19\x53\x15\x3e\x7f\x46\x66\x90\x9c\x1c\xf3\x12\xdc\xdd\xfd\x1a\xa5\x19\x06\xd0\x46\x68\x24\xf3\xcd\xb3\x84\x55\x8a\x39\x12\xa5\x35\xfa\x09\xb7\x37\x97\x9f\x96\xab\xa3\x0e\x16\x1b\x72\x5e\x59\x83\x67\x14\x5d\x40\x56\xe1\x27\x64\x25\x12\x24\x78\x46\x4c\x99\xcf\xa5\xcd\x73\x79\xc6\xfa\x9e\x4c\x39\x0c\xec\xc0\xde\x1c\x8c\x68\x91\x8a\xda\x22\x59\xce\x49\x59\xc5\x00\x5c\x24\x4c\xaa\x80\xc2\x9a\x4a\x49\x64\x99\xd4\x76\x2d\xf4\x98\x4d\x1e\x7f\xfc\x45\x9e\x4b\xdb\xde\x4b\xae\x0c\xaf\xac\xd6\xf6\xe1\x86\x4a\xe5\xa8\x08\x1e\xc1\x75\xc4\xfa\x5e\x55\xe0\x73\xd1\x61\x60\x7d\xef\x84\x91\x74\x08\x4d\x7d\xfb\x9e\x0f\x43\xc2\x4e\xc7\xff\xbd\xd4\xbe\x27\xed\x69\xc7\x92\x76\x5c\x57\x16\xc0\x73\xce\x79\xc2\x8e\x80\xc3\x7b\xc7\x56\xa7\x6c\xec\x9d\xae\xa6\xf0\x6e\x8d\xf6\x7d\xfe\x0a\x89\xaa\xb6\xa1\x56\xde\xfc\x18\xb6\xa1\xec\x3c\xa1\x11\xf7\x04\x3b\x7d\xae\x25\x55\xa2\xd3\x01\xd2\x6a\x31\xa7\x1d\x85\x6d\x1a\x61\x4a\x9e\xe0\x55\x3e\x2a\x88\xce\xa7\x26\x87\xc6\x77\xc8\xf7\xf9\x9e\x73\x43\x58\x7c\xfc\xef\x4a\xdc\x53\xa5\x34\xe1\xee\x0e\xcf\xcf\x33\xfe\x22\xf8\x7e\xf5\x7b\xf3\x12\x7e\x00\xee\x12\x17\x31\x16\xbb\x31\x60\x12\xbd\xb0\xef\x3b\xa5\xcb\x4b\x27\xfd\x30\x30\x40\xda\xc9\x61\xdf\x7f\x73\xb4\xd3\xb9\x67\x65\x9b\x79\x17\xe3\xe9\xe8\xa6\x52\xdf\x2c\xe6\xff\x00\x00\x00\xff\xff\xed\x66\x03\x53\x63\x05\x00\x00")

func assetsEntrypointGoBytes() ([]byte, error) {
	return bindataRead(
		_assetsEntrypointGo,
		"assets/entrypoint-go",
	)
}

func assetsEntrypointGo() (*asset, error) {
	bytes, err := assetsEntrypointGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/entrypoint-go", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsEntrypointPython = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x52\xd1\x6a\xdb\x4a\x14\x7c\xdf\xaf\x98\x2b\x85\x90\x80\x25\x25\x79\xb8\x17\x7c\xc9\x43\x69\x53\xea\xb6\x69\x42\xe2\x16\x4a\x1a\xc8\x5a\x3a\x92\x16\xe4\xb3\xdb\xdd\x23\xc7\x46\xf8\xdf\x8b\x64\xd9\x29\x69\x9f\xc4\xce\x8e\xe6\xcc\xce\x9c\xf8\x9f\x6c\x61\x38\x5b\xe8\x50\xab\x58\xc5\x20\x16\xbf\x71\xd6\xb0\x24\x6e\x23\xb5\xe5\x01\x7d\x6b\xdd\xc6\x9b\xaa\x16\x9c\xe4\xa7\xb8\x38\x3b\xff\x37\xb9\x38\x3b\xff\x0f\x1f\x5b\x76\x64\xf0\x49\x3f\xeb\xa5\x15\x3b\x70\xe7\xb5\x09\x08\xb6\x94\x67\xed\x09\x26\xc0\x53\x43\x3a\x50\x81\x96\x0b\xf2\x90\x9a\x70\x3d\x9b\xe3\xb3\xc9\x89\x03\xa5\xc3\x4f\xb5\x88\x9b\x66\x99\x75\xc4\xc1\xb6\x3e\xa7\xd4\xfa\x2a\x6b\x76\x94\x90\x2d\x8d\x24\xe3\x21\x75\xb5\x53\xb1\x52\x31\xee\x5a\xc6\x82\x4a\xdb\x4f\xe1\x20\xba\x69\x10\x84\x5c\x98\xbc\x3e\x8e\xa4\x90\x7b\xe3\x64\x0f\x6a\x2e\x46\x44\xc5\x3b\x0c\x86\x21\xbd\x79\xeb\x0b\xf2\xe9\xee\x21\x43\x1e\x18\x02\xd9\x0b\x94\x36\x6f\x03\x05\x58\x86\x6f\x99\x0d\x57\x10\x0a\x12\x50\x5a\xaf\x62\xdc\x0e\xb1\x41\x3b\xd7\x98\x5c\x8b\xb1\x1c\x52\xcc\x4a\x68\xde\xc0\x69\xaf\x97\x24\xe4\x03\xfa\x70\x2a\xb3\x22\x9e\x80\xd6\x94\xb7\x42\x43\x32\x2f\x0c\x15\x0f\xef\x20\x5d\xa4\x78\x6f\x3d\x68\xad\x97\xae\xa1\x09\xc4\xa2\xa0\x45\x5b\xed\xcc\xee\x4c\x4d\x7a\x2b\x3b\xe0\xd9\x48\x8d\xa7\xbe\xd0\x27\xe8\xd0\x8b\xaa\xf8\x45\x76\x78\x77\x45\x02\x8d\x50\x53\xd3\x40\xe7\x39\x85\x90\xaa\xae\x2b\xa8\x34\x4c\x88\x88\x57\xd1\x76\xab\x6e\xdf\xcc\x3f\x5c\x46\x99\xb7\x56\xb2\xd4\x6d\x88\x57\xfd\xa6\x4c\x8f\x7a\x3c\x52\xb4\xd2\x0d\xa2\xa3\x93\xe1\x02\x86\x8d\x20\x39\x8d\x70\x7c\x8c\x1f\xaf\xee\x56\xc6\x4b\xab\x1b\xe2\x55\xb2\xa7\xa9\xae\x23\x2e\xb6\x5b\xf5\xdb\x54\xe7\xc9\x69\x4f\xfd\xe4\xf9\xd5\xdd\xf5\xe5\x5a\xc8\x2f\x95\x27\x5d\x58\x6e\x36\xb8\xfd\xfe\xed\xea\xee\x7e\x76\xf3\xe5\x72\xaf\xdb\xd8\x5c\x37\xa7\x7f\x93\x1a\xeb\xef\xa5\x28\xaf\x2d\xa2\xd9\xb8\x0f\xf7\x7d\xcd\xd3\x48\x99\x12\x0f\x0f\x48\x08\x9e\x7e\xb6\xc6\xd3\x92\x58\x42\x2a\x6b\xc1\xe3\xe3\xff\x7d\x62\xac\x00\x67\xdc\x61\x91\x12\xff\x07\x55\x95\x46\x75\x9d\xd7\x5c\x11\xd2\x71\xc0\x61\x60\xd7\xa5\xdb\x6d\xa4\x4e\x86\xef\x8b\xc7\x83\xd7\x78\x7e\xf3\xee\x66\x8a\xaf\x81\x10\xda\xc5\xd8\x84\xd8\xa1\xc4\xc3\xe2\xaa\x5f\x01\x00\x00\xff\xff\x1c\x00\xaf\xbf\x9c\x03\x00\x00")

func assetsEntrypointPythonBytes() ([]byte, error) {
	return bindataRead(
		_assetsEntrypointPython,
		"assets/entrypoint-python",
	)
}

func assetsEntrypointPython() (*asset, error) {
	bytes, err := assetsEntrypointPythonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/entrypoint-python", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/Dockerfile": assetsDockerfile,
	"assets/Dockerfile-go": assetsDockerfileGo,
	"assets/Dockerfile-python": assetsDockerfilePython,
	"assets/entrypoint": assetsEntrypoint,
	"assets/entrypoint-go": assetsEntrypointGo,
	"assets/entrypoint-python": assetsEntrypointPython,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"assets": &bintree{nil, map[string]*bintree{
		"Dockerfile": &bintree{assetsDockerfile, map[string]*bintree{}},
		"Dockerfile-go": &bintree{assetsDockerfileGo, map[string]*bintree{}},
		"Dockerfile-python": &bintree{assetsDockerfilePython, map[string]*bintree{}},
		"entrypoint": &bintree{assetsEntrypoint, map[string]*bintree{}},
		"entrypoint-go": &bintree{assetsEntrypointGo, map[string]*bintree{}},
		"entrypoint-python": &bintree{assetsEntrypointPython, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

