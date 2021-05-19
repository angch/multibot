package bothandler

import (
	"fmt"
	"log"
	"sync"

	"github.com/angch/discordbot/pkg/engineersmy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Implements MessagePlatform
// New way, using Slack Rtm
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
		log.Printf("%+v\n", update)
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s %d %d", update.Message.From.UserName, update.Message.Text, update.Message.From.ID, update.Message.Chat.ID)
		s.KnownUsersLock.Lock()
		s.KnownUsers[update.Message.From.UserName] = *update.Message.From
		s.KnownUsersLock.Unlock()

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID
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
			r := v(content)
			if r != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, r)
				msg.ReplyToMessageID = update.Message.MessageID
				_, err := s.Client.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (s *TelegramMessagePlatform) Send(text string) {
	if s == nil {
		return
	}
	err := s.ChannelMessageSend("", text)
	if err != nil {
		log.Println(err)
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
