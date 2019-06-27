package cliFE

import (
	"bufio"
	"runtime"
	"strings"

	"github.com/virepri/Spinner/common"
)

var cliReader *bufio.Reader

func getCLIInput() []string {
	buffer, err := cliReader.ReadString(byte('\n'))
	buffer = buffer[:len(buffer)-1] // Cut \n

	if err != nil {
		common.GetLifecycleManager().Suicide(err, common.EExitCode.FailedDiscordGo())
		return nil
	}

	output := string(buffer)
	if runtime.GOOS == "windows" { // Sanity check
		output = strings.TrimSuffix(output, "\r")
	}

	return strings.Split(output, " ")
}
