package discordHelpers

import (
	"errors"
	"fmt"
	"strings"
	"tobio/reacto/config"
	"tobio/reacto/helpers"

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

func FetchMember(s *discordgo.Session, userDetails string) {
	guildMembers, err := s.GuildMembers(config.GuildID, "", 1000)
	helpers.RaiseError(err)
	fmt.Println(guildMembers)
}
