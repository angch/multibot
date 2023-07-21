package bothandler

import (
	"fmt"
	"log"
	"strings"

	"github.com/flytam/filenamify"
)

// We could use multibot's handler system, but we have own wrappers here to make
// multiplatform bots work.

type Request struct {
	Content  string
	Platform string
	// ClientId string
	Channel string
	From    string
}

type MessageHandler func() string
type MessageWithInputHandler func(Request) string
type CatchallHandler func(Request) string
type ImageHandler func(string, Request) string

type SendOptions struct {
	Silent bool
}

type ExtendedMessage struct {
	Text  string
	Image []byte
}
type CatchallExtendedHandler func(ExtendedMessage) *ExtendedMessage

type MessagePlatform interface {
	Send(string)
	SendWithOptions(string, SendOptions)
	ProcessMessages()
	Close()
	ChannelMessageSend(channel string, message string) error
}

type AddMessagePlatform func(MessagePlatform)

var Handlers = map[string]MessageHandler{}
var MsgInputHandlers = map[string]MessageWithInputHandler{}
var CatchallHandlers = []CatchallHandler{}
var CatchallExtendedHandlers = []CatchallExtendedHandler{}
var ImageHandlers = []ImageHandler{}
var AddMessagePlatforms = []AddMessagePlatform{}
var ActiveMessagePlatforms = []MessagePlatform{}

func RegisterMessageWithInputHandler(m string, h MessageWithInputHandler) {
	MsgInputHandlers[m] = h
}

func RegisterMessageHandler(m string, h MessageHandler) {
	Handlers[m] = h
}

func RegisterCatchallHandler(h CatchallHandler) {
	CatchallHandlers = append(CatchallHandlers, h)
}
func RegisterCatchallExtendeHandler(h CatchallExtendedHandler) {
	CatchallExtendedHandlers = append(CatchallExtendedHandlers, h)
}

func RegisterImageHandler(h ImageHandler) {
	ImageHandlers = append(ImageHandlers, h)
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
		v(m)
	}
	ActiveMessagePlatforms = append(ActiveMessagePlatforms, m)
}

func RegisterPassiveMessagePlatform(m MessagePlatform) {
	ActiveMessagePlatforms = append(ActiveMessagePlatforms, m)
}

func Shutdown() {
	for _, v := range ActiveMessagePlatforms {
		v.Close()
	}
}

func ChannelMessageSend(channelId string, message string) error {
	for _, v := range ActiveMessagePlatforms {
		err := v.ChannelMessageSend(channelId, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func sanitizeFilename(f string, extension string) string {
	f = strings.ReplaceAll(f, " ", "_")
	if len(f) > 94 {
		f = f[:94]
	}
	filename := fmt.Sprintf("%s.%s", f, extension)
	filename, err := filenamify.Filenamify(filename, filenamify.Options{Replacement: "_"})
	if err != nil {
		log.Println(err)
		return "badfilename." + extension
	}
	return filename
}
