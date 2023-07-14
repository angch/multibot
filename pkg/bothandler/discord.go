package bothandler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/angch/multibot/pkg/engineersmy"
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
		Channels: engineersmy.KnownDiscordChannels,
	}, nil
}

// Send to default channel
func (dg *DiscordMessagePlatform) Send(text string) {
	if dg != nil {
		dg.SendWithOptions(text, SendOptions{})
	}
}

func (dg *DiscordMessagePlatform) SendWithOptions(text string, options SendOptions) {
	if dg == nil {
		return
	}
	if options.Silent {
		// FIXME: Figuure out how to use ChannelMessageSendComplex to send silent messages
		_, err := dg.Session.ChannelMessageSend(dg.Channels[""], text)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := dg.Session.ChannelMessageSend(dg.Channels[""], text)
		if err != nil {
			log.Println(err)
		}
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

	// Can be better to decouple 1 to 1 of message : response
	for _, v := range CatchallExtendedHandlers {
		r := v(ExtendedMessage{Text: m.Content})
		if r != nil {
			if r.Text != "" {
				_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content:   r.Text,
					Reference: m.Reference(),
				})
				if err != nil {
					log.Println(err)
				}
			}
			if r.Image != nil {
				fileImage := discordgo.File{
					Name: sanitizeFilename(m.Content, "png"),
					// ContentType: "image/jpeg",
					Reader: bytes.NewReader(r.Image),
				}
				msg := &discordgo.MessageSend{
					Content:   m.Content,
					Reference: m.Reference(),
					Files:     []*discordgo.File{&fileImage},
				}
				_, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
				if err != nil {
					log.Println(err)
				}
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

	if m.Attachments != nil {
		targetPixels := 512 * 512
		best := &discordgo.MessageAttachment{}
		bestSize := 100000000

		for _, v := range m.Attachments {
			if v == nil {
				continue
			}
			pixels := v.Height * v.Width
			diff := targetPixels - pixels
			if diff < 0 {
				diff = -diff
			}
			if diff < bestSize {
				bestSize = diff
				best = v
			}
			log.Printf("photosize %+v\n", best)
			// FIXME:
			filename := "tmp/" + best.ID
			err := botDownload(best, filename)
			if err != nil {
				log.Println(err)
			}
			for _, v := range ImageHandlers {
				username := ""
				if m.Author != nil {
					username = m.Author.Username
				}

				req := Request{m.Content, "discord", m.ChannelID, username}
				r := v(filename, req)
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

			if false {
				os.Remove(filename)
			}
		}
	}
}

func (dg *DiscordMessagePlatform) ProcessMessages() {
	// fmt.Println("Discord Bot is now running.  Press CTRL-C to exit.")

	dg.Session.AddHandler(discordMessageCreate)
}

func (dg *DiscordMessagePlatform) Close() {
	dg.Session.Close()
}

func (s *DiscordMessagePlatform) ChannelMessageSend(channel, message string) error {
	channelId, ok := engineersmy.KnownDiscordChannels[channel]
	if !ok {
		log.Println("Unknown channel", channel)
		return fmt.Errorf("Unknown channel %s", channel)
	}
	_, err := s.Session.ChannelMessageSend(channelId, message)

	return err
}

func botDownload(attachment *discordgo.MessageAttachment, localFilename string) error {
	if attachment == nil {
		return fmt.Errorf("No attachment")
	}

	downloadUrl := attachment.URL
	log.Println("Downloading", downloadUrl)

	get, err := http.Get(downloadUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	reader := get.Body
	defer reader.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	out, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	return nil
}
