package xkcd

import (
	"fmt"
	"net/http"

	"github.com/angch/discordbot/pkg/bothandler"
)

func init() {
	bothandler.RegisterMessageWithInputHandler("!xkcd", GetXKCD)
	bothandler.RegisterMessageWithInputHandler("!explainxkcd", GetXKCDExplained)
}

func GetXKCD(request bothandler.Request) string {
	input := request.Content

	base_url := "https://www.xkcd.com/"

	resp_url := fmt.Sprintf("%s%s/", base_url, input)

	resp, err := http.Get(resp_url)

	if err != nil {
		fmt.Println("Error retrieving explainxkcd link: ", err)
		return ""
	}

	defer resp.Body.Close()

	message := resp_url
	return message
}

func GetXKCDExplained(request bothandler.Request) string {
	input := request.Content

	base_url := "https://www.explainxkcd.com/wiki/index.php/"

	resp_url := fmt.Sprintf("%s%s", base_url, input)

	resp, err := http.Get(resp_url)

	if err != nil {
		fmt.Println("Error retrieving explainxkcd link: ", err)
		return ""
	}

	defer resp.Body.Close()

	message := resp_url
	return message
}
