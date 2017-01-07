//
// command/run_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import "testing"

func TestGetRepository(t *testing.T) {

	expected := "github.com/jkawamoto/loci"

	if res := getRepository("ssh://git@github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("ssh://git:pw@github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("http://github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("http://username@github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("http://username:pw@github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("https://github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("https://username@github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

	if res := getRepository("https://username:pw@github.com/jkawamoto/loci.git"); res != expected {
		t.Error("Wrong repository path:", res)
	}

}
