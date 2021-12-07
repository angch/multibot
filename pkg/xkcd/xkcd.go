package xkcd

import (
	"fmt"
	"strconv"

	"github.com/angch/discordbot/pkg/bothandler"
)

func init() {
	bothandler.RegisterMessageWithInputHandler("!xkcd", GetXKCD)
	bothandler.RegisterMessageWithInputHandler("!explainxkcd", GetXKCDExplained)
}

func sanitize(input string) int {
	n, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	return n
}

func GetXKCD(request bothandler.Request) string {
	num := sanitize(request.Content)
	if num <= 0 {
		return ""
	}

	resp_url := fmt.Sprintf("https://www.xkcd.com/%d/", num)
	message := resp_url
	return message
}

func GetXKCDExplained(request bothandler.Request) string {
	num := sanitize(request.Content)
	if num <= 0 {
		return ""
	}

	resp_url := fmt.Sprintf("https://www.explainxkcd.com/wiki/index.php?title=%d", num)
	message := resp_url
	return message
}
