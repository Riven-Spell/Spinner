package cliFE

import (
	"fmt"

	"github.com/virepri/Spinner/common"
)

var commands = map[string]func([]string){
	"shutdown":   Shutdown,
	"listguilds": ListGuilds,
	"help":       Help,
}

func Shutdown(args []string) {
	fmt.Println("Shutting down...")
	lcm.Shutdown(common.ExitCode{Reason: "Operator requested shutdown."})
}

func ListGuilds(args []string) {
	for k, v := range common.BotVars.Guilds {
		fmt.Printf("ID %s: %s\n", k, v.Name)
	}
	fmt.Println()
}

func Help(args []string) {
	var text string
	var ok = false

	if len(args) > 0 {
		text, ok = helpText[args[0]]
	}

	if ok {
		fmt.Printf("%s:\n%s\n", args[0], text)
	} else {
		for k, v := range helpText {
			fmt.Printf("%s:\n%s\n\n", k, v)
		}
	}
}
