package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"tobio/reacto/config"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = "!"

type BotCommand struct {
	command     string
	instruction []string
}

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	raiseError(err)

	dg.AddHandler(ready)
	dg.AddHandler(message)
	dg.AddHandler(reaction)
	// Looks like ready is an expected name, and so this line registers a function to be performed on ready.

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsAllWithoutPrivileged
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

	channels, _ := s.GuildChannels(config.Guild)

	ch, err := getChannelByName(&channels, "general")
	raiseError(err)

	s.ChannelMessageSend(ch.ID, "SUP")

}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// channel := m.ChannelID
	msg := m.Message

	if strings.HasPrefix(msg.Content, PREFIX) {
		c := parseCommand(msg.Content)

		if strings.Contains(c.command, "!clear") {
			// This is redundant because it only clears the message that you use to tell it to clear something
			// But it is a proof of concept lol

			fmt.Println(c.command)
			howMany, _ := strconv.Atoi(c.instruction[0])

			fmt.Println(howMany)

			messages, _ := s.ChannelMessages(m.ChannelID, howMany, msg.ID, "", "")

			var messageIDs []string

			for _, m := range messages {
				messageIDs = append(messageIDs, m.ID)
			}

			err := s.ChannelMessagesBulkDelete(msg.ChannelID, messageIDs)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	// fmt.Println(msg.Content)

}

func reaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	fmt.Println(m.MessageID)
}

func parseCommand(commandString string) BotCommand {
	instructions := strings.Split((commandString), " ")
	return BotCommand{command: instructions[0], instruction: instructions[1:]}
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
