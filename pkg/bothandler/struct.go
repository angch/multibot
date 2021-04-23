package bothandler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/angch/discordbot/pkg/engineersmy"
	"github.com/bwmarrin/discordgo"
)

// We could use discordbot's handler system, but we have own wrappers here to make
// multiplatform bots work.

type MessageHandler func() string
type CatchallHandler func(string) string

type MessagePlatform interface {
	Send(string)
}

type AddMessagePlatform func(MessagePlatform)

var Handlers = map[string]MessageHandler{}
var CatchallHandlers = []CatchallHandler{}
var AddMessagePlatforms = []AddMessagePlatform{}

func RegisterMessageHandler(m string, h MessageHandler) {
	Handlers[m] = h
}

func RegisterCatchallHandler(h CatchallHandler) {
	CatchallHandlers = append(CatchallHandlers, h)
}

// Yes, weird. All the modules register themselves,
// Then all the platforms (discord, telegram, slack) calls these
// to tell the modules these platforms are available for them
// to call asynchronously.
// We need these to get the moduls to run themselves asychronously
// apart from the main messaging platform event loop
func RegisterPlatformRegisteration(h AddMessagePlatform) {
	AddMessagePlatforms = append(AddMessagePlatforms, h)
}

func RegisterMessagePlatform(m MessagePlatform) {
	for _, v := range AddMessagePlatforms {
		// log.Printf("Registering %+v\n", m)
		v(m)
	}
}

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

type DevMessagePlatform struct {
}

func NewMessagePlatformFromDev() *DevMessagePlatform {
	return &DevMessagePlatform{}
}

func (s *DevMessagePlatform) Send(text string) {
	log.Println(text)
}
