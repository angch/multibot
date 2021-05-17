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

	"github.com/angch/discordbot/pkg/bothandler"
	"github.com/spf13/cobra"
)

// sendmsgCmd represents the sendmsg command
var sendmsgCmd = &cobra.Command{
	Use:   "sendmsg",
	Short: "Send a message to channel as bot, outside of the event loop",
	Long:  `Send a message to channel as bot, outside of the event loop`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Println("Not enough params")
			return
		}
		channel := args[0]
		mesg := strings.Join(args[1:], " ")

		discordtoken := os.Getenv("DISCORDTOKEN")
		if discordtoken != "" {
			n, err := bothandler.NewMessagePlatformFromDiscord(discordtoken)
			if err != nil {
				log.Fatal(err)
			}
			bothandler.RegisterPassiveMessagePlatform(n)
		}

		slackAppToken := os.Getenv("SLACK_APP_TOKEN")
		slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
		if slackAppToken != "" && slackBotToken != "" {
			if !strings.HasPrefix(slackAppToken, "xapp-") {
				fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
			}
			if !strings.HasPrefix(slackBotToken, "xoxb-") {
				fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
			}

			s, err := bothandler.NewMessagePlatformFromSlack(slackBotToken, slackAppToken)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Slack bot is now running.")
			bothandler.RegisterPassiveMessagePlatform(s)
			go s.ProcessMessages()
		}

		err := bothandler.ChannelMessageSend(channel, mesg)
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
