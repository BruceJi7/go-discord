package actions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func DeleteMessages(howMany int, s *discordgo.Session, channel string, messageID string) error {

	messages, err := s.ChannelMessages(channel, howMany, messageID, "", "")
	if err != nil {
		fmt.Println("Error getting messages to delete")
		return err
	}
	// fmt.Println(messages)
	var messageIDs []string

	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	messageIDs = append(messageIDs, messageID)

	err = s.ChannelMessagesBulkDelete(channel, messageIDs)
	if err != nil {
		return err
	}

	return nil
}
