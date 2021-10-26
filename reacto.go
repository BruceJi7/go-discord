package main

import (
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
const MSG_TO_WATCH = "902547426215358575"

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	help.RaiseError(err)

	dg.AddHandler(onReady)
	dg.AddHandler(onMessage)
	dg.AddHandler(onReaction)
	// Looks like ready is an expected name, and so this line registers a function to be performed on ready.

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsAllWithoutPrivileged
	err = dg.Open() // Open the websocket
	help.RaiseOrPrint(err, "Bot is running. CTRL-C to exit.")

	// Not sure what this stuff is doing. Concurrency?
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func onReady(s *discordgo.Session, event *discordgo.Ready) {

	channels, _ := s.GuildChannels(config.Guild)

	ch, err := disc.GetChannelByName(&channels, "general")
	help.RaiseError(err)

	s.ChannelMessageSend(ch.ID, "SUP")

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

	}

	// fmt.Println(msg.Content)

}

func onReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	if m.MessageID == MSG_TO_WATCH {
		print("Pretend role given")
	}
}
