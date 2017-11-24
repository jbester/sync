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

package startgroup

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	<-time.After(time.Millisecond)
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

func Benchmark_Release(b *testing.B) {
	var sg = MakeStartGroup()
	for n := 0; n < b.N; n++ {
		sg.Release()
	}
}
