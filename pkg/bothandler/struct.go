package bothandler

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

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
type ImageListHandler func([]string, Request) string

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
var ImageListHandlers = []ImageListHandler{}
var AddMessagePlatforms = []AddMessagePlatform{}
var ActiveMessagePlatforms = []MessagePlatform{}

func RegisterMessageWithInputHandler(m string, h MessageWithInputHandler) {
	MsgInputHandlers[m] = h
}

func GetMsgInputHandler(content string) (MessageWithInputHandler, string, bool) {
	slicedContent := strings.SplitN(content, " ", 2)
	command := slicedContent[0]
	if ih, ok := MsgInputHandlers[command]; ok {
		actualContent := ""
		if len(slicedContent) > 1 {
			actualContent = strings.TrimSpace(slicedContent[1])
		}
		return ih, actualContent, true
	}

	var bestKey string
	var bestHandler MessageWithInputHandler
	for key, handler := range MsgInputHandlers {
		if strings.HasPrefix(content, key) {
			if len(key) > len(bestKey) {
				bestKey = key
				bestHandler = handler
			}
		}
	}

	if bestKey != "" {
		actualContent := strings.TrimPrefix(content, bestKey)
		actualContent = strings.TrimSpace(actualContent)
		return bestHandler, actualContent, true
	}

	return nil, "", false
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

func RegisterImageListHandler(h ImageListHandler) {
	ImageListHandlers = append(ImageListHandlers, h)
}

func RegisterMessagePlatform(m MessagePlatform) {
	ActiveMessagePlatforms = append(ActiveMessagePlatforms, m)
}

func RegisterPassiveMessagePlatform(m MessagePlatform) {
	RegisterMessagePlatform(m)
}

func Shutdown() {
	for _, v := range ActiveMessagePlatforms {
		v.Close()
	}
}

func ChannelMessageSend(channelId string, message string) error {
	var errs []error
	for _, v := range ActiveMessagePlatforms {
		err := v.ChannelMessageSend(channelId, message)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// splitMessage splits text into rune-safe chunks of at most limit bytes,
// preferring to break at whitespace.
func splitMessage(text string, limit int) []string {
	var chunks []string
	for len(text) > 0 {
		if len(text) <= limit {
			chunks = append(chunks, text)
			break
		}
		cut := limit
		// Don't split in the middle of a UTF-8 rune.
		for cut > 0 && !utf8.RuneStart(text[cut]) {
			cut--
		}
		chunk := text[:cut]
		rest := text[cut:]
		if spaceIndex := strings.LastIndex(chunk, " "); spaceIndex > 0 {
			rest = text[spaceIndex+1:]
			chunk = text[:spaceIndex]
		}
		chunks = append(chunks, chunk)
		text = rest
	}
	return chunks
}

func sanitizeFilename(f string, extension string) string {
	f = strings.ReplaceAll(f, " ", "_")
	if len(f) > 94 {
		cut := 94
		// Don't split in the middle of a UTF-8 rune.
		for cut > 0 && !utf8.RuneStart(f[cut]) {
			cut--
		}
		f = f[:cut]
	}
	filename := fmt.Sprintf("%s.%s", f, extension)
	filename, err := filenamify.Filenamify(filename, filenamify.Options{Replacement: "_"})
	if err != nil {
		log.Println(err)
		return "badfilename." + extension
	}
	return filename
}
