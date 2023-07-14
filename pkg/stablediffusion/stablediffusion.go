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

	"github.com/angch/multibot/pkg/bothandler"
	"github.com/angch/multibot/pkg/stablediffusion/sdapi"
)

var sd_url *url.URL
var sdapi_server *sdapi.Server

func init() {
	bothandler.RegisterCatchallExtendeHandler(GetMessage)
	sdapi_url, sd_urlString := os.Getenv("SDAPI_URL"), os.Getenv("SD_URL")

	if sd_urlString == "" && sdapi_url == "" {
		log.Fatal("Need valid env SD_URL or SDAPI_URL")
	}
	if sdapi_url != "" {
		sdapi_server = sdapi.NewServer(sdapi_url)
		sd_url = nil
	} else {
		sd, err := url.Parse(sd_urlString)
		if err != nil {
			log.Fatal("Need valid env SD_URL")
		}
		if sd.Scheme != "http" && sd.Scheme != "https" {
			log.Fatal("Need valid env SD_URL")
		}
		sd_url = sd
		sdapi_server = nil
	}
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

	var body []byte
	var err error
	if sdapi_server != nil {
		log.Println("sdapi")
		body, err = sdapi_server.Txt2Img(i)
	} else if sd_url != nil {
		u := *sd_url
		q := u.Query()
		q.Set("q", i)
		u.RawQuery = q.Encode()
		var resp *http.Response
		resp, err = http.Get(u.String())

		if err != nil {
			fmt.Println("error retrieving", err)
			return &bothandler.ExtendedMessage{
				Text:  err.Error(),
				Image: nil,
			}
		}

		defer resp.Body.Close()

		// This is when you don't want a stream, so you have a copy you can debug
		body, err = io.ReadAll(resp.Body)
	}
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
