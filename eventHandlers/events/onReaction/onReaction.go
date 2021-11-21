package onReaction

import (
	"fmt"
	"tobio/reacto/config"
	disc "tobio/reacto/discordHelpers"
	"tobio/reacto/eventHandlers/events/onReaction/reactForRole"

	"github.com/bwmarrin/discordgo"
)

func ParseReactionAdded(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	emojiUsed := m.Emoji.MessageFormat()

	fmt.Println(emojiUsed)

	member, err := disc.FetchMember(s, m.UserID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added:")
		fmt.Println("Error finding user")
		fmt.Println(err)
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == config.RFRPostID {
		reactForRole.RFRAdd(s, member, emojiUsed)
	} else {
		//If not, might be learning-related
		learningChannel, _ := disc.GetChannelByName(s, "learning-discussion")
		resourcesChannel, _ := disc.GetChannelByName(s, "learning-resources")

		if m.ChannelID == learningChannel.ID && emojiUsed == "💡" {

			message, err := s.ChannelMessage(learningChannel.ID, m.MessageID)
			if err != nil {
				fmt.Println("Whilst parsing reaction added")
				fmt.Println("Whilst handling learning-discussion reaction")
				fmt.Println("Error finding message")
				fmt.Println(err)
				return
			}
			fmt.Println(resourcesChannel.ID)
			fmt.Println(message.ChannelID)
			// TODO - count the amount of reactions
			// If over 5, copy the message over to learning-resources

		}

	}

}

func ParseReactionRemoved(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	emojiUsed := m.Emoji.MessageFormat()

	member, err := disc.FetchMember(s, m.UserID)
	if err != nil {
		fmt.Println("Whilst parsing reaction removed:")
		fmt.Println("Error finding user")
		fmt.Println(err)
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == config.RFRPostID {
		reactForRole.RFRRemove(s, member, emojiUsed)
	}
}
