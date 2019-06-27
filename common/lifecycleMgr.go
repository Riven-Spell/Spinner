package common

import (
	"errors"
	"fmt"
	"io"
	"os"
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
		shutdownChan:    make(chan ExitCode, 10), // There will probably never be more than 1 in shutdownChan or 10 in suicideResults
		suicideResults:  make(chan error, 500),   // But it's best to be safer than sorry.
		logDestination:  os.Stdout,
		minimumLogLevel: ELogLevel.Information(),
	}
}

type LifecycleManager interface {
	// ===== Routine Wrangling =====

	// SurrenderControl blocks until Shutdown or Suicide is called.
	SurrenderControl() (*ExitCode, []error)

	// Shutdown gracefully shuts down the app.
	Shutdown(exitCode ExitCode)
	// Suicide gracefully shuts down the app and displays an error.
	Suicide(err error, exitCode ExitCode)
	// WatchForShutdown returns the underlying shutdownChan. Useful inside select statements.
	WatchForShutdown() chan ExitCode

	// CreateRoutine adds a waiter to the underlying WaitGroup.
	// This is intended for primary routines only. Child routines should use their own suicide methods.
	CreateRoutine(func())

	// ===== Logging =====

	// Sets log destination to an io.Writer
	SetLogDestination(writer io.Writer)
	// Sets minimum log level
	SetMinimumLogLevel(minimum LogLevel)
	// Write to the log
	Log(Input string, severity LogLevel)
}

type lifecycleManager struct {
	// Routine wranglers
	routines       sync.WaitGroup
	shutdownChan   chan ExitCode
	suicideResults chan error

	// Logging
	logDestination  io.Writer
	minimumLogLevel LogLevel
}

func (lcm *lifecycleManager) SurrenderControl() (*ExitCode, []error) {
	lcm.routines.Wait()

	//noinspection GoPrintFunctions
	fmt.Println("\n") // Everything's closing out, create a clean exit line.

	errorOut := make([]error, len(lcm.suicideResults))
	errorIndex := 0
	for len(lcm.suicideResults) != 0 {
		err := <-lcm.suicideResults

		if err != nil {
			lcm.Log(err.Error(), ELogLevel.Fatal())
		}

		errorOut[errorIndex] = err
		errorIndex++
	}

	if len(lcm.shutdownChan) != 0 {
		x := <-lcm.shutdownChan
		fmt.Printf("Exiting with code %d: %s\n", x.Code, x.Reason)
		return &x, errorOut
	}

	return nil, errorOut
}

func (lcm *lifecycleManager) Shutdown(exitCode ExitCode) {
	lcm.shutdownChan <- exitCode
}

func (lcm *lifecycleManager) Suicide(err error, exitCode ExitCode) {
	lcm.suicideResults <- err
	lcm.shutdownChan <- exitCode
}

func (lcm *lifecycleManager) WatchForShutdown() chan ExitCode {
	return lcm.shutdownChan
}

func (lcm *lifecycleManager) CreateRoutine(routine func()) {
	if routine == nil {
		lcm.Suicide(errors.New("empty routine supplied to CreateRoutine"), EExitCode.NilRoutine())
		return
	}

	lcm.routines.Add(1)
	go func() {
		defer lcm.routines.Done()
		routine()
	}()
}

func (lcm *lifecycleManager) SetLogDestination(writer io.Writer) {
	lcm.logDestination = writer
}

func (lcm *lifecycleManager) SetMinimumLogLevel(minimum LogLevel) {
	lcm.minimumLogLevel = minimum
}

func (lcm *lifecycleManager) Log(Input string, severity LogLevel) {
	if CanLog(lcm.minimumLogLevel, severity) {
		splitDest := false
		if _, splitDest = lcm.logDestination.(*RedirectWriter); (splitDest || lcm.logDestination == os.Stdout) && CliVars.Initialized {
			fmt.Println()
		}
		_, err := lcm.logDestination.Write([]byte(fmt.Sprintf("%s: %s\n", severity.Type, Input)))
		if (splitDest || lcm.logDestination == os.Stdout) && CliVars.Initialized {
			fmt.Print("> ") // Get the console out of the way.
		}

		if err != nil {
			fmt.Println("WARNING: Failed to write following line to log destination...")
			fmt.Printf("%s: %s\n", severity.Type, Input)
		}
	}
}
