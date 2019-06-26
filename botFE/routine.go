package botFE

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/virepri/Spinner/common"
)

type BotParameters struct {
	Token string
}

func BakeRunBot(params BotParameters) func() {
	return func() {
		runBot(params)
	}
}

func runBot(params BotParameters) {
	discord, err := discordgo.New("Bot " + params.Token)
	defer discord.Close()
	lcm := common.GetLifecycleManager()

	if err != nil {
		lcm.Suicide(err, common.EExitCode.FailedDiscordGo())
		return
	}

	err = discord.Open()
	if err != nil {
		lcm.Suicide(fmt.Errorf("failed to open a connection to Discord: %s", err), common.EExitCode.FailedDiscordGo())
		return
	}

	for {
		select {
		case exitCode := <-lcm.WatchForShutdown():
			lcm.Shutdown(exitCode)
			return
		}
	}
}
