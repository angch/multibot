package standarddiffusion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

var sd_url *url.URL

func init() {
	bothandler.RegisterCatchallExtendeHandler(GetMessage)
	sd_urlString := os.Getenv("SD_URL")
	if sd_urlString == "" {
		log.Fatal("Need valid env SD_URL")
	}
	sd, err := url.Parse(sd_urlString)
	if err != nil {
		log.Fatal("Need valid env SD_URL")
	}
	if sd.Scheme != "http" && sd.Scheme != "https" {
		log.Fatal("Need valid env SD_URL")
	}
	sd_url = sd
}

/*
{"id":21,"body":"The soul becomes dyed with the color of its thoughts.","author_id":1,"author":"Marcus Aurelius"}
*/

type JsonResponse struct {
	Error string `json:"error"`
}

func GetMessage(input bothandler.ExtendedMessage) *bothandler.ExtendedMessage {
	i := strings.ToLower(input.Text)

	if strings.HasPrefix(i, "!sd ") {
		i = i[4:]
	} else {
		return nil
	}

	u := *sd_url
	q := u.Query()
	q.Set("q", i)
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
		return &bothandler.ExtendedMessage{
			Text:  "Zzzz server is sleeping",
			Image: nil,
		}
	}

	if len(body) > 0 && body[0] == '{' {
		// It's likely an error
		msg := JsonResponse{}
		err := json.Unmarshal(body, &msg)
		if msg.Error != "" && err == nil {
			return &bothandler.ExtendedMessage{
				Text:  msg.Error,
				Image: nil,
			}
		} else {
			return &bothandler.ExtendedMessage{
				Text:  "Zzzz server is sleeping",
				Image: nil,
			}
		}
	}

	return &bothandler.ExtendedMessage{
		Text:  "",
		Image: body,
	}
}
