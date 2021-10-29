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

func SendLog(s *discordgo.Session, logMessage string) {

	channels, _ := s.GuildChannels(config.GuildID)
	ch, err := GetChannelByName(&channels, config.LogChannelName)
	if err != nil {
		s.ChannelMessageSend(ch.ID, logMessage)
	}
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

func FetchMember(s *discordgo.Session, userDetails string) (member *discordgo.Member, err error) {
	guildMembers, err := s.GuildMembers(config.GuildID, "", 1000)

	helpers.RaiseError(err)
	for _, member := range guildMembers {
		if member.User.ID == userDetails {
			return member, nil
		}
	}
	return guildMembers[0], errors.New("member not found")
}

func IsAdmin(s *discordgo.Session, m *discordgo.Member) bool {

	guildRoles, _ := s.GuildRoles(config.GuildID)

	for _, role := range guildRoles {
		fmt.Println(role.Permissions)
	}

	for _, role := range m.Roles {
		fmt.Println(role)
	}

	return false
}
