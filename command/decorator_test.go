//
// command/decorator_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import "testing"

func TestNewDecorator(t *testing.T) {

	d := NewDecorator()
	if d.Bold == nil {
		t.Error("NewDecorator returns an uninitialized decorator.")
	}

}

func TestNewNoopDecorator(t *testing.T) {

	d := NewNoopDecorator()
	if d.Bold == nil {
		t.Error("NewNoopDecorator returns an uninitialized decorator.")
	}
	if d.Bold("abcdefg") != "abcdefg" {
		t.Error("NewNoopDecorator decorates strings.")
	}

}
