//
// commands.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

// GlobalFlags defines global flags.
var GlobalFlags = []cli.Flag{
	// TODO: Remote docker host.
	// TODO: delete name option because lots of containers will be created.
	cli.StringFlag{
		Name: "name, n",
		Usage: "creating a container named `NAME` to run tests. " +
			"If name is given, container will not be deleted.",
	},
	cli.StringFlag{
		Name:  "select, s",
		Usage: "select specific runtime `VERSION` where tests running on.",
	},
	cli.StringFlag{
		Name:  "tag, t",
		Usage: "specify a `TAG` name of the docker image to be build.",
	},
	cli.IntFlag{
		Name:  "max-processors",
		Usage: "Max processors used to run tests.",
		Value: runtime.NumCPU() - 2,
	},
	cli.BoolFlag{
		Name:  "log, l",
		Usage: "store logging information to files.",
	},
	cli.StringFlag{
		Name:  "base, b",
		Usage: "use image `TAG` as the base image.",
		Value: "ubuntu:latest",
	},
	cli.StringFlag{
		Name:   "apt-proxy",
		Usage:  "`URL` for a proxy server of apt repository.",
		EnvVar: "APT_PROXY",
	},
	cli.StringFlag{
		Name:   "pypi-proxy",
		Usage:  "`URL` for a proxy server of pypi repository.",
		EnvVar: "PYPI_PROXY",
	},
	cli.StringFlag{
		Name:   "http-proxy",
		Usage:  "`URL` for a http proxy server.",
		EnvVar: "HTTP_PROXY",
	},
	cli.StringFlag{
		Name:   "https-proxy",
		Usage:  "`URL` for a https proxy server.",
		EnvVar: "HTTPS_PROXY",
	},
	cli.StringFlag{
		Name:   "no-proxy",
		Usage:  "Comma separated URL `LIST` for which proxies won't be used.",
		EnvVar: "NO_PROXY",
	},
	cli.BoolFlag{
		Name:  "no-build-cache",
		Usage: "Do not use cache when building the image.",
	},
	cli.BoolFlag{
		Name:  "no-color",
		Usage: "Omit to print color codes.",
	},
}

// Commands defines sub-commands.
var Commands = []cli.Command{}

// CommandNotFound prints an error message when a given command is not supported.
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(
		os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
