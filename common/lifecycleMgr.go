package common

import (
	"errors"
	"fmt"
	"sync"
)

var lcm LifecycleManager

func GetLifecycleManager() LifecycleManager {
	if lcm == nil {
		makeLifecycleManager()
	}
	
	return lcm
}

func makeLifecycleManager() {
	lcm = &lifecycleManager{
		shutdownChan: make(chan int, 10), // There will probably never be more than 1 in shutdownChan or 10 in suicideResults
		suicideResults: make(chan error, 500), // But it's best to be safer than sorry.
	}
}

type LifecycleManager interface {
	// SurrenderControl blocks until Shutdown or Suicide is called.
	SurrenderControl() (*int, []error)
	
	// Shutdown gracefully shuts down the app.
	Shutdown(exitCode int)
	// Suicide gracefully shuts down the app and displays an error.
	Suicide(err error, exitCode int)
	// WatchForShutdown returns the underlying shutdownChan. Useful inside select statements.
	WatchForShutdown() chan int
	
	// CreateRoutine adds a waiter to the underlying WaitGroup.
	// This is intended for primary routines only. Child routines should use their own suicide methods.
	CreateRoutine(func())
	
	// TODO: Add logging to lcm
}

type lifecycleManager struct {
	routines     sync.WaitGroup
	shutdownChan chan int
	suicideResults chan error
}

func (lcm *lifecycleManager) SurrenderControl() (*int, []error) {
	lcm.routines.Wait()
	
	errorOut := make([]error, len(lcm.suicideResults))
	errorIndex := 0
	for len(lcm.suicideResults) != 0 {
		err := <- lcm.suicideResults
		
		if err != nil {
			fmt.Println("================ START FATAL ERROR ================")
			fmt.Println(err)
			fmt.Println("================  END FATAL ERROR  ================")
			fmt.Println()
		}
		
		errorOut[errorIndex] = err
		errorIndex++
	}
	
	if len(lcm.shutdownChan) != 0 {
		x := <- lcm.shutdownChan
		return &x, errorOut
	}
	
	return nil, errorOut
}

func (lcm *lifecycleManager) Shutdown(exitCode int) {
	lcm.shutdownChan <- exitCode
}

func (lcm *lifecycleManager) Suicide(err error, exitCode int) {
	lcm.suicideResults <- err
	lcm.shutdownChan <- exitCode
}

func (lcm *lifecycleManager) WatchForShutdown() chan int {
	return lcm.shutdownChan
}

func (lcm *lifecycleManager) CreateRoutine(routine func()) {
	if routine == nil {
		lcm.Suicide(errors.New("empty routine supplied to CreateRoutine"), -2)
		return
	}
	
	lcm.routines.Add(1)
	go func() {
		defer lcm.routines.Done()
		routine()
	}()
}