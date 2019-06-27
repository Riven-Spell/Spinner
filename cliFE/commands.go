package cliFE

import (
	"fmt"

	"github.com/virepri/Spinner/common"
)

var commands = map[string]func([]string){
	"shutdown":   Shutdown,
	"listguilds": ListGuilds,
	"sethost":    SetHostGuild,
	"help":       Help,
}

func SetHostGuild(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: sethost [guild id]\nSee `listguilds` to find guild IDs.")
		return
	}

	if g, ok := common.BotVars.Guilds[args[0]]; ok {
		common.BotVars.HostGuild = args[0]
		fmt.Printf("%s is now the home guild.\n\n", g.Name)
	} else {
		fmt.Printf("%s is not a valid guild ID.\nSee `listguilds` to find guild IDs.\n\n", args[0])
	}
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
