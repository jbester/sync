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

package semaphores

import (
	"time"
)

type countingSemaphore struct {
	items chan empty
	max   int
}

// Create a counting semaphore.  The give operation increments the semaphore.
// A take operation decrements the semaphore.
func MakeCountingSemaphore(initial int, max int) Semaphore {
	var semaphore = &countingSemaphore{
		items: make(chan empty, max),
		max:   max,
	}
	if initial > max {
		panic("semaphore create with initial larger than maximum")
	}

	for initial > 0 {
		semaphore.items <- empty{}
		initial--
	}
	return semaphore
}

func (semaphore *countingSemaphore) Take() {
	<-semaphore.items
}

func (semaphore *countingSemaphore) TryTake(duration time.Duration) bool {
	select {
	case <-semaphore.items:
		return true
	case <-time.After(duration):
		return false
	}
}

func (semaphore *countingSemaphore) Give() bool {
	var ok = false
	if !semaphore.IsFull() {
		semaphore.items <- empty{}
		ok = true
	}
	return ok
}

func (semaphore *countingSemaphore) IsEmpty() bool {
	return len(semaphore.items) == 0
}

func (semaphore *countingSemaphore) IsFull() bool {
	return len(semaphore.items) == semaphore.max
}
func (semaphore *countingSemaphore) Count() int {
	return len(semaphore.items)
}
