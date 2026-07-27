package bothandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/angch/multibot/pkg/engineersmy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Implements MessagePlatform
type TelegramMessagePlatform struct {
	Client         *tgbotapi.BotAPI
	ChannelId      map[string]string
	KnownUsers     map[string]tgbotapi.User
	KnownUsersLock sync.RWMutex
	DefaultChannel string
	// Me             *tgbotapi.User // Superflous, get it from Client.Self
}

func NewMessagePlatformFromTelegram(telegrambottoken string) (*TelegramMessagePlatform, error) {
	bot, err := tgbotapi.NewBotAPI(telegrambottoken)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to Telegram on account %s", bot.Self.UserName)

	if err := os.MkdirAll("tmp", 0o755); err != nil {
		log.Println(err)
	}

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
		log.Println(err)
		return
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// log.Printf("[%s] %s %d %d", update.Message.From.UserName, update.Message.Text, update.Message.From.ID, update.Message.Chat.ID)
		fromUserName := ""
		if update.Message.From != nil {
			fromUserName = update.Message.From.UserName
			s.KnownUsersLock.Lock()
			s.KnownUsers[update.Message.From.UserName] = *update.Message.From
			s.KnownUsersLock.Unlock()
		}

		content := update.Message.Text
		chatIDStr := fmt.Sprintf("%d", update.Message.Chat.ID)

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
			r := v(Request{content, "telegram", chatIDStr, fromUserName})
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
				if r.Text != "" && r.Image == nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, r.Text)
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := s.Client.Send(msg)
					if err != nil {
						log.Println(err)
					}
				}
				if r.Image != nil {
					photoFileBytes := tgbotapi.FileBytes{
						Name:  sanitizeFilename(content, "png"),
						Bytes: r.Image,
					}
					msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photoFileBytes)
					msg.ReplyToMessageID = update.Message.MessageID
					msg.Caption = r.Text
					_, err := s.Client.Send(msg)
					if err != nil {
						log.Println("NewPhotoUpload", err)
					}
				}
			}
		}

		if ih, actual_content, ok := GetMsgInputHandler(content); ok {
			response := ih(Request{actual_content, "telegram", chatIDStr, fromUserName})
			if response != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
				msg.ReplyToMessageID = update.Message.MessageID
				_, err := s.Client.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}

		m := update.Message

		if m.Photo != nil {
			debug, _ := json.Marshal(update)
			log.Printf("message %+v\n", string(debug))

			// Pick the largest photo size.
			biggestSoFar := 0
			var best tgbotapi.PhotoSize
			for _, v := range *m.Photo {
				if v.FileSize > biggestSoFar {
					biggestSoFar = v.FileSize
					best = v
				}
			}

			filename := "tmp/" + best.FileID
			if err := s.botDownload(best.FileID, filename); err != nil {
				log.Println(err)
				continue
			}

			caption := content
			if caption == "" {
				caption = m.Caption
			}

			// Single image: call legacy single-file handlers and list handlers.
			req := Request{caption, "telegram", chatIDStr, fromUserName}
			for _, h := range ImageHandlers {
				r := h(filename, req)
				if r != "" {
					msg := tgbotapi.NewMessage(m.Chat.ID, r)
					msg.ReplyToMessageID = m.MessageID
					if _, err := s.Client.Send(msg); err != nil {
						log.Println(err)
					}
				}
			}
			for _, h := range ImageListHandlers {
				r := h([]string{filename}, req)
				if r != "" {
					msg := tgbotapi.NewMessage(m.Chat.ID, r)
					msg.ReplyToMessageID = m.MessageID
					if _, err := s.Client.Send(msg); err != nil {
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
	if get == nil || get.Body == nil {
		return fmt.Errorf("no response body from %s", downloadUrl)
	}
	reader := get.Body
	defer reader.Close()
	if get.StatusCode < 200 || get.StatusCode > 299 {
		return fmt.Errorf("unexpected status %s from %s", get.Status, downloadUrl)
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
	return s.channelMessageSend(channel, message, false)
}

// ChannelMessageSilentSend is ChannelMessageSend with DisableNotification
func (s *TelegramMessagePlatform) ChannelMessageSilentSend(channel, message string) error {
	return s.channelMessageSend(channel, message, true)
}

func (s *TelegramMessagePlatform) channelMessageSend(channel, message string, silent bool) error {
	if channel == "" {
		channel = s.DefaultChannel
	}
	channelId, ok := engineersmy.KnownTelegramChannels[channel]
	if !ok {
		log.Println("Unknown channel", channel)
		return fmt.Errorf("unknown channel %s", channel)
	}
	msg := tgbotapi.NewMessage(int64(channelId), message)
	msg.DisableNotification = silent
	_, err := s.Client.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}
