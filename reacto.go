package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tobio/reacto/commands"
	"tobio/reacto/commands/handlers"
	"tobio/reacto/config"

	"github.com/bwmarrin/discordgo"
)

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	if err != nil {
		fmt.Println("Error starting up:")
		fmt.Println(err)
	}

	dg.AddHandler(handlers.OnReady)
	dg.AddHandler(handlers.OnNewMember)
	dg.AddHandler(handlers.AdminCommands)
	dg.AddHandler(handlers.OnReaction)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged
	err = dg.Open() // Open the websocket
	if err != nil {
		fmt.Println("Error initialising websocket:")
		fmt.Println(err)
	}

	_, err = dg.ApplicationCommandCreate(config.AppID, config.GuildID, commands.EraseCommand)
	if err != nil {
		fmt.Println("Error adding erase command:")
		fmt.Println(err)
	} else {
		fmt.Println("Erase command added")
	}
	_, err = dg.ApplicationCommandCreate(config.AppID, config.GuildID, commands.ForceLogCommand)
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
