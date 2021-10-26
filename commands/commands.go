package commands

import (
	"fmt"
	"strconv"

	"tobio/reacto/discordHelpers"

	"github.com/bwmarrin/discordgo"
)

func DeleteMessages(c discordHelpers.BotCommand, s *discordgo.Session, m *discordgo.MessageCreate) {

	fmt.Println(c.Command)
	howMany, _ := strconv.Atoi(c.Instruction[0])

	fmt.Println(howMany)

	messages, _ := s.ChannelMessages(m.ChannelID, howMany, m.Message.ID, "", "")

	var messageIDs []string

	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	messageIDs = append(messageIDs, m.Message.ID)

	err := s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
	// s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	if err != nil {
		fmt.Println(err)
	}
}
