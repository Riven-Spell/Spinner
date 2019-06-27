package botFE

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/virepri/Spinner/common"
)

var lcm common.LifecycleManager

type BotParameters struct {
	Token     string
	HostGuild string // TODO: Config?
}

func BakeRunBot(params BotParameters) func() {
	return func() {
		runBot(params)
	}
}

func runBot(params BotParameters) {
	discord, err := discordgo.New("Bot " + params.Token)
	defer discord.Close()
	lcm = common.GetLifecycleManager()

	if err != nil {
		lcm.Suicide(err, common.EExitCode.FailedDiscordGo())
		return
	}

	discord.AddHandler(ListenMessages)

	err = discord.Open() // TODO (maybe): Sharding? Probably not necessary for this _yet_. Stretch goal if I have time.
	if err != nil {
		lcm.Suicide(fmt.Errorf("failed to open a connection to Discord: %s", err), common.EExitCode.FailedDiscordGo())
		return
	}

	// Perform general post-connect setup
	SetupGuildList(discord)

	// TODO: Setup guild functions.
	if _, ok := common.BotVars.Guilds[params.HostGuild]; params.HostGuild != "" && !ok {
		lcm.Log(fmt.Sprintf("Could not locate host guild %s.", params.HostGuild), common.ELogLevel.Warning())
	} else {
		// TODO: Create listmissing command
		lcm.Log("No host guild has been selected. Some features will not be available. Run the command \"listmissing\" to find out what you're missing out on.", common.ELogLevel.Warning())
	}

	for {
		select {
		// TODO: Receive commands from channel
		case exitCode := <-lcm.WatchForShutdown():
			lcm.Shutdown(exitCode)
			return
		}
	}
}
