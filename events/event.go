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

// Package events provides a single synchronization primitive - the Event.  The event is used to notify that
// a condition has occurred to routines.  Each event has two states - set and unset.  Set state indicates the condition has occurred.
//
// Multiple routines can wait on a condition.   _All_ routines unblock once the event is set to the set state.
// A routine that waits on an event that is already set will not block.
//
// The event primitive is similar to the event in the pSOS or ARINC 653 APIs.
package events

import (
	"sync/atomic"
	"time"

	"bitbucket.org/jbester/sync/startgroup"
)

type Event struct {
	state      int32
	notifyList *startgroup.StartGroup
}

// Creates an event object for use by any routine.  Upon creation the event is set to the unset state.
func MakeEvent() *Event {
	return &Event{state: 0, notifyList: startgroup.MakeStartGroup()}
}

//  Set changes the specified event to the set state. All the routines waiting
//  on the event will stop waiting.  Any subsequent routine attempting to wait will not block until the
//  event is set to the unset state again.
//  Returns true if the event changed.
func (evt *Event) Set() bool {
	var ok = atomic.CompareAndSwapInt32(&evt.state, 0, 1)
	if ok {
		evt.notifyList.Release()
	}
	return ok
}

//  Checks if the specified event is in set state.
func (evt *Event) IsSet() bool {
	return atomic.LoadInt32(&evt.state) == 1
}

//  Resets the specified event to the unset state.   Once reset, any routine attempting to wait
//  will block until the event is set again.  Returns true if the event changed.
func (evt *Event) Reset() bool {
	return atomic.CompareAndSwapInt32(&evt.state, 1, 0)
}

//  Wait for the event to be in the set state.  Any routine that attempts to wait on an event
//  already in the set state will not block.
func (evt *Event) Wait() {
	if evt.IsSet() {
		return
	} else {
		evt.notifyList.Wait()
	}
}

//  Wait for the event to be in the set state up to the given timeout.  Any routine that attempts to wait on an event
//  already in the set state will not block.
func (evt *Event) TimedWait(timeout time.Duration) bool {
	if evt.IsSet() {
		return true
	} else {
		return evt.notifyList.TimedWait(timeout)
	}
}
