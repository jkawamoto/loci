//
// command/errset.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"sort"
	"sync"
)

// ErrorSet is a goroutine-safe set of errors. Each error in this set has a key
// and GetList function returns a list of errors according to the order of keys.
type ErrorSet struct {
	errors map[string]error
	mutex  sync.Mutex
}

// NewErrorSet creates a new ErrorSet.
func NewErrorSet() *ErrorSet {
	return &ErrorSet{
		errors: make(map[string]error),
	}
}

// Add a new error with a key string.
func (e *ErrorSet) Add(key string, err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.errors[key] = err
}

// GetList returns a list of errors in this set; the list is sorted by the order
// of keys.
func (e *ErrorSet) GetList() []error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	keys := make(sort.StringSlice, 0, len(e.errors)) // Not call Size() here, because of the lock.
	for k := range e.errors {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	res := make([]error, 0, len(e.errors)) // Not call Size() here, because of the lock.
	for _, k := range keys {
		res = append(res, e.errors[k])
	}
	return res
}

// Size returns the number of errors in this set.
func (e *ErrorSet) Size() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return len(e.errors)
}
