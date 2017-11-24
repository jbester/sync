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
type StartGroup struct {
	lock       *sync.Mutex
	notifyList chan empty
}

//  Create a StartGroup.
func MakeStartGroup() *StartGroup {
	return &StartGroup{&sync.Mutex{}, make(chan empty)}
}

//  Release all Waiting goroutines
func (group *StartGroup) Release() {
	// mutex is to prevent multiple goroutines from trying to release simultaneously
	group.lock.Lock()
	// store off the current channel
	var ch = group.notifyList

	// replace it - in case the listener immediately waits again
	group.notifyList = make(chan empty)

	// close it - this will wake up all waiting goroutines
	close(ch)
	group.lock.Unlock()
}

// Wait for a release event
func (group *StartGroup) Wait() {
	<-group.notifyList
}

// Wait for a release event or timeout.   If the release event occurs,
// return true.
func (group *StartGroup) TryWait(timeout time.Duration) bool {
	select {
	case <-group.notifyList:
		return true

	case <-time.After(timeout):
		return false
	}
}
