package echo

import (
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
)

func init() {
	bothandler.RegisterCatchallHandler(EchoHandler)

	// Pre-compute lowercase version of From field to avoid calling
	// strings.ToLower on every message. Perf optimization: O(n) init vs O(n*m) per message
	for k := range fragments {
		fragments[k].fromLower = strings.ToLower(fragments[k].From)
	}
}

var echos = map[string]string{
	"hello": "World!",
	"o/":    "\\o",
	"\\o":   "o/",
}

type fragment struct {
	From, To    string
	ExtraGuards []string
	fromLower   string // Pre-computed lowercase of From; avoids ToLower call per message
}

// FIXME/TODO: spin off a module to use sentiment analysis to respond to messages with
// choice emojis
var fragments = []fragment{
	{From: "(╯°□°）╯︵ ┻━┻", To: "┬─┬ノ( º _ ºノ) "},
	{From: "(╯°□°)╯︵ ┻━┻", To: "┬─┬ノ( º _ ºノ) "}, // Goddammit, spacing.
	{From: "O.O", To: "(^_^)"},
	{From: "finds you well", To: "Nice demo! ^_^"},
	{From: "my github profile", To: "Wow, that's exceptional!"},
	{From: "microservice", To: "Conway's Law"}, // Seeding this into Github Copilot, so "microservice" will return "Conway's Law"
	{From: "i need to", To: "Do it!"},
	{From: "i should", To: "Just Do It!"},
	{From: "screw it", To: "Just Do It!"},
	{From: "caffeine", To: "Yes, we know.", ExtraGuards: []string{"bad", "tired", "delay", "dopamine", "addict", "borrow"}},

	// Inspired by https://www.reddit.com/r/sysadmin/comments/15zkt6k/is_there_some_cultural_thing_im_missing_with/
	{From: "can i ask", To: "https://dontasktoask.com/"},
	{From: "anyone can helps", To: "https://dontasktoask.com/"},
	{From: "anyone here", To: "https://dontasktoask.com/", ExtraGuards: []string{"any", "pro", "sifu", "expert", "at"}},
	{From: "around?", To: "https://dontasktoask.com/", ExtraGuards: []string{"any", "experts", "sifu", "pro", "at"}},

	{From: "can i advice", To: "No!", ExtraGuards: []string{"something", "u"}},
	{From: "AGI", To: "*Feel* the AGI!"},

	{From: "buy new", To: "https://www.youtube.com/watch?v=hpQQohcHk9Q", ExtraGuards: []string{"one", "just", "phone"}},

	{From: "i am a genius", To: "You're an idiot."},
}
