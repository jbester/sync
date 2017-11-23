// Package events provides a single synchronization primitive the Event.  The event is used to notify the occurrence of
// a condition to routines.  Each event has two states - up and down.  Up indicates the condition has occurred.
//
// Multiple routines can wait on a condition.   All routines unblock once the condition occurs.
// A routine that waits on a condition that has already occurred will not block.
//
// The event primitive is similar to the event in the pSOS or ARINC 653 API sets.
package events

import (
	"sync/semaphores"
	"sync/startgroup"
	"time"
)

type Event struct {
	state      semaphores.Semaphore
	notifyList *startgroup.StartGroup
}

// Creates an event object for use by any routine.  Upon creation the event is set to the down state.
func MakeEvent() *Event {
	return &Event{semaphores.MakeBinarySemaphore(false), startgroup.MakeStartGroup()}
}

//  Set request sets the specified event to the up state. All the routines waiting
//  on the event will stop waiting.  Any routine attempting to wait will not block.
func (evt *Event) Set() {
	if !evt.IsSet() {
		evt.state.Give()
		evt.notifyList.Release()
	}
}

//  Checks if the specified event is in up state.
func (evt *Event) IsSet() bool {
	return evt.state.IsFull()
}

//  Resets the specified event to the down state.   Any routine attempting to wait
//  will block until the event is set.
func (evt *Event) Reset() {
	if evt.IsSet() {
		evt.state.Take()
	}
}

//  Wait for the event to be in the up state.  Any routine that attemps to wait on an event
//  already in the up state will not block.
func (evt *Event) Wait() {
	if evt.IsSet() {
		return
	} else {
		evt.notifyList.Wait()
	}
}

//  Wait for the vent to be in the up state.  Any routine waiting on an event
//  in the up state will not block.   Returns true if the event is in the up state and
//  false if a timeout occurred.
func (evt *Event) TryWait(timeout time.Duration) bool {
	var ok bool
	if evt.IsSet() {
		ok = true
	} else {
		ok = evt.notifyList.TryWait(timeout)
	}

	return ok
}
