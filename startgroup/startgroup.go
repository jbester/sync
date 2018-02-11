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
	"time"
)

type empty struct{}

// A StartGroup provides a mechanism for a collection of goroutines to wait for a release event.
// When released, all blocked routines simultaneously.
//
// A typical use is when multiple routines need to know when a resource is available but do
// not need exclusive access to the resource.
type StartGroup interface {
	//  Release all Waiting goroutines
	Release()

	// Wait for a release event
	Wait()

	// Wait for a release event for up to a timeout
	TimedWait(timeout time.Duration) bool
}

type startGroup struct {
	lock       *sync.RWMutex
	notifyList *sync.WaitGroup
}

//  Create a StartGroup.
func MakeStartGroup() StartGroup {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	return &startGroup{lock: &sync.RWMutex{}, notifyList: wg}
}

func (group *startGroup) Release() {
	// create a wg waitgroup
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	// mutex is to prevent multiple goroutines from trying to release simultaneously
	group.lock.Lock()
	// swap the wg waitgroup in for the old one
	var old = group.notifyList
	group.notifyList = wg
	group.lock.Unlock()

	// release the old one
	old.Done()
}

func (group *startGroup) Wait() {
	group.lock.RLock()
	var waitList = group.notifyList
	group.lock.RUnlock()
	waitList.Wait()
}

func (group *startGroup) TimedWait(timeout time.Duration) bool {
	group.lock.RLock()
	var waitList = group.notifyList
	group.lock.RUnlock()
	var ch = make(chan empty)
	defer close(ch)
	go func() {
		waitList.Wait()
		ch <- empty{}
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}
