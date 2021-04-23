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
	"io"
	"log"
	"os"

	"github.com/angch/discordbot/pkg/bothandler"
	"github.com/chzyer/readline"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// testbotCmd represents the testbot command
var testbotCmd = &cobra.Command{
	Use:   "testbot",
	Short: "Test the bot on the command line, without connecting to discord/slack",
	Long:  `Test the bot on the command line, without connecting to discord/slack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Test Bot is now running.  Press CTRL-C to exit.")

		n := bothandler.NewMessagePlatformFromDev()
		bothandler.RegisterMessagePlatform(n)

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		historyfile := home + "/." + rootCmd.Use + "_history"

		l, err := readline.NewEx(&readline.Config{
			Prompt:          "\033[31m»\033[0m ",
			HistoryFile:     historyfile,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",

			HistorySearchFold: true,
		})
		if err != nil {
			panic(err)
		}
		defer l.Close()

		for {
			line, err := l.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			content := line
			h, ok := bothandler.Handlers[content]
			if ok {
				response := h()
				log.Println("Bot says", response)
			}

			// Can be better to decouple 1 to 1 of message : response
			for _, v := range bothandler.CatchallHandlers {
				r := v(content)
				if r != "" {
					log.Println("Bot says", r)
				}
			}

			switch {
			case line == "bye", line == "quit":
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
