package echo

import (
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

func EchoHandler(request bothandler.Request) string {
	i := strings.ToLower(request.Content)
	uwu := uwucheck(i)
	if uwu != "" {
		return uwu
	}

	r, ok := echos[i]
	if ok {
		return r
	}

	for _, v := range fragments {
		if strings.Contains(i, v.From) {
			if len(i) <= len(v.From)*7 {
				return v.To
			}
		}
	}
	return ""
}
