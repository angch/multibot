package bothandler

// We could use discordbot's handler system, but we have own wrappers here to make
// multiplatform bots work.

type MessageHandler func() string

var Handlers = map[string]MessageHandler{}

func RegisterMessageHandler(m string, h MessageHandler) {
	Handlers[m] = h
}
