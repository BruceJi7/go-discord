package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	comm "tobio/reacto/commands"
	"tobio/reacto/config"
	"tobio/reacto/constant"
	disc "tobio/reacto/discordHelpers"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = "!"

var MSG_TO_WATCH string = ""

var Log = disc.NewLogPrefixes()

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	if err != nil {
		fmt.Println("Error starting up:")
		fmt.Println(err)
	}

	dg.AddHandler(onReady)
	dg.AddHandler(onNewMember)
	dg.AddHandler(adminCommands)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged
	err = dg.Open() // Open the websocket
	if err != nil {
		fmt.Println("Error initialising websocket:")
		fmt.Println(err)
	}

	eraseCommand := &discordgo.ApplicationCommand{
		Name:        "erase",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Erase messages in a channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "multiple",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Specify amount to erase",
			},
		},
	}

	forceLogCommand := &discordgo.ApplicationCommand{
		Name:        "forcelog",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Force Bot to Log Something",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "message",
				Type:        discordgo.ApplicationCommandOptionString,
				Description: "Specify log message",
				Required:    true,
			},
		},
	}

	_, err = dg.ApplicationCommandCreate(config.AppID, config.GuildID, eraseCommand)
	if err != nil {
		fmt.Println("Error adding erase command:")
		fmt.Println(err)
	} else {
		fmt.Println("Erase command added")
	}
	_, err = dg.ApplicationCommandCreate(config.AppID, config.GuildID, forceLogCommand)
	if err != nil {
		fmt.Println("Error adding forcelog command:")
		fmt.Println(err)
	} else {
		fmt.Println("Forcelog command added")
	}

	// Create channel, hold it open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	disc.SendLog(s, "init")
}

func onNewMember(s *discordgo.Session, memberJoinEvent *discordgo.GuildMemberAdd) {

	var newUserName string
	if memberJoinEvent.Member.Nick != "" {
		newUserName = memberJoinEvent.Member.Nick
	} else {
		newUserName = memberJoinEvent.User.Username
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	greeting := constant.RandomGreeting(r)
	suggestion := constant.RandomSuggestion(r)
	secondSuggestion := constant.RandomSuggestion(r)
	closing := constant.RandomClosing(r)

	botWelcomeScript := fmt.Sprintf("%s, %s! %s introduce yourself, tell us your coding story.\n %s check out the react-for-roles channel and let us know where you're based!\n %s", greeting, newUserName, suggestion, secondSuggestion, closing)

	channels, _ := s.GuildChannels(config.GuildID)
	welcomeChannel, err := disc.GetChannelByName(&channels, "off-topic")
	if err != nil {
		fmt.Println("Error finding off-topic channel")
		fmt.Println(err)
	} else {
		s.ChannelMessageSend(welcomeChannel.ID, botWelcomeScript)
	}
}

func onReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	emojiUsed := m.Emoji.MessageFormat()

	if emojiUsed == "🔥" {
		if MSG_TO_WATCH == "" {
			fmt.Println("MSG TO WATCH was set to ", m.MessageID)
			MSG_TO_WATCH = m.MessageID
		}
	}

	fmt.Println(m.Emoji.APIName())

	if m.MessageID == MSG_TO_WATCH {
		fmt.Println("Pretend role given")
	}
}

func adminCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options

	interactionID := i.Interaction.ID
	interactionChannel, _ := disc.GetChannelByIDFromSession(s, i.ChannelID)
	interactionMember := i.Member

	interactionMemberIsAdmin, err := disc.IsAdmin(s, config.GuildID, interactionMember.User.ID)
	if err != nil {
		fmt.Println("Error on evaluating admin permissions:")
		fmt.Println(err)
	} else {
		if !interactionMemberIsAdmin {
			return
		}
	}

	switch data.Name {
	case "erase":

		if len(options) == 0 {
			// Triggered single erase mode

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
				deleteErr := comm.DeleteMessages(1, s, interactionChannel.ID, interactionID)
				if deleteErr != nil {
					logmessage := fmt.Sprintf(Log.Error+Log.EraseOne+"User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
					disc.SendLog(s, logmessage)
					fmt.Println("Error deleting one message")
					fmt.Println(deleteErr)
				} else {
					logmessage := fmt.Sprintf(Log.EraseOne+"User %s | channel %s", interactionMember.User.Username, interactionChannel.Name)
					disc.SendLog(s, logmessage)
				}
			}

		} else {
			// Multiple erase mode:
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
				deleteErr := comm.DeleteMessages(int(eraseAmount), s, interactionChannel.ID, interactionID)
				if deleteErr != nil {
					logmessage := fmt.Sprintf(Log.Error+Log.EraseMulti+"User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
					disc.SendLog(s, logmessage)
					fmt.Println("Error deleting messages")
					fmt.Println(deleteErr)
				} else {
					logmessage := fmt.Sprintf(Log.EraseMulti+"User %s | %d messages | channel %s", interactionMember.User.Username, eraseAmount, interactionChannel.Name)
					disc.SendLog(s, logmessage)
				}

			}

		}

	case "forcelog":
		fmt.Println("Force Log")

		err := s.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Log made in log channel", Flags: 1 << 6},
			})

		if err != nil {
			fmt.Println("Error responding to command Forcelog")
			fmt.Println(err)
		} else {
			logString := options[0].StringValue()
			logmessage := fmt.Sprintf(Log.Forcelog+"By User %s: %s", interactionMember.User.Username, logString)
			disc.SendLog(s, logmessage)
			fmt.Println("Force log: ", logString)
		}
	}
}
