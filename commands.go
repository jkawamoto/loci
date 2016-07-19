//
// commands.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package main

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
	"github.com/urfave/cli"
)

// GlobalFlags defines global flags.
var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "name, n",
		Usage: "creating a container named `NAME` to run tests.",
		// If name is given, continer will not be deleted.
	},
	cli.StringFlag{
		Name:  "tag, t",
		Usage: "creating an image named `TAG`.",
	},
}

// Commands defines sub-commands.
var Commands = []cli.Command{}

// CommandNotFound prints an error message when a given command is not supported.
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(
		os.Stderr, chalk.Red.Color("%s: '%s' is not a %s command. See '%s --help'."),
		c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
