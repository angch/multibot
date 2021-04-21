package echo

// (ꈍᴗꈍ)
import (
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

func init() {
	bothandler.RegisterCatchallHandler(UwuHandler)

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
}

func UwuHandler(input string) string {
	i := strings.ToLower(input)
	if strings.Contains(i, "uwu") || strings.Contains(input, "(ꈍᴗꈍ)") {
		return "(ꈍᴗꈍ)"
	}

	r, ok := echos[i]
	if ok {
		return r
	}

	for _, v := range fragments {
		if strings.Contains(i, v.From) {
			return v.To
		}
	}
	return ""
}
