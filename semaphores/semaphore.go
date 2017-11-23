// Package semaphores provides a go implementation of binary and counting semaphores.
// The underlying implementation is built on the channel primitive.
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
