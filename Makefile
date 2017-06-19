#
# Makefile
#
# Copyright (c) 2016 Junpei Kawamoto
#
# This software is released under the MIT License.
#
# http://opensource.org/licenses/mit-license.php
#
VERSION = snapshot
GHRFLAGS =

default: build

.PHONY: asset
asset:
	go-bindata -pkg command -o command/assets.go -nometadata assets

.PHONY: build
build: asset
	goxc -os="darwin linux windows" -d=pkg -pv=$(VERSION)

.PHONY: release
release:
	ghr  -u jkawamoto $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)

.PHONY: test
test: asset
	go test -v ./...

.PHONY: local
local: asset
	go build
	go install

.PHONY: get-deps
get-deps:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -d -t -v .
