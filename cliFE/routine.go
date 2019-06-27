package cliFE

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/virepri/Spinner/common"
)

var lcm common.LifecycleManager

func RunCLI() {
	cliReader = bufio.NewReader(os.Stdin)
	cmdChan := make(chan []string)
	common.CliVars.Initialized = true
	lcm = common.GetLifecycleManager() // Gets global lifecycle manager

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
			cmd[0] = strings.ToLower(cmd[0])

			if c, ok := commands[cmd[0]]; ok {
				c(cmd[1:])
			} else {
				lcm.Log(fmt.Sprintf("%s is not a valid command.", cmd[0]), common.ELogLevel.Information())
			}
		}
	}
}
