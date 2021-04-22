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
	"strings"

	"github.com/angch/discordbot/pkg/engineersmy"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// sendmsgCmd represents the sendmsg command
var sendmsgCmd = &cobra.Command{
	Use:   "sendmsg",
	Short: "Send a message to channel as bot, outside of the event loop",
	Long:  `Send a message to channel as bot, outside of the event loop`,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("TOKEN")
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}

		if len(args) < 2 {
			log.Println("Not enough params")
			return
		}
		channel := args[0]
		mesg := strings.Join(args[1:], " ")

		channelId, ok := engineersmy.KnownChannels[channel]
		if !ok {
			log.Println("Unknown channel", channel)
			return
		}

		_, err = dg.ChannelMessageSend(channelId, mesg)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendmsgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendmsgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendmsgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
