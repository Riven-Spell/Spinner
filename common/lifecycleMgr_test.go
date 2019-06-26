package common_test

import (
	"errors"
	"testing"
	"time"

	chk "gopkg.in/check.v1"

	. "github.com/virepri/Spinner/common"
)

func Test(t *testing.T) { chk.TestingT(t) }

type LcmTestSuite struct{}

var _ = chk.Suite(&LcmTestSuite{})

func (s *LcmTestSuite) TestShutdown(c *chk.C) {
	lcm := GetLifecycleManager()

	killChan := make(chan bool, 1)

	// Create primary suicide routine
	lcm.CreateRoutine(func() {
		<-killChan
		c.Log("Kill initiated")
		lcm.Shutdown(0)
	})

	for i := 1; i <= 5; i++ {
		// Create domino suicide routines
		lcm.CreateRoutine(func() {
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
		<-time.After(time.Millisecond * 100)
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
		lcm.CreateRoutine(func() {
			select {
			case <-killChan:
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
		<-time.After(time.Millisecond * 100)
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
		<-killChan
		c.Log("Suicide initiated")
		lcm.Suicide(suicideErr, 1)
	})

	for i := 1; i <= 5; i++ {
		// Create domino suicide routines
		lcm.CreateRoutine(func() {
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
		<-time.After(time.Millisecond * 100)
		killChan <- true
	}()

	exitCode, errorOut := lcm.SurrenderControl()

	c.Check(len(errorOut), chk.Equals, 1)
	c.Check(*exitCode, chk.Equals, 1)
}

type SpoofWriter struct {
	LastInput string
}

func (s *SpoofWriter) Write(p []byte) (n int, err error) {
	s.LastInput = string(p)
	return len(p), nil
}

func (s *LcmTestSuite) TestLogging(c *chk.C) {
	lcm := GetLifecycleManager()
	sw := SpoofWriter{}

	// Set logging destination to a spoof writer for testing
	lcm.SetLogDestination(&sw)

	// Try outputting a log and check if the spoof writer received it.
	lcm.Log("testing logs", ELogLevel.Fatal())
	c.Assert(sw.LastInput, chk.Equals, "FATAL: testing logs\n")

	// Set logging level to Warning (Level 1)
	lcm.SetMinimumLogLevel(ELogLevel.Warning())

	// Test under the minimum.
	lcm.Log("Do not log me", ELogLevel.Information())
	c.Assert(sw.LastInput, chk.Equals, "FATAL: testing logs\n")

	// Test on & over the minimum
	lcm.Log("I should be logged.", ELogLevel.Warning())
	c.Assert(sw.LastInput, chk.Equals, "WARNING: I should be logged.\n")
	lcm.Log("I should also be logged.", ELogLevel.Fatal())
	c.Assert(sw.LastInput, chk.Equals, "FATAL: I should also be logged.\n")

	// Test ELogLevel.None()
	lcm.SetMinimumLogLevel(ELogLevel.None())
	lcm.Log("BIG FATAL ERROR", ELogLevel.Fatal()) // Fatal is the highest severity.
	c.Assert(sw.LastInput, chk.Equals, "FATAL: I should also be logged.\n")
}
