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
