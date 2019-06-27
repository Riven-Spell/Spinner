package botFE

import (
	"fmt"
)

var helpText = map[string]string{
	"help":   `Displays this text.`,
	"invite": `Get an invite for this bot to your server.`,
	"ping":   `Pong!`,
}

func craftHelpMessages() []string {
	out := make([]string, 0)
	currentMessage := "```"

	for k, v := range helpText {
		toAdd := fmt.Sprintf("%s:\n%s\n", k, v)

		if len(currentMessage)+len(toAdd) > 2000-3 {
			out = append(out, currentMessage+"```")
			currentMessage = "```"
		}

		currentMessage += toAdd
	}

	out = append(out, currentMessage+"```")

	return out
}
