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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// testbotCmd represents the testbot command
var testbotCmd = &cobra.Command{
	Use:   "testbot",
	Short: "Test the bot on the command line, without connecting to discord/slack",
	Long:  `Test the bot on the command line, without connecting to discord/slack`,
	Run: func(cmd *cobra.Command, args []string) {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		if true {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			historyfile := home + "/." + rootCmd.Use + "_history"

			n, err := bothandler.NewMessagePlatformFromReadline(historyfile, sc)
			if err != nil {
				log.Fatal(err)
			}
			bothandler.RegisterMessagePlatform(n)
			go n.ProcessMessages()
			fmt.Println("Test Bot is now running.  Press CTRL-C to exit.")
		}

		<-sc
		bothandler.Shutdown()
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
