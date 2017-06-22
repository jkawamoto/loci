//
// command/errset_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"fmt"
	"io"
	"testing"
)

func TestErrorSet(t *testing.T) {

	errs := NewErrorSet()
	if errs.Size() != 0 {
		t.Error("New ErrorSet contains some error already.")
	}

	errs.Add("eof", io.EOF)
	if errs.Size() != 1 {
		t.Error("ErrorSet returns a wrong size.")
	}

	var ret []error
	ret = errs.GetList()
	if len(ret) != 1 {
		t.Error("Length of a slice returned from an ErrorSet is wrong.")
	} else if ret[0] != io.EOF {
		t.Error("Slice returned from an ErrorSet consists of wrong errors.")
	}

	another := fmt.Errorf("Another error")
	errs.Add("another", another)
	ret = errs.GetList()
	if len(ret) != 2 {
		t.Error("Length of a slice returned from an ErrorSet is wrong.")
	} else if ret[0] != another || ret[1] != io.EOF {
		t.Error("Slice returned from an ErrorSet consists in a wrong order.")
	}

}
