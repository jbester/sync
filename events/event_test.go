package events

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime"
	"sync"
	"sync/semaphores"
	"testing"
	"time"
)

type TestEventSuite struct {
	suite.Suite
	evt *Event
	wg  *sync.WaitGroup
}

func (suite *TestEventSuite) SetupTest() {
	suite.evt = MakeEvent()
	suite.wg = &sync.WaitGroup{}
}

type callback func()

//  Spawn a routine to wait on the event
func (suite *TestEventSuite) asyncWait(onWaitComplete callback) {
	go func() {
		suite.wg.Add(1)
		defer suite.wg.Done()
		suite.evt.Wait()
		onWaitComplete()
	}()
	runtime.Gosched()
}

//  Spawn a routine to trywait on the event
func (suite *TestEventSuite) asyncTryWait(timeout time.Duration, onWaitComplete callback) {
	go func() {
		suite.wg.Add(1)
		defer suite.wg.Done()
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
	suite.wg.Wait()
	assert.True(suite.T(), eventReceived)
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

func (suite *TestEventSuite) Test_TryWait() {
	var eventReceived = false
	suite.asyncTryWait(time.Millisecond, func() {
		eventReceived = true
	})
	suite.wg.Wait()
	assert.False(suite.T(), eventReceived)
}

func (suite *TestEventSuite) Test_TryWaitWhenSet() {
	var eventReceived = false
	suite.evt.Set()
	suite.asyncTryWait(time.Millisecond, func() {
		eventReceived = true
	})
	suite.wg.Wait()
	assert.True(suite.T(), eventReceived)
}

func (suite *TestEventSuite) Test_WaitWhenReset() {
	var eventReceived = false
	suite.evt.Set()
	suite.evt.Reset()
	suite.asyncTryWait(time.Millisecond, func() {
		eventReceived = true
	})
	assert.False(suite.T(), eventReceived)
}

func (suite *TestEventSuite) Test_WaitMultiple() {
	var eventCount = semaphores.MakeCountingSemaphore(0, 2)
	suite.asyncWait(func() {
		eventCount.Give()
	})
	suite.asyncWait(func() {
		eventCount.Give()
	})
	suite.evt.Set()
	suite.wg.Wait()
	assert.True(suite.T(), eventCount.IsFull())
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(TestEventSuite))
}
