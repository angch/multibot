package ynot

import (
	"math/rand"
	"strings"
	"sync"

	"github.com/angch/discordbot/pkg/bothandler"
)

var triggers = [][]string{
	{"why don't", "y not", "why dont"},
	{"just", "use", "try", "pay", "buy", "pick"},
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

var randomBufferLock = sync.Mutex{}
var randomBuffer = []int{}
var randomBufferIdx = 0

func init() {
	bothandler.RegisterCatchallHandler(YNotHandler)
	randomBuffer = make([]int, len(excuses)/2)
	for i := 0; i < len(randomBuffer); i++ {
		randomBuffer[i] = -1
	}
}

func ynot(i string) bool {
	if len(i) > 50 {
		return false
	}
	i = strings.ToLower(i)

	count := 0
	for k, v := range triggers {
		for _, v2 := range v {
			idx := strings.Index(i, v2)

			if idx >= 0 {
				if k > 0 {
					if idx >= 100 {
						continue
					}
				}
				count++
				i = i[idx:]
			}
		}
	}

	return count >= 3
}

// getShuffleRandom returns a random number, but minimizes repeats
// https://www.ripleys.com/weird-news/is-shuffle-random/
func getShuffleRandom() int {
	r := -1
	randomBufferLock.Lock()
	defer randomBufferLock.Unlock()

again:
	for retries := 0; retries < 20; retries++ {
		r = rand.Intn(len(excuses))
		for i := 0; i < len(randomBuffer); i++ {
			if randomBuffer[i] == r {
				continue again
			}
		}
		break
	}
	randomBuffer[randomBufferIdx] = r
	randomBufferIdx = (randomBufferIdx + 1) % len(randomBuffer)

	return r
}

func YNotHandler(r bothandler.Request) string {
	input := r.Content
	if ynot(input) {
		random := getShuffleRandom()
		return excuses[random]
	}

	return ""
}
