package bothandler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Implements MessagePlatform
type SlackMessagePlatform struct {
	SlackWebHook string
}

func NewMessagePlatformFromSlack(slackwebhook string) *SlackMessagePlatform {
	return &SlackMessagePlatform{
		SlackWebHook: slackwebhook,
	}
}

func (s *SlackMessagePlatform) Send(text string) {
	content := bytes.NewBuffer([]byte(fmt.Sprintf("{\"text\":\"%s\"}", text)))
	_, err := http.Post(s.SlackWebHook, "Content-type: application/json", content)
	if err != nil {
		log.Println(err)
	}
}
