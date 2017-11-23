package semaphores

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_BinaryLock(t *testing.T) {
	var semaphore = MakeBinarySemaphore(false)
	semaphore.Give()
	assert.True(t, semaphore.IsFull())
}

func Test_BinaryLockUnlock(t *testing.T) {
	var semaphore = MakeBinarySemaphore(false)
	semaphore.Give()
	semaphore.Take()
	assert.True(t, semaphore.IsEmpty())
}

func Test_BinaryGiveFull(t *testing.T) {
	var semaphore = MakeBinarySemaphore(true)
	assert.False(t, semaphore.Give())
}

func Test_BinaryTakeEmpty(t *testing.T) {
	var semaphore = MakeBinarySemaphore(false)
	assert.False(t, semaphore.TryTake(time.Millisecond))
}
