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
	lcm := common.GetLifecycleManager() // Gets global lifecycle manager

	for {
		go func() {
			cmdChan <- getCLIInput()
		}()

		select {
		case exitCode := <-lcm.WatchForShutdown():
			lcm.Shutdown(exitCode)
		case cmd := <-cmdChan:
			fmt.Println(cmd) // TODO: parse and handle CLI commands
		}
	}
}
