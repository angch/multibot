package stoic

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/angch/multibot/pkg/bothandler"
)

type StoicResponse struct {
	Id       int    `json:"id"` // cannot unmarshal number into Go struct field StoicResponse.id of type string
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
	Author   string `json:"author"`
}

func init() {
	bothandler.RegisterMessageHandler("!stoic", GetMessage)
}

/*
{"id":21,"body":"The soul becomes dyed with the color of its thoughts.","author_id":1,"author":"Marcus Aurelius"}
*/

func GetMessage() string {
	url := "https://stoicquotesapi.com/v1/api/quotes/random"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("error retrieving stoicquotesapi", err)
		return ""
	}

	defer resp.Body.Close()

	var respBody StoicResponse

	// This is when you don't want a stream, so you have a copy you can debug
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.Println("error decoding stoicquotesapi response", err, string(body))
		return ""
	}

	message := respBody.Body + " — " + respBody.Author
	return message
}
