package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tobio/reacto/config"
	eventhandlers "tobio/reacto/eventHandlers"

	"github.com/bwmarrin/discordgo"
)

func main() {

	dg, err := discordgo.New("Bot " + config.Key)
	if err != nil {
		fmt.Println("Error starting up:")
		fmt.Println(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged
	err = dg.Open() // Open the websocket
	if err != nil {
		fmt.Println("Error initialising websocket:")
		fmt.Println(err)
	}

	//Add all event handlers
	eventhandlers.AddEventHandlers(dg)

	// Initialise slash commands
	eventhandlers.CreateCommands(dg)

	// Create channel, hold it open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}
