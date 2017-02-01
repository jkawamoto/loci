//
// command/decorator.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import "github.com/ttacon/chalk"

// TextDecorator provides functions to decorate strings.
type TextDecorator struct {
	Bold func(string) string
}

func noDecoration(str string) string {
	return str
}

// NewDecorator returns a decorator.
func NewDecorator() *TextDecorator {

	return &TextDecorator{
		Bold: chalk.Bold.TextStyle,
	}

}

// NewNoopDecorator returns a decorator which doesn't decorate any strings.
func NewNoopDecorator() *TextDecorator {

	return &TextDecorator{
		Bold: noDecoration,
	}

}
