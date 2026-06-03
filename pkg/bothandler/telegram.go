package bothandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/angch/multibot/pkg/engineersmy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type mediaGroupEntry struct {
	chatID    int64
	replyToID int
	caption   string
	from      string
	filenames []string
	timer     *time.Timer
}

// Implements MessagePlatform
type TelegramMessagePlatform struct {
	Client         *tgbotapi.BotAPI
	ChannelId      map[string]string
	KnownUsers     map[string]tgbotapi.User
	KnownUsersLock sync.RWMutex
	DefaultChannel string
	mediaGroupsMu  sync.Mutex
	mediaGroups    map[string]*mediaGroupEntry
	// Me             *tgbotapi.User // Superflous, get it from Client.Self
}

func NewMessagePlatformFromTelegram(telegrambottoken string) (*TelegramMessagePlatform, error) {
	bot, err := tgbotapi.NewBotAPI(telegrambottoken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Connected to Telegram on account %s", bot.Self.UserName)

	return &TelegramMessagePlatform{
		Client:      bot,
		ChannelId:   map[string]string{},
		KnownUsers:  map[string]tgbotapi.User{},
		mediaGroups: map[string]*mediaGroupEntry{},
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
						Name:  sanitizeFilename(content, ".png"),
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
			chatIDStr := fmt.Sprintf("%d", m.Chat.ID)

			if m.MediaGroupID != "" {
				// Album: buffer photos until they stop arriving, then dispatch.
				key := fmt.Sprintf("%d_%s", m.Chat.ID, m.MediaGroupID)
				s.bufferAlbumPhoto(key, filename, caption, m.From.UserName, m.Chat.ID, m.MessageID)
			} else {
				// Single image: call legacy single-file handlers and list handlers.
				req := Request{caption, "telegram", chatIDStr, m.From.UserName}
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
	if err != nil || get == nil || get.Body == nil {
		log.Println(err)
		return err
	}
	reader := get.Body
	defer reader.Close()

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

// bufferAlbumPhoto accumulates photos from a Telegram media group (album).
// It resets a 800 ms timer on each arrival; when the timer fires all buffered
// filenames are dispatched to ImageListHandlers in one call.
func (s *TelegramMessagePlatform) bufferAlbumPhoto(key, filename, caption, from string, chatID int64, replyToID int) {
	s.mediaGroupsMu.Lock()
	defer s.mediaGroupsMu.Unlock()

	entry, exists := s.mediaGroups[key]
	if !exists {
		entry = &mediaGroupEntry{
			chatID:    chatID,
			replyToID: replyToID,
			from:      from,
		}
		s.mediaGroups[key] = entry
	}
	entry.filenames = append(entry.filenames, filename)
	if caption != "" && entry.caption == "" {
		entry.caption = caption
	}

	if entry.timer != nil {
		entry.timer.Stop()
	}
	captured := entry
	entry.timer = time.AfterFunc(800*time.Millisecond, func() {
		s.mediaGroupsMu.Lock()
		delete(s.mediaGroups, key)
		s.mediaGroupsMu.Unlock()

		req := Request{
			Content:  captured.caption,
			Platform: "telegram",
			Channel:  fmt.Sprintf("%d", captured.chatID),
			From:     captured.from,
		}
		for _, h := range ImageListHandlers {
			r := h(captured.filenames, req)
			if r != "" {
				msg := tgbotapi.NewMessage(captured.chatID, r)
				msg.ReplyToMessageID = captured.replyToID
				if _, err := s.Client.Send(msg); err != nil {
					log.Println(err)
				}
			}
		}
	})
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
		return fmt.Errorf("unknown channel %s", channel)
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
		return fmt.Errorf("unknown channel %s", channel)
	}
	msg := tgbotapi.NewMessage(int64(channelId), message)
	msg.DisableNotification = true
	_, err := s.Client.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}
