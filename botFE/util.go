package botFE

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/virepri/Spinner/common"
)

func SetupGuildList(discord *discordgo.Session) {
	common.BotVars.Guilds = make(map[string]*discordgo.UserGuild)

	guilds, err := discord.UserGuilds(-1, "", "")
	if err != nil {
		lcm.Log(fmt.Sprintf("Failed to setup internal guild list: %s", err), common.ELogLevel.Error())
	}

	for _, v := range guilds {
		common.BotVars.Guilds[v.ID] = v
	}
}
