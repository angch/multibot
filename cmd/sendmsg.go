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
	"net/url"
	"os"
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
	"github.com/spf13/cobra"
	"gopkg.in/irc.v3"
)

// sendmsgCmd is a Cobra command that sends a message to a specified channel on messaging platforms
// outside of the main event loop. This command is useful for sending one-off messages without
// starting the full bot event processing.
//
// Usage:
//
//	sendmsg <platform> <channel> <message>
//
// Parameters:
//   - platform: The messaging platform to send the message to. Supported values:
//   - "discord": Send to Discord (requires DISCORDTOKEN environment variable)
//   - "slack": Send to Slack (requires SLACK_APP_TOKEN and SLACK_BOT_TOKEN environment variables)
//   - "telegram": Send to Telegram (requires TELEGRAM_BOT_TOKEN environment variable)
//   - "mattermost": Send to Mattermost (requires MATTERMOST_BOT_TOKEN and MATTERMOST_URL environment variables)
//   - "irc": Send to IRC (requires IRC_CONN environment variable with connection URL)
//   - "all": Send to all configured platforms
//   - channel: The target channel name or ID to send the message to
//   - message: The message content to send (multiple words will be joined with spaces)
//
// Environment Variables:
//   - DISCORDTOKEN: Discord bot token
//   - SLACK_APP_TOKEN: Slack app token (must start with "xapp-")
//   - SLACK_BOT_TOKEN: Slack bot token (must start with "xoxb-")
//   - TELEGRAM_BOT_TOKEN: Telegram bot token
//   - MATTERMOST_BOT_TOKEN: Mattermost bot token
//   - MATTERMOST_URL: Mattermost server URL
//   - IRC_CONN: IRC connection URL (format: irc://username:password@host/channel)
//
// Examples:
//
//	sendmsg discord general "Hello Discord!"
//	sendmsg slack random "Test message"
//	sendmsg all announcements "Message to all platforms"
var sendmsgCmd = &cobra.Command{
	Use:     "sendmsg",
	Short:   "Send a message to channel as bot, outside of the event loop",
	Long:    `Send a message to channel as bot, outside of the event loop, params are platform, channel, and message.`,
	Example: `  sendmsg discord general "Hello Discord!"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Println("Not enough params")
			return
		}
		platform := args[0]
		channel := args[1]
		mesg := strings.Join(args[2:], " ")
		sc := make(chan os.Signal, 1)

		if platform == "discord" || platform == "all" {
			discordtoken := os.Getenv("DISCORDTOKEN")
			if discordtoken != "" {
				n, err := bothandler.NewMessagePlatformFromDiscord(discordtoken)
				if err != nil {
					log.Fatal(err)
				}
				bothandler.RegisterPassiveMessagePlatform(n)
			}
		}

		if platform == "slack" || platform == "all" {
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
				s.DefaultChannel = "random"
				bothandler.RegisterPassiveMessagePlatform(s)
			}
		}

		if platform == "telegram" || platform == "all" {
			telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
			if telegramBotToken != "" {
				s, err := bothandler.NewMessagePlatformFromTelegram(telegramBotToken)
				if err != nil {
					log.Fatal(err)
				}
				s.DefaultChannel = "offtopic"
				log.Println("Telegram bot is now running.")
				bothandler.RegisterMessagePlatform(s)
				go s.ProcessMessages()
			}
		}

		if platform == "mattermost" || platform == "all" {
			mattermostBotToken := os.Getenv("MATTERMOST_BOT_TOKEN")
			mattermostURL := os.Getenv("MATTERMOST_URL")
			if mattermostBotToken != "" && mattermostURL != "" {
				s, err := bothandler.NewMessagePlatformFromMattermost(mattermostBotToken, mattermostURL)
				if err != nil {
					log.Fatal(err)
				}
				mattermost_channel := os.Getenv("MATTERMOST_CHANNEL")
				if mattermost_channel != "" {
					s.DefaultChannel = mattermost_channel
				}
				log.Println("Mattermost bot is now running.")
				bothandler.RegisterMessagePlatform(s)
				go s.ProcessMessages()
			}
		}

		if platform == "irc" || platform == "all" {
			ircConn := os.Getenv("IRC_CONN")
			if ircConn != "" {
				ircParams, err := url.Parse(ircConn)
				if err == nil {
					password, _ := ircParams.User.Password()
					username := ircParams.User.Username()
					config := irc.ClientConfig{
						User: username,
						Nick: username,
						Name: username,
						Pass: password,
					}
					s, err := bothandler.NewMessagePlatformFromIrc(ircParams.Host, &config, sc)
					if err != nil {
						log.Fatal(err)
					}
					s.DefaultChannel = strings.TrimPrefix(ircParams.Path, "/")

					log.Println("Irc bot is now running.")
					bothandler.RegisterMessagePlatform(s)
					go s.ProcessMessages()
				}
			}
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
