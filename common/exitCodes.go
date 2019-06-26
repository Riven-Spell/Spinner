package common

import (
	"fmt"
)

type ExitCode struct {
	Code   int
	Reason string
}

// Enumeration pattern
var EExitCode = ExitCode{0, ""}

func (ExitCode) CleanExit() ExitCode {
	return ExitCode{0, ""}
}

func (ExitCode) FailedVerify() ExitCode {
	return ExitCode{1, "Failed to verify arguments supplied."}
}

func (ExitCode) NilRoutine() ExitCode {
	return ExitCode{2, "A nil routine was supplied to LifecycleManager.CreateRoutine()"}
}

func (ExitCode) FailedDiscordGo() ExitCode {
	return ExitCode{3, "Failed to create a new DiscordGo session."}
}

func (ExitCode) FailedCLIRead() ExitCode {
	return ExitCode{4, "Failed to read a line from the command line."}
}

func (e ExitCode) String() string {
	return fmt.Sprintf("Error code %d: %s", e.Code, e.Reason)
}
