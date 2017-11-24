package events

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TestEventSuite struct {
	suite.Suite
	evt       *Event
	waitGroup *sync.WaitGroup
}

func (suite *TestEventSuite) SetupTest() {
	suite.evt = MakeEvent()
	suite.waitGroup = &sync.WaitGroup{}
}

type callback func()

//  Spawn a routine to wait on the event
func (suite *TestEventSuite) asyncWait(onWaitComplete callback) {
	go func() {
		suite.waitGroup.Add(1)
		defer suite.waitGroup.Done()
		suite.evt.Wait()
		onWaitComplete()
	}()
	runtime.Gosched()
}

//  Spawn a routine to trywait on the event
func (suite *TestEventSuite) asyncTryWait(timeout time.Duration, onWaitComplete callback) {
	go func() {
		suite.waitGroup.Add(1)
		defer suite.waitGroup.Done()
		if suite.evt.TryWait(timeout) {
			onWaitComplete()
		}
	}()
	runtime.Gosched()
}

//  Test that the routine wakes up when set
func (suite *TestEventSuite) Test_Wait() {
	var eventReceived = false
	suite.asyncWait(func() {
		eventReceived = true
	})
	suite.evt.Set()
	suite.waitGroup.Wait()
	assert.True(suite.T(), eventReceived)
	assert.True(suite.T(), suite.evt.IsSet())
}

//  Test that when set it remains set
func (suite *TestEventSuite) Test_Set() {
	suite.asyncWait(func() {
	})
	assert.True(suite.T(), suite.evt.Set())
	// event is already set
	assert.False(suite.T(), suite.evt.Set())
	// wait until the event has been read
	suite.waitGroup.Wait()
	// verify event remains set
	assert.True(suite.T(), suite.evt.IsSet())
}

//  Test that when reset the event is cleared
func (suite *TestEventSuite) Test_Reset() {
	assert.True(suite.T(), suite.evt.Set())
	assert.True(suite.T(), suite.evt.Reset())
	assert.False(suite.T(), suite.evt.IsSet())
}

//  Test that the routine doesn't block when set before wait
func (suite *TestEventSuite) Test_WaitWhenAlreadySet() {
	var eventReceived = false
	suite.evt.Set()
	suite.asyncWait(func() {
		eventReceived = true
	})
	assert.True(suite.T(), eventReceived)
}

// Test try wait will trigger a timeout
func (suite *TestEventSuite) Test_TryWaitTimeout() {
	var eventReceived = false
	suite.asyncTryWait(time.Millisecond, func() {
		eventReceived = true
	})
	suite.waitGroup.Wait()
	assert.False(suite.T(), eventReceived)
}

// Test try wait will return immediately if already set
func (suite *TestEventSuite) Test_TryWaitWhenSet() {
	var eventReceived = false
	suite.evt.Set()
	suite.asyncTryWait(time.Millisecond, func() {
		eventReceived = true
	})
	assert.True(suite.T(), eventReceived)
}

// Test try wait will wait on a reset event
func (suite *TestEventSuite) Test_WaitWhenReset() {
	var eventReceived = false
	suite.evt.Set()
	suite.evt.Reset()
	suite.asyncTryWait(time.Millisecond, func() {
		eventReceived = true
	})
	assert.False(suite.T(), eventReceived)
}

// Test that the event will wake up multiple routines
func (suite *TestEventSuite) Test_WaitMultiple() {
	var eventCount int32 = 0
	suite.asyncWait(func() {
		atomic.AddInt32(&eventCount, 1)
	})
	suite.asyncWait(func() {
		atomic.AddInt32(&eventCount, 1)
	})
	suite.evt.Set()
	suite.waitGroup.Wait()
	assert.Equal(suite.T(), int32(2), eventCount)
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(TestEventSuite))
}
