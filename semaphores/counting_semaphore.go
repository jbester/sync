// Copyright 2018 Jeffrey Bester <jbester@gmail.com>
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

package semaphores

import (
	"sync"
	"sync/atomic"
	"time"
)

type countingSemaphore struct {
	signal  chan empty
	lock    *sync.Mutex
	current int32
	waiting int32
	max     int32
}

// Create a counting semaphore.  The give operation increments the semaphore.
// A take operation decrements the semaphore.
func MakeCountingSemaphore(initial int32, max int32) Semaphore {
	var semaphore = &countingSemaphore{
		signal:  make(chan empty, 1),
		lock:    &sync.Mutex{},
		current: initial,
		max:     max,
	}
	if initial > max {
		panic("semaphore create with initial larger than maximum")
	}
	return semaphore
}

func (semaphore *countingSemaphore) tryAcquire() bool {
	var ok = false
	var count = semaphore.Count()
	if count > 0 {
		ok = atomic.CompareAndSwapInt32(&semaphore.current, count, count-1)
		if ok && count == semaphore.max {
			semaphore.notify()
		}
	}
	return ok
}

func (semaphore *countingSemaphore) tryGive() bool {
	var ok = false
	var count = semaphore.Count()
	if count < semaphore.max {
		ok = atomic.CompareAndSwapInt32(&semaphore.current, count, count+1)
		if ok && count == 0 {
			semaphore.notify()
		}
	}
	return ok
}

func (semaphore *countingSemaphore) notify() {
	var numWaiting = atomic.LoadInt32(&semaphore.waiting)
	if numWaiting > 0 {
		// swap channel
		semaphore.lock.Lock()
		var old, new chan empty
		old, new = semaphore.signal, make(chan empty, 1)
		semaphore.signal = new
		semaphore.lock.Unlock()

		// wake everyone
		close(old)
	}
}

func (semaphore *countingSemaphore) wait() {
	atomic.AddInt32(&semaphore.waiting, 1)
	<-semaphore.signal
	atomic.AddInt32(&semaphore.waiting, -1)
}

func (semaphore *countingSemaphore) timedWait(timeout *time.Duration) bool {
	atomic.AddInt32(&semaphore.waiting, 1)
	var ok = false
	var start = time.Now()
	select {
	case <-semaphore.signal:
		// decrement timeout by time elapsed
		var timeElapsed = time.Now().Sub(start)
		if timeElapsed < *timeout {
			*timeout -= timeElapsed
		} else {
			*timeout = 0
		}
		ok = true

	case <-time.After(*timeout):
	}
	atomic.AddInt32(&semaphore.waiting, -1)
	return ok
}

func (semaphore *countingSemaphore) Take() {
	var ok = false
	for !ok {
		// if empty wait
		if semaphore.IsEmpty() {
			semaphore.wait()
		}

		ok = semaphore.tryAcquire()
	}
}

func (semaphore *countingSemaphore) TryTake(timeout time.Duration) bool {
	var ok = false
	for !ok {
		// if empty wait
		if semaphore.IsEmpty() {
			if !semaphore.timedWait(&timeout) {
				return false
			}
		}

		ok = semaphore.tryAcquire()
	}
	return ok
}

func (semaphore *countingSemaphore) Give() bool {
	var ok = false
	// if not full give
	if !semaphore.IsFull() {
		ok = semaphore.tryGive()
	}
	return ok
}

func (semaphore *countingSemaphore) IsEmpty() bool {
	return semaphore.Count() == 0
}

func (semaphore *countingSemaphore) IsFull() bool {
	return semaphore.Count() == semaphore.max
}

func (semaphore *countingSemaphore) Count() int32 {
	return atomic.LoadInt32(&semaphore.current)
}
