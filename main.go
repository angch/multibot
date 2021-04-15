package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Mostly cripped from the discordgo examples

func main() {
	token := os.Getenv("TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

type StoicResponse struct {
	id		string
	body		string
	author_id	int
	author		string
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "World!")
	}

	if m.Content == "o/" {
		s.ChannelMessageSend(m.ChannelID, "\\o")
	}

	if m.Content == "!stoic" {
		url := "https://stoicquotesapi.com/v1/api/quotes/random"

		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("error retrieving stoicquotesapi", err)
			return
		}

		defer resp.Body.Close()

		var respBody StoicResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&respBody)

		if decodeErr != nil {
			fmt.Println("error decoding stoicquotesapi response", decodeErr)
			return
		}

		message := respBody.body + " — " + respBody.author
		s.ChannelMessageSend(m.ChannelID, message)
	}
}
