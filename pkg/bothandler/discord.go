package bothandler

import (
	"fmt"
	"log"
	"strings"

	"github.com/angch/discordbot/pkg/engineersmy"
	"github.com/bwmarrin/discordgo"
)

// Implements MessagePlatform
type DiscordMessagePlatform struct {
	Session  *discordgo.Session
	Channels map[string]string
}

func NewMessagePlatformFromDiscord(discordtoken string) (*DiscordMessagePlatform, error) {
	dg, err := discordgo.New("Bot " + discordtoken)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return nil, err
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return nil, err
	}

	return &DiscordMessagePlatform{
		Session:  dg,
		Channels: engineersmy.KnownChannels,
	}, nil
}

// Send to default channel
func (dg *DiscordMessagePlatform) Send(m string) {
	_, err := dg.Session.ChannelMessageSend(dg.Channels[""], m)
	if err != nil {
		log.Println(err)
	}
}

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// FIXME: This can be better
	// Part of first stage refac
	h, ok := Handlers[m.Content]
	if ok {
		response := h()
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   response,
			Reference: m.Reference(),
		})
		if err != nil {
			log.Println(err)
		}
	}

	// Can be better to decouple 1 to 1 of message : response
	for _, v := range CatchallHandlers {
		username := ""
		if m.Author != nil {
			username = m.Author.Username
		}
		r := v(Request{m.Content, "discord", m.ChannelID, username})
		if r != "" {
			_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:   r,
				Reference: m.Reference(),
			})
			if err != nil {
				log.Println(err)
			}
		}
	}

	sliced_content := strings.SplitN(m.Content, " ", 2)
	if len(sliced_content) > 1 {
		command := sliced_content[0]
		actual_content := sliced_content[1]

		username := ""
		if m.Author != nil {
			username = m.Author.Username
		}

		ih, ok := MsgInputHandlers[command]
		if ok {
			response := ih(Request{actual_content, "discord", m.ChannelID, username})
			if response != "" {
				_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content:   response,
					Reference: m.Reference(),
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (dg *DiscordMessagePlatform) ProcessMessages() {
	fmt.Println("Discord Bot is now running.  Press CTRL-C to exit.")

	dg.Session.AddHandler(discordMessageCreate)
}

func (dg *DiscordMessagePlatform) Close() {
	dg.Session.Close()
}

func (s *DiscordMessagePlatform) ChannelMessageSend(channel, message string) error {
	channelId, ok := engineersmy.KnownChannels[channel]
	if !ok {
		log.Println("Unknown channel", channel)
		return fmt.Errorf("Unknown channel %s", channel)
	}
	_, err := s.Session.ChannelMessageSend(channelId, message)

	return err
}
