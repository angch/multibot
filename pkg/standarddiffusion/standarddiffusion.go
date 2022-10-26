package standarddiffusion

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

var sd_url url.URL

func init() {
	bothandler.RegisterCatchallExtendeHandler(GetMessage)
	sd_url := os.Getenv("SD_URL")
	if sd_url == "" {
		log.Fatal("Need valid env SD_URL")
	}
	_, err := url.Parse(sd_url)
	if err != nil {
		log.Fatal("Need valid env SD_URL")
	}
}

/*
{"id":21,"body":"The soul becomes dyed with the color of its thoughts.","author_id":1,"author":"Marcus Aurelius"}
*/

func GetMessage(input bothandler.ExtendedMessage) *bothandler.ExtendedMessage {
	i := strings.ToLower(input.Text)

	if strings.HasPrefix(i, "!sd ") {
		i = i[4:]
	} else {
		return nil
	}

	u := sd_url
	q := u.Query()
	q.Add("q", i)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		fmt.Println("error retrieving", err)
		return &bothandler.ExtendedMessage{
			Text:  err.Error(),
			Image: nil,
		}
	}

	defer resp.Body.Close()

	// This is when you don't want a stream, so you have a copy you can debug
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return &bothandler.ExtendedMessage{
		Text:  "",
		Image: body,
	}
}
