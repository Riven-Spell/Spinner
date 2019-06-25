package common_test

import (
	"errors"
	"testing"
	"time"
	
	chk "gopkg.in/check.v1"
	
	. "github.com/virepri/Spinner/common"
)

func Test(t *testing.T) {chk.TestingT(t)}

type LcmTestSuite struct{}

var _ = chk.Suite(&LcmTestSuite{})

func (s *LcmTestSuite) TestShutdown(c *chk.C) {
	lcm := GetLifecycleManager()
	
	killChan := make(chan bool, 1)
	
	// Create primary suicide routine
	lcm.CreateRoutine(func() {
		<- killChan
		c.Log("Kill initiated")
		lcm.Shutdown(0)
	})
	
	for i := 1; i <= 5; i++ {
		// Create domino suicide routines
		lcm.CreateRoutine(func () {
			select {
				case code := <-lcm.WatchForShutdown():
					c.Log("Received shutdown call.")
					lcm.Shutdown(code)
					return
				case <-time.After(5 * time.Second):
					c.Log("TestShutdown timed out.")
					c.Fail()
			}
		})
	}
	
	// Send the kill call in 0.1 seconds
	go func() {
		<- time.After(time.Millisecond * 100)
		killChan <- true
	}()
	
	exitCode, _ := lcm.SurrenderControl()
	if exitCode == nil {
		c.FailNow()
	}
	
	c.Assert(*exitCode, chk.Equals, 0)
}

func (s *LcmTestSuite) TestAllExit(c *chk.C) {
	lcm := GetLifecycleManager()
	
	killChan := make(chan bool, 1)
	
	for i := 1; i <= 5; i++ {
		// Create domino suicide routines
		lcm.CreateRoutine(func () {
			select {
			case <- killChan:
				killChan <- true
				c.Log("Received kill call.")
				return
			case <-time.After(5 * time.Second):
				c.Log("TestShutdown timed out.")
				c.Fail()
			}
		})
	}
	
	// Send the kill call in 0.1 seconds
	go func() {
		<- time.After(time.Millisecond * 100)
		killChan <- true
	}()
	
	exitCode, _ := lcm.SurrenderControl()
	c.Assert(exitCode, chk.Equals, (*int)(nil))
}

func (s *LcmTestSuite) TestSuicide(c *chk.C) {
	lcm := GetLifecycleManager()
	
	killChan := make(chan bool, 1)
	
	suicideErr := errors.New("committing suicide")
	
	// Create primary suicide routine
	lcm.CreateRoutine(func() {
		<- killChan
		c.Log("Suicide initiated")
		lcm.Suicide(suicideErr, 1)
	})
	
	for i := 1; i <= 5; i++ {
		// Create domino suicide routines
		lcm.CreateRoutine(func () {
			select {
			case code := <-lcm.WatchForShutdown():
				c.Log("Received shutdown call.")
				lcm.Shutdown(code)
				return
			case <-time.After(5 * time.Second):
				c.Log("TestShutdown timed out.")
				c.Fail()
			}
		})
	}
	
	// Send the kill call in 0.1 seconds
	go func() {
		<- time.After(time.Millisecond * 100)
		killChan <- true
	}()
	
	exitCode, errorOut := lcm.SurrenderControl()
	
	c.Check(len(errorOut), chk.Equals, 1)
	c.Check(*exitCode, chk.Equals, 1)
}