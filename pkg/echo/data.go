package echo

import (
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

func init() {
	bothandler.RegisterCatchallHandler(EchoHandler)

	for k, v := range fragments {
		vl := strings.ToLower(v.From)
		if vl != v.From {
			fragments[k].From = vl
		}
	}
}

var echos = map[string]string{
	"hello": "World!",
	"o/":    "\\o",
	"\\o":   "o/",
}

type fragment struct {
	From, To string
}

// FIXME/TODO: spin off a module to use sentiment analysis to respond to messages with
// choice emojis
var fragments = []fragment{
	{"(╯°□°）╯︵ ┻━┻", "┬─┬ノ( º _ ºノ) "},
	{"O.O", "(^_^)"},
	{"omaewamoushindeiru", "nani?!"},
	{"お前はもう死んでいる", "Haven’t merged yet, please try again later"},
}
