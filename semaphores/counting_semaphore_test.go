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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	assert.Equal(t, 3, semaphore.Count())
	semaphore.Give()
	assert.Equal(t, 4, semaphore.Count())
}
