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
	"os/signal"
	"strings"
	"syscall"

	"github.com/angch/multibot/pkg/bothandler"
	"github.com/spf13/cobra"
	"gopkg.in/irc.v3"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the discordbot",
	Long:  `Run the discordbot`,
	Run: func(cmd *cobra.Command, args []string) {
		sc := make(chan os.Signal, 1)

		discordtoken := os.Getenv("DISCORDTOKEN")
		if discordtoken != "" {
			n, err := bothandler.NewMessagePlatformFromDiscord(discordtoken)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Discord Bot is now running.")
			bothandler.RegisterMessagePlatform(n)
			go n.ProcessMessages()
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
			s.DefaultChannel = "random"
			log.Println("Slack bot is now running.")
			bothandler.RegisterMessagePlatform(s)
			go s.ProcessMessages()
		}

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

		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		bothandler.Shutdown()
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
