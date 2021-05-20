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

	"github.com/angch/discordbot/pkg/engineersmy"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// discordlsCmd represents the discordls command
var discordlsCmd = &cobra.Command{
	Use:   "discordls",
	Short: "List members in a discord Guild",
	Long:  `List all members in a discord guild`,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("TOKEN")
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}
		discordChan, err := dg.Channel(engineersmy.KnownChannels["general"])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v\n", discordChan)
		guildId := discordChan.GuildID
		members, err := dg.GuildMembers(guildId, "", 1000)
		if err != nil {
			log.Fatal(err)
		}
		output := "`select * from users;`\n"
		output += "`|----------------------|----------------------|`\n"
		output += fmt.Sprintf("`| %-20s | %-20s |`\n", "username", "nick")
		output += "`|----------------------|----------------------|`\n"
		for _, v := range members {
			log.Printf("%+v\n", *v)
			log.Println(v.User.Username, v.Nick)
			output += fmt.Sprintf("`| %-20s | %-20s |`\n", v.User.Username, v.Nick)
		}
		output += "`|----------------------|----------------------|`\n"

		_, err = dg.ChannelMessageSend(engineersmy.KnownChannels["general"], output)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(discordlsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discordlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discordlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
