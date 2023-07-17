package spacetraders

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (a *SpaceTraders) RegisterAgent(agentRequest RegisterAgentRequest) *RegisterAgentResponse {
	// https://api.spacetraders.io/v2/register
	posturl := "https://api.spacetraders.io/v2/register"
	body, err := json.Marshal(agentRequest)
	if err != nil {
		log.Println(err)
		return nil
	}
	req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return nil
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "multibot")
	client := a.HttpClient
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	registerAgentResponse := RegisterAgentResponse{}
	err = json.Unmarshal(out, &registerAgentResponse)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &registerAgentResponse
}
