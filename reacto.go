package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	comm "tobio/reacto/commands"
	"tobio/reacto/config"
	disc "tobio/reacto/discordHelpers"
	help "tobio/reacto/helpers"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = "!"

var MSG_TO_WATCH string = ""

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	if err != nil {
		fmt.Println("Error starting up:")
		fmt.Println(err)
	}

	dg.AddHandler(onReady)
	dg.AddHandler(commands)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged
	err = dg.Open() // Open the websocket
	if err != nil {
		fmt.Println("Error initialising websocket:")
		fmt.Println(err)
	}

	command := &discordgo.ApplicationCommand{
		Name:        "command-name",
		Type:        discordgo.ChatApplicationCommand,
		Description: "This is command description",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "subcommand",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This is subcommand description",
			},
			{
				Name:        "subcommand-group",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Description: "This is subcommand group description",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "subcommand",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "This is subcommand description",
					},
				},
			},
		},
	}
	_, err = dg.ApplicationCommandCreate(config.AppID, config.GuildID, command)
	if err != nil {
		fmt.Println("Error creating command:")
		fmt.Println(err)
	} else {
		fmt.Println("Command added")
	}

	// Create channel, hold it open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func onReady(s *discordgo.Session, _ *discordgo.Ready) {

	channels, _ := s.GuildChannels(config.GuildID)

	ch, err := disc.GetChannelByName(&channels, "general")
	help.RaiseError(err)
	s.ChannelMessageSend(ch.ID, "SUP")

	disc.SendLog(s, "init")

}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// channel := m.ChannelID
	msg := m.Message
	if strings.HasPrefix(msg.Content, PREFIX) {
		c := disc.ParseCommand(msg.Content)

		if strings.Contains(c.Command, "!clear") {
			comm.DeleteMessages(c, s, m)
		}
		if strings.Contains(c.Command, "!members") {
			member, err := disc.FetchMember(s, m.Author.ID)
			help.RaiseError(err)
			disc.IsAdmin(s, member)
		}

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

func commands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	switch data.Options[0].Name {
	case "command":
		// Do something
	case "subcommand-group":
		data := data.Options[0]
		switch data.Options[0].Name {
		case "subcommand":
			// Do something
			fmt.Println("subcommand 1")
			err := s.InteractionRespond(i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{TTS: true, Content: "A reply"},
				})

			fmt.Println("subcommand error")
			fmt.Println(err)
		}
	case "subcommand":
		data := data.Options[0]
		switch data.Options[0].Name {
		case "subcommand":
			// Do something
			fmt.Println("subcommand 2")
		}
	}
}
