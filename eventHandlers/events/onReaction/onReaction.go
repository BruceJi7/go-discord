package onReaction

import (
	"fmt"
	"tobio/reacto/config"
	"tobio/reacto/constant"
	disc "tobio/reacto/discordHelpers"
	"tobio/reacto/eventHandlers/events/onReaction/learningResources"
	"tobio/reacto/eventHandlers/events/onReaction/reactForRole"

	"github.com/bwmarrin/discordgo"
)

func ParseReactionAdded(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	emojiUsed := m.Emoji.MessageFormat()

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
		learningDiscussionChannel, _ := disc.GetChannelByName(s, "learning-discussion")
		learningResourcesChannel, _ := disc.GetChannelByName(s, "learning-resources")

		if m.ChannelID == learningDiscussionChannel.ID && emojiUsed == constant.LearningEmoji {

			learningResources.LearningResourcePost(s, m, learningDiscussionChannel, learningResourcesChannel)

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
