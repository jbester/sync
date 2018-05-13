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
	"testing"
	"time"

	"sync/atomic"

	"sync"

	"github.com/stretchr/testify/assert"
)

func Test_CreateInvalid(t *testing.T) {
	assert.Panics(t, func() {
		var semaphore = MakeCountingSemaphore(2, 1)
		semaphore.Give()
		assert.True(t, semaphore.IsFull())
	})
}

func Test_CountingLock(t *testing.T) {
	var semaphore = MakeCountingSemaphore(0, 1)
	semaphore.Give()
	assert.True(t, semaphore.IsFull())
}

func Test_CountingLockUnlock(t *testing.T) {
	var semaphore = MakeCountingSemaphore(0, 1)
	semaphore.Give()
	semaphore.Take()
	assert.True(t, semaphore.IsEmpty())
}

func Test_CountingGiveFull(t *testing.T) {
	var semaphore = MakeCountingSemaphore(1, 1)
	assert.False(t, semaphore.Give())
}

func Test_CountingTakeEmpty(t *testing.T) {
	var semaphore = MakeCountingSemaphore(0, 1)
	assert.False(t, semaphore.TryTake(time.Millisecond))
}

func Test_CountingTryTake(t *testing.T) {
	var semaphore = MakeCountingSemaphore(1, 1)
	assert.True(t, semaphore.TryTake(time.Millisecond))
}

func Test_CountingGetCount(t *testing.T) {
	var semaphore = MakeCountingSemaphore(3, 5)
	assert.Equal(t, int32(3), semaphore.Count())
	semaphore.Give()
	assert.Equal(t, int32(4), semaphore.Count())
}

func Test_CountingTakeEmptyMultiple(t *testing.T) {
	var semaphore = MakeCountingSemaphore(0, 1)
	var takes int32 = 0
	var started = &sync.WaitGroup{}
	var done = &sync.WaitGroup{}
	started.Add(2)
	done.Add(1) // only one expected to finish

	var worker = func() {
		// mark started
		started.Done()
		// take
		semaphore.Take()
		// increment number of takes
		atomic.AddInt32(&takes, 1)
		// mark done
		done.Done()
	}

	go worker()
	go worker()

	// wait for all coroutines start
	started.Wait()
	// verify precondition
	assert.Equal(t, int32(0), takes)
	// give the semaphore
	semaphore.Give()
	// wait for one of the coroutines finishes
	done.Wait()
	// verify count
	assert.Equal(t, int32(1), takes)
}

func Test_CountingTimedTakeEmptyMultiple(t *testing.T) {
	var semaphore = MakeCountingSemaphore(0, 1)
	var takes int32 = 0
	var started = &sync.WaitGroup{}
	var done = &sync.WaitGroup{}
	started.Add(3)
	done.Add(3)
	var worker = func() {
		// mark started
		started.Done()
		var start = time.Now()
		const timeout = time.Millisecond * 100
		// try take
		if semaphore.TryTake(timeout) {
			// increment if successful
			atomic.AddInt32(&takes, 1)
		}
		var delta = time.Now().Sub(start)
		// check if timeout expired within timeout + 5%
		assert.True(t, delta < timeout+time.Millisecond*5)
		// mark done
		done.Done()
	}
	go worker()
	go worker()
	go worker()
	// wait for all coroutines start
	started.Wait()
	// verify precondition
	assert.Equal(t, int32(0), takes)
	// give the semaphore
	semaphore.Give()
	<-time.After(50 * time.Millisecond)
	// give the semaphore
	semaphore.Give()
	// wait for one of the coroutines finishes
	done.Wait()
	// verify takes
	assert.Equal(t, int32(2), takes)
}
