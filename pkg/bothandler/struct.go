package bothandler

// We could use discordbot's handler system, but we have own wrappers here to make
// multiplatform bots work.

type MessageHandler func() string
type CatchallHandler func(string) string

var Handlers = map[string]MessageHandler{}
var CatchallHandlers = []CatchallHandler{}

func RegisterMessageHandler(m string, h MessageHandler) {
	Handlers[m] = h
}

func RegisterCatchallHandler(h CatchallHandler) {
	CatchallHandlers = append(CatchallHandlers, h)
}
