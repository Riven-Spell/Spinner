package cliFE

import (
	"bufio"
	"runtime"
	"strings"
)

var cliReader *bufio.Reader

func getCLIInput() []string {
	buffer, _, err := cliReader.ReadLine()

	if err != nil {
		return nil
	}

	output := string(buffer)
	if runtime.GOOS == "windows" { // Sanity check
		output = strings.TrimSuffix(output, "\r")
	}

	return strings.Split(output, " ")
}
