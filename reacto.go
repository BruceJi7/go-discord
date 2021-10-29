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
	help.RaiseError(err)

	dg.AddHandler(onReady)
	dg.AddHandler(onMessage)
	dg.AddHandler(onReaction)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged
	err = dg.Open() // Open the websocket
	help.RaiseOrPrint(err, "Bot is running. CTRL-C to exit.")

	// Not sure what this stuff is doing. Concurrency?
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func onReady(s *discordgo.Session, event *discordgo.Ready) {

	channels, _ := s.GuildChannels(config.GuildID)

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
