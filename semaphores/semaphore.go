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

// Package semaphores provides a go implementation of binary and counting semaphores.
// The underlying implementation is built on the channel primitive.  As such, it doesn't
// offer any advantages over using a channel except for readability.
package semaphores

import "time"

// Semaphore interface.
type Semaphore interface {
	// Take (decrement) a semaphore.  Routine will block until the semaphore is available.
	Take()

	// Take (decrement) a semaphore.  Routine will block until the timeout has occurred
	// or the semaphore becomes available
	TryTake(timeout time.Duration) bool

	// Release (increment) the semaphore.  Returns false if semaphore is full.
	Give() bool

	// Test if the semaphore is full.
	IsFull() bool

	// Test if the semaphore is empty.
	IsEmpty() bool

	//  Returns the count of a semaphore.
	Count() int
}
