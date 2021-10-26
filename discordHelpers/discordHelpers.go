package discordHelpers

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type BotCommand struct {
	Command     string
	Instruction []string
}

func ParseCommand(commandString string) BotCommand {
	instructions := strings.Split((commandString), " ")
	return BotCommand{Command: instructions[0], Instruction: instructions[1:]}
}

// // Get channel by name
func GetChannelByName(chans *[]*discordgo.Channel, name string) (c *discordgo.Channel, err error) {
	for _, c := range *chans {
		if c.Name == name {
			return c, nil
		}
	}
	return (*chans)[0], errors.New("channel not found")
}
