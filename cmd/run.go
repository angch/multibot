/*
Copyright © 2021 Ang Chin Han <ang.chin.han@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/angch/discordbot/pkg/bothandler"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// FIXME: This can be better
	// Part of first stage refac
	h, ok := bothandler.Handlers[m.Content]
	if ok {
		response := h()
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Println(err)
		}
	}

	// Can be better to decouple 1 to 1 of message : response
	for _, v := range bothandler.CatchallHandlers {
		r := v(m.Content)
		if r != "" {
			_, err := s.ChannelMessageSend(m.ChannelID, r)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the discordbot",
	Long:  `Run the discordbot`,
	Run: func(cmd *cobra.Command, args []string) {
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
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		dg.Close()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
