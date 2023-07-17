package spacetraders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (a *SpaceTraders) RegisterAgent(ctx context.Context, pc PlatformChannel, agentRequest RegisterAgentRequest) (*RegisterAgentResponse, error) {
	// https://api.spacetraders.io/v2/register
	posturl := "https://api.spacetraders.io/v2/register"
	body, err := json.Marshal(agentRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "multibot")
	client := a.HttpClient
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	a.LogRequest(ctx, &RequestLog{
		Platform:           pc.Platform,
		Channel:            pc.Channel,
		Type:               "RegisterAgent",
		URL:                posturl,
		Data:               string(body),
		Response:           string(out),
		ResponseStatusCode: res.StatusCode,
	}, req)

	registerAgentResponse := RegisterAgentResponse{}
	err = json.Unmarshal(out, &registerAgentResponse)
	if err != nil {
		return nil, err
	}

	if registerAgentResponse.Error != nil {
		return nil, fmt.Errorf("%s", registerAgentResponse.Error.Message)
	}
	return &registerAgentResponse, nil
}
