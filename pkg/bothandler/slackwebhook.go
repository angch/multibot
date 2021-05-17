package bothandler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Implements MessagePlatform
// Old way, using webhooks
type SlackWebhookMessagePlatform struct {
	SlackWebHook string
}

func NewMessagePlatformFromSlackWebhook(slackwebhook string) *SlackWebhookMessagePlatform {
	return &SlackWebhookMessagePlatform{
		SlackWebHook: slackwebhook,
	}
}

func (s *SlackWebhookMessagePlatform) Send(text string) {
	content := bytes.NewBuffer([]byte(fmt.Sprintf("{\"text\":\"%s\"}", text)))
	_, err := http.Post(s.SlackWebHook, "Content-type: application/json", content)
	if err != nil {
		log.Println(err)
	}
}

func (s *SlackWebhookMessagePlatform) ProcessMessages() {
}

func (s *SlackWebhookMessagePlatform) Close() {
}

func (s *SlackWebhookMessagePlatform) ChannelMessageSend(channelId, message string) error {
	return nil
}
