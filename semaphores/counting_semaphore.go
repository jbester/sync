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
