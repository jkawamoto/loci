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

default: build

.PHONY: asset
asset:
	go-bindata -pkg util -o command/util/assets.go assets/startup.sh
	go-bindata -pkg command -o command/assets.go asset

.PHONY: build
build: asset
	goxc -os="darwin linux windows" -d=pkg -pv=$(VERSION)
