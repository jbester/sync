// Copyright 2017 Jeffrey Bester <jbester@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
// to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of
// the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS
// OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package startgroup provides a single synchronization primitive the StartGroup.
//
package startgroup

import (
	"sync"
)

type empty struct{}

// A StartGroup provides a mechanism for a collection of goroutines to wait for a release event.
// When released, all blocked routines simultaneously.
//
// A typical use is when multiple routines need to know when a resource is available but do
// not need exclusive access to the resource.
type StartGroup struct {
	lock       *sync.Mutex
	notifyList *sync.WaitGroup
}

//  Create a StartGroup.
func MakeStartGroup() *StartGroup {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	return &StartGroup{lock: &sync.Mutex{}, notifyList: wg}
}

//  Release all Waiting goroutines
func (group *StartGroup) Release() {
	// mutex is to prevent multiple goroutines from trying to release simultaneously
	group.lock.Lock()
	// create a new waitgroup and swap it
	var new = &sync.WaitGroup{}
	var old = group.notifyList
	new.Add(1)
	group.notifyList = new
	// release the old one
	old.Done()
	group.lock.Unlock()
}

// Wait for a release event
func (group *StartGroup) Wait() {
	group.notifyList.Wait()
}
