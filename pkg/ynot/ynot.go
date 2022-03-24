package ynot

import (
	"math/rand"
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

var triggers = [][]string{
	{"why don't", "y not"},
	{"just", "use", "try", "pay", "buy"},
	{"?", "in"},
}

// https://www.brendangregg.com/blog/2022-03-19/why-dont-you-use.html
var excuses = []string{
	"It performs poorly.",
	"It is too expensive.",
	"It is not open source.",
	"It lacks features.",
	"It lacks a community.",
	"It lacks debug tools.",
	"It has serious bugs.",
	"It is poorly documented.",
	"It lacks timely security fixes.",
	"It lacks subject matter expertise.",
	"It's developed for the wrong audience.",
	"Our custom internal solution is good enough.",
	"Its longevity is uncertain: Its startup may be dead or sold soon.",
	"We know under NDA of a better solution.",
	"We know other bad things under NDA.",
	"Key contributors told us it was doomed.",
	"It made sense a decade ago but doesn't today.",
	"It made false claims in articles/blogs/talks and has lost credibility.",
	"It tolerates brilliant jerks and has no effective CoC.",
	"Our lawyers won't accept its T&Cs or license.",
}

func init() {
	bothandler.RegisterCatchallHandler(YNotHandler)
}

func ynot(i string) bool {
	i = strings.ToLower(i)

	count := 0
	for _, v := range triggers {
		for _, v2 := range v {
			idx := strings.Index(i, v2)

			if idx >= 0 {
				count++
				i = i[idx:]
			}
		}
	}

	return count >= 3

}

func YNotHandler(r bothandler.Request) string {
	input := r.Content
	if ynot(input) {
		random := rand.Intn(len(excuses))
		return excuses[random]
	}

	return ""
}
