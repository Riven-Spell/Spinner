package botFE

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/virepri/Spinner/common"
)

var commands = map[string]func(*discordgo.Session, *discordgo.MessageCreate, []string){
	"ping":   Ping, // Sample command to make sure this structure works.
	"invite": Invite,
}

func Ping(discord *discordgo.Session, message *discordgo.MessageCreate, params []string) {
	_, err := discord.ChannelMessageSend(message.ChannelID, "Pong!")

	if err != nil {
		lcm.Log(fmt.Sprintf("Failed to respond to ping command, %s", err.Error()), common.ELogLevel.Error())
	}
}

func Invite(discord *discordgo.Session, message *discordgo.MessageCreate, params []string) {
	_, err := discord.ChannelMessageSend(message.ChannelID, "Invite this Spinner instance to your servers: \nhttps://discordapp.com/api/oauth2/authorize?client_id=592955706861420545&permissions=8&scope=bot\n\nNote: It will not run civilizations on your chat.")

	if err != nil {
		lcm.Log(fmt.Sprintf("Failed to respond to invite command, %s", err.Error()), common.ELogLevel.Error())
	}
}
