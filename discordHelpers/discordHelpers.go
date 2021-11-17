package discordHelpers

import (
	"errors"
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
	ch, err := GetChannelByName(&channels, "bot-logs")

	if err != nil {
		panic(err)
	} else {
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

func GetChannelByIDFromSession(s *discordgo.Session, cID string) (c *discordgo.Channel, err error) {

	channels, _ := s.GuildChannels(config.GuildID)
	for _, c := range channels {
		if c.ID == cID {
			return c, nil
		}
	}
	return channels[0], errors.New("channel not found")
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

func IsAdmin(s *discordgo.Session, guildID string, userID string) (bool, error) {
	return memberHasPermission(s, guildID, userID, discordgo.PermissionAdministrator)
}

func memberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}
