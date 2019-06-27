package botFE

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/virepri/Spinner/common"
)

func ListenMessages(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Message.Content

	if !strings.HasPrefix(content, "!") { // TODO: Make the prefix configurable from command line.
		return
	}

	cmd := strings.Split(strings.TrimPrefix(content, "!"), " ")
	cmd[0] = strings.ToLower(cmd[0])

	if c, ok := commands[cmd[0]]; ok {
		c(discord, message, cmd[1:])
	} else {
		_, err := discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s, I'm afraid that %s is not a valid or existing command!", message.Author.Username, cmd[0]))

		if err != nil {
			lcm.Log(fmt.Sprintf("Failed to respond to invalid command, %s", err.Error()), common.ELogLevel.Error())
		}
	}
}
