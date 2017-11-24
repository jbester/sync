package startgroup

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type StartGroupTestSuite struct {
	suite.Suite
	startGroup *StartGroup
	waitGroup  *sync.WaitGroup
}

func (suite *StartGroupTestSuite) SetupTest() {
	suite.startGroup = MakeStartGroup()
	suite.waitGroup = &sync.WaitGroup{}
}

type callback func()

//  Spawn a routine to wait on the startgroup
func (suite *StartGroupTestSuite) asyncWait(onWaitComplete callback) {
	go func() {
		suite.waitGroup.Add(1)
		defer suite.waitGroup.Done()
		suite.startGroup.Wait()
		onWaitComplete()
	}()
	runtime.Gosched()
}

//  Spawn a routine to trywait on the startgroup
func (suite *StartGroupTestSuite) asyncTryWait(timeout time.Duration, onWaitComplete callback) {
	go func() {
		suite.waitGroup.Add(1)
		defer suite.waitGroup.Done()
		if suite.startGroup.TryWait(timeout) {
			onWaitComplete()
		}
	}()
	runtime.Gosched()
}

func (suite *StartGroupTestSuite) Test_Wait() {
	var done = false
	suite.asyncWait(func() {
		done = true
	})
	suite.startGroup.Release()
	suite.waitGroup.Wait()
	assert.True(suite.T(), done)
}

func (suite *StartGroupTestSuite) Test_TryWaitSuccess() {
	var done = false
	suite.asyncTryWait(time.Second, func() {
		done = true
	})
	suite.startGroup.Release()
	suite.waitGroup.Wait()
	assert.True(suite.T(), done)
}

func (suite *StartGroupTestSuite) Test_TryWaitTimeout() {
	var done = false
	suite.asyncTryWait(time.Millisecond*5, func() {
		done = true
	})
	suite.waitGroup.Wait()
	assert.False(suite.T(), done)
}

func (suite *StartGroupTestSuite) Test_ReleaseGroup() {
	var done int32 = 0
	suite.asyncWait(func() {
		atomic.AddInt32(&done, 1)
	})
	suite.asyncWait(func() {
		atomic.AddInt32(&done, 1)
	})
	suite.startGroup.Release()
	suite.waitGroup.Wait()
	assert.Equal(suite.T(), int32(2), done)
}

func TestStartGroupTestSuite(t *testing.T) {
	suite.Run(t, new(StartGroupTestSuite))
}
