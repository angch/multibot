package meme

import (
	"fmt"
	"strings"
	"math/rand"
	"time"

	"github.com/angch/discordbot/pkg/bothandler"
)

var gifLinks = []string {
	"https://tenor.com/La2h.gif",
	"https://tenor.com/Zx43.gif",
	"https://tenor.com/bCcZV.gif",
	"https://tenor.com/14jR.gif",
}

var responseMap = map[string]string{
	"omae wa mou shindeiru"	: "Nani?!",
	"omae wa mo shindeiru"	: "Nani?!",
	"omaewamoshindeiru"	: "Nani?!",
	"お前はもう死んでいる"	: "なに?!",
}

func init() {
	bothandler.RegisterCatchallHandler(responseHandler)
}

func responseHandler(input string) string {
	i := strings.ToLower(input)
	response := "";
	
	rand.Seed(time.Now().Unix())
	gifLink := gifLinks[rand.Intn(len(gifLinks))]

	value, exists := responseMap[i]
	if exists {
		response = fmt.Sprintf("%s %s", value, gifLink)
	}

	return response
}
