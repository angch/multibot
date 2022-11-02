package bothandler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/angch/discordbot/pkg/engineersmy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Implements MessagePlatform
type TelegramMessagePlatform struct {
	Client         *tgbotapi.BotAPI
	ChannelId      map[string]string
	KnownUsers     map[string]tgbotapi.User
	KnownUsersLock sync.RWMutex
	DefaultChannel string
}

func NewMessagePlatformFromTelegram(telegrambottoken string) (*TelegramMessagePlatform, error) {
	bot, err := tgbotapi.NewBotAPI(telegrambottoken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramMessagePlatform{
		Client:     bot,
		ChannelId:  map[string]string{},
		KnownUsers: map[string]tgbotapi.User{},
	}, nil
}

func (s *TelegramMessagePlatform) ProcessMessages() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := s.Client.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// log.Printf("[%s] %s %d %d", update.Message.From.UserName, update.Message.Text, update.Message.From.ID, update.Message.Chat.ID)
		s.KnownUsersLock.Lock()
		s.KnownUsers[update.Message.From.UserName] = *update.Message.From
		s.KnownUsersLock.Unlock()

		content := update.Message.Text

		h, ok := Handlers[content]
		if ok {
			response := h()

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := s.Client.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}

		// Can be better to decouple 1 to 1 of message : response
		for _, v := range CatchallHandlers {
			// FIXME
			r := v(Request{content, "telegram", "", ""})
			if r != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, r)
				msg.ReplyToMessageID = update.Message.MessageID
				_, err := s.Client.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}

		// Can be better to decouple 1 to 1 of message : response
		for _, v := range CatchallExtendedHandlers {
			r := v(ExtendedMessage{Text: content})
			if r != nil {
				if r.Text != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, r.Text)
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := s.Client.Send(msg)
					if err != nil {
						log.Println(err)
					}
				}
				if r.Image != nil {
					photoFileBytes := tgbotapi.FileBytes{
						Name:  sanitizeFilename(content, ".png"),
						Bytes: r.Image,
					}
					msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photoFileBytes)
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := s.Client.Send(msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		sliced_content := strings.SplitN(content, " ", 2)
		if len(sliced_content) > 1 {
			command := sliced_content[0]
			actual_content := sliced_content[1]

			ih, ok := MsgInputHandlers[command]
			if ok {
				response := ih(Request{actual_content, "telegram", "", ""})
				if response != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := s.Client.Send(msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		m := update.Message

		if m.Photo != nil {
			// FIXME: Actual exec path never goes through here.
			targetPixels := 512 * 512
			best := &tgbotapi.PhotoSize{}
			bestSize := 100000000

			for _, v := range *m.Photo {
				pixels := v.Height * v.Width
				diff := targetPixels - pixels
				if diff < 0 {
					diff = -diff
				}
				if diff < bestSize {
					bestSize = diff
					best = &v
				}
			}
			// log.Printf("photosize %+v\n", best)
			// FIXME:
			filename := "tmp/" + best.FileID
			err := s.botDownload(best.FileID, filename)
			if err != nil {
				log.Println(err)
			}
			for _, v := range ImageHandlers {
				// log.Printf("%+v\n", m)
				c := content
				if c == "" {
					c = m.Caption
				}
				r := v(filename, Request{c, "telegram", "", ""})
				if r != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, r)
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := s.Client.Send(msg)
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

func (s *TelegramMessagePlatform) botDownload(fileId string, localFilename string) error {
	bot := s.Client
	downloadUrl, err := bot.GetFileDirectURL(fileId)
	if err != nil {
		log.Println(err)
		return err
	}
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

func (s *TelegramMessagePlatform) Send(text string) {
	if s == nil {
		return
	}
	s.SendWithOptions(text, SendOptions{})

}

func (s *TelegramMessagePlatform) SendWithOptions(text string, options SendOptions) {
	if s == nil {
		return
	}
	if options.Silent {
		err := s.ChannelMessageSilentSend("", text)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := s.ChannelMessageSend("", text)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *TelegramMessagePlatform) Close() {
}

func (s *TelegramMessagePlatform) ChannelMessageSend(channel, message string) error {
	if channel == "" {
		channel = s.DefaultChannel
	}
	channelId, ok := engineersmy.KnownTelegramChannels[channel]
	if !ok {
		log.Println("Unknown channel", channel)
		return fmt.Errorf("Unknown channel %s", channel)
	}
	msg := tgbotapi.NewMessage(int64(channelId), message)
	_, err := s.Client.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}

// ChannelMessageSilentSend is FIXME: dupe of ChannelMessageSend with DisableNotification
func (s *TelegramMessagePlatform) ChannelMessageSilentSend(channel, message string) error {
	if channel == "" {
		channel = s.DefaultChannel
	}
	channelId, ok := engineersmy.KnownTelegramChannels[channel]
	if !ok {
		log.Println("Unknown channel", channel)
		return fmt.Errorf("Unknown channel %s", channel)
	}
	msg := tgbotapi.NewMessage(int64(channelId), message)
	msg.DisableNotification = true
	_, err := s.Client.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}
