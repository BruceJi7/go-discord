package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"tobio/reacto/config"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = "!"

// const KEY string = "Njk3NzgwMDAwNTg5MzQ4ODY2.Xo8Qug.xUbrZGy8fco3dQJzy8m-qcd42tY"

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	raiseError(err)

	dg.AddHandler(ready)
	dg.AddHandler(message)
	dg.AddHandler(discordMessageReactionAdd)
	// Looks like ready is an expected name, and so this line registers a function to be performed on ready.

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMessageReactions

	err = dg.Open() // Open the websocket
	raiseOrPrint(err, "Bot is running. CTRL-C to exit.")

	// Not sure what this stuff is doing. Concurrency?
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	channels, _ := s.GuildChannels(config.Guild)

	ch, err := getChannelByName(&channels, "general")
	raiseError(err)

	s.ChannelMessageSend(ch.ID, "SUP")

}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// channel := event.ChannelID
	msg := m.Message

	if strings.HasPrefix(msg.Content, PREFIX) {
		fmt.Println("COMMAND")
	}

	fmt.Println(msg.Content)

}

func discordMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	fmt.Println(m.MessageID)
}

// // Get channel by name
func getChannelByName(chans *[]*discordgo.Channel, name string) (c *discordgo.Channel, err error) {
	for _, c := range *chans {
		if c.Name == name {
			return c, nil
		}
	}
	return (*chans)[0], errors.New("channel not found")
}

func raiseError(err error) {
	if err != nil {
		panic(err)
	}
}
func raiseOrPrint(err error, msg string) {
	if err != nil {
		panic(err)
	} else {
		fmt.Println(msg)
	}
}
