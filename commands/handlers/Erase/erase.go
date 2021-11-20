package erase

import (
	"fmt"

	act "tobio/reacto/commands/actions"
	disc "tobio/reacto/discordHelpers"

	"github.com/bwmarrin/discordgo"
)

func SingleErase(s *discordgo.Session, i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member) {

	err := s.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Messages Erased", Flags: 1 << 6},
		})

	if err != nil {
		fmt.Println("Error responding to command Erase")
		fmt.Println(err)
	} else {
		fmt.Println("Trigger Erase Command")
		deleteErr := act.DeleteMessages(1, s, interactionChannel.ID, interactionID)
		if deleteErr != nil {
			logmessage := fmt.Sprintf(disc.Log.Error+disc.Log.EraseOne+"User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
			disc.SendLog(s, logmessage)
			fmt.Println("Error deleting one message")
			fmt.Println(deleteErr)
		} else {
			logmessage := fmt.Sprintf(disc.Log.EraseOne+"User %s | channel %s", interactionMember.User.Username, interactionChannel.Name)
			disc.SendLog(s, logmessage)
		}
	}

}

func MultiErase(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member) {

	eraseAmount := options[0].IntValue()
	fmt.Println(eraseAmount)
	err := s.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Messages Erased", Flags: 1 << 6},
		})
	if err != nil {
		fmt.Println("Error responding to command Erase")
		fmt.Println(err)
	} else {
		fmt.Println("Trigger Multiple Erase Command: ", eraseAmount)
		deleteErr := act.DeleteMessages(int(eraseAmount), s, interactionChannel.ID, interactionID)
		if deleteErr != nil {
			logmessage := fmt.Sprintf(disc.Log.Error+disc.Log.EraseMulti+"User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
			disc.SendLog(s, logmessage)
			fmt.Println("Error deleting messages")
			fmt.Println(deleteErr)
		} else {
			logmessage := fmt.Sprintf(disc.Log.EraseMulti+"User %s | %d messages | channel %s", interactionMember.User.Username, eraseAmount, interactionChannel.Name)
			disc.SendLog(s, logmessage)
		}

	}
}
