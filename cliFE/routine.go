package cliFE

import (
	"bufio"
	"fmt"
	"os"

	"github.com/virepri/Spinner/common"
)

func RunCLI() {
	cliReader = bufio.NewReader(os.Stdin)
	cmdChan := make(chan []string)
	common.CliVars.Initialized = true
	lcm := common.GetLifecycleManager() // Gets global lifecycle manager

	for {
		go func() {
			fmt.Print("> ")
			cmdChan <- getCLIInput()
		}()

		select {
		case exitCode := <-lcm.WatchForShutdown():
			lcm.Shutdown(exitCode)
			return
		case cmd := <-cmdChan:
			fmt.Println("command ", cmd) // TODO: parse and handle CLI commands
		}
	}
}
