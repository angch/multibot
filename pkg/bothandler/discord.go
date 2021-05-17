package bothandler

import (
	"log"

	"github.com/angch/discordbot/pkg/engineersmy"
	"github.com/bwmarrin/discordgo"
)

// Implements MessagePlatform
type DiscordMessagePlatform struct {
	Session  *discordgo.Session
	Channels map[string]string
}

func NewMessagePlatformFromDiscord(dg *discordgo.Session) *DiscordMessagePlatform {
	return &DiscordMessagePlatform{
		Session:  dg,
		Channels: engineersmy.KnownChannels,
	}
}

// Send to default channel
func (dg *DiscordMessagePlatform) Send(m string) {
	_, err := dg.Session.ChannelMessageSend(dg.Channels[""], m)
	if err != nil {
		log.Println(err)
	}
}
