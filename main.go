package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "jkawamoto"
	app.Email = ""
	app.Usage = ""

	// app.Flags = GlobalFlags
	// app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Action = Run

	app.Run(os.Args)
}
