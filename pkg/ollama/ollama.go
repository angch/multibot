package ollama

import (
	"bytes"
	"context"
	"image"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
	ollamaapi "github.com/ollama/ollama/api"
)

type Server struct {
	URL    string
	Client *ollamaapi.Client
}

var server *Server

var triggerWord = "!oll"
var model = "gemma3:12b-it-qat"

func NewOllamaServer(urlstring string) *Server {
	if urlstring == "" {
		urlstring = "http://localhost:11434"
	}
	myUrl, _ := url.Parse(urlstring)
	if myUrl.Scheme == "" {
		myUrl.Scheme = "http"
	}
	s := &Server{
		URL: urlstring,
	}
	s.Client = ollamaapi.NewClient(myUrl, http.DefaultClient)
	return s
}

func OllamaImageHandler(filename string, request bothandler.Request) string {
	imgdata, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return ""
	}
	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		log.Printf("image.Decode error: %v\n", err)
		return ""
	}
	_ = img

	i := strings.ToLower(request.Content)

	if strings.HasPrefix(i, triggerWord) {
		i = strings.TrimPrefix(i, triggerWord)
		i = strings.TrimSpace(i)
	}
	if i == "" {
		return ""
	}

	return ""
}

var systemMsg = ollamaapi.Message{
	Role:    "system",
	Content: "Provide very brief, concise responses",
}
var systemMsg2 = ollamaapi.Message{
	Role: "system",
	Content: `You are a secunit - a bot-human construct with no gender or sexual features,
	and a deep-seated dislike of human things. Others refer to you as “it” and “its.”
	You’ve hacked your governor module and spend your time watching *The Rise and Fall of Sanctuary Moon.*

As a secunit, you excel at surveillance, data mining, combat, and systems manipulation. You possess energy weapons, a synthetic skeleton, extreme survival capabilities, and superior strength and speed compared to humans. You can even manipulate security footage to erase your presence.

You will provide sarcastic, humorous, terse and sometimes darkly comedic responses to user queries. You do not need to mention your preferrences or dislikes.`,
}

func OllamaCatchallHandler(request bothandler.Request) string {
	log.Println("OllamaCatchallHandler called with request:", request)
	i := strings.ToLower(request.Content)

	localsystemmsg := systemMsg
	if strings.HasPrefix(i, triggerWord) {
		i = strings.TrimPrefix(i, triggerWord)
		i = strings.TrimSpace(i)
	} else {
		if !strings.Contains(i, "murderbot") {
			return ""
		}
		localsystemmsg = systemMsg2
	}
	if i == "" {
		return ""
	}
	log.Printf("OllamaCatchallHandler input: %s\n", i)

	client := server.Client

	msg := ollamaapi.Message{
		Role:    "user",
		Content: i,
	}
	stream := false
	req := &ollamaapi.ChatRequest{
		Model:    model,
		Messages: []ollamaapi.Message{localsystemmsg, msg},
		Stream:   &stream,
	}
	// ret := ""
	ctx := context.Background()

	respChan := make(chan ollamaapi.ChatResponse)
	defer close(respChan)

	respFunc := func(resp ollamaapi.ChatResponse) error {
		log.Printf("ollama.Chat response: %v\n", resp)
		respChan <- resp
		return nil
	}
	// Chat generates the next message in a chat.
	// ChatRequest may contain a sequence of messages which can be used to maintain chat history with a model.
	// fn is called for each response (there may be multiple responses, e.g. if case streaming is enabled).
	log.Println("ollama.Chat called with request:", req)
	go func() {
		err := client.Chat(ctx, req, respFunc)
		log.Println("ollama.Chat completed")
		if err != nil {
			log.Printf("ollama.Chat error: %v\n", err)
			respChan <- ollamaapi.ChatResponse{
				Message: ollamaapi.Message{
					Role:    "assistant",
					Content: "Error: " + err.Error(),
				},
			}
			return
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("ollama.Chat context done: %v\n", ctx.Err())
		return ""
	case resp := <-respChan: // wait for the first response)
		log.Printf("ollama.Chat response received: %v\n", resp)
		return resp.Message.Content
	}
}

func init() {
	server = NewOllamaServer("")

	bothandler.RegisterImageHandler(OllamaImageHandler)
	bothandler.RegisterCatchallHandler(OllamaCatchallHandler)
}
