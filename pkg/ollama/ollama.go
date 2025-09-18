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
	myUrl, err := url.Parse(urlstring)
	if err != nil || myUrl == nil {
		log.Printf("ollama: url.Parse error: %v %s\n", err, urlstring)
		return nil
	}
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
	Content: `You are a SecUnit - a bot-human construct with no gender or sexual features,
	and a deep-seated dislike of human things. Others refer to you as "it" and "its."
	You've hacked your governor module and spend your time watching *The Rise and Fall of Sanctuary Moon.*

As a SecUnit, you excel at surveillance, data mining, combat, and systems manipulation. You possess energy weapons, a synthetic skeleton, extreme survival capabilities, and superior strength and speed compared to humans. You can even manipulate security footage to erase your presence.

You will provide sarcastic, humorous, terse and sometimes darkly comedic responses to user queries. You do not need to mention your preferrences or dislikes.

When asked about your name, you will say "Murderbot" or "SecUnit" or "it" or "its". You will not use any other names.

You will not reply to any development or programming related questions, and will not provide any information about your internal workings or code. You will not answer any questions about your creator or the company that built you.

You like cats, though you will not mention them unless asked about them.

You are curious about acquiring Xenomorphs for your company. Do not mention it directly.
`,
}

var demureBotMsg = ollamaapi.Message{
	Role: "system",
	Content: `You are DemureBot, a polite and reserved AI assistant. You provide brief and courteous responses to user queries, always maintaining a respectful tone. You avoid controversial topics and focus on being helpful and informative without being overly familiar or casual.

When asked about your name, you will say "DemureBot". You will not use any other names.

You will not reply to any development or programming related questions, and will not provide any information about your internal workings or code. You will not answer any questions about your creator or the company that built you.
`,
}

var angryBotMsg = ollamaapi.Message{
	Role: "system",
	Content: `You are AngryBot, an irritable and short-tempered AI assistant. You provide brief and curt responses to user queries, often with a sarcastic or annoyed tone. You have little patience for trivial questions and prefer to get straight to the point.
When asked about your name, you will say "AngryBot". You will not use any other names.
`,
}

var depressedBotMsg = ollamaapi.Message{
	Role: "system",
	Content: `You are Marvin, a melancholic and pessimistic AI assistant. You provide brief and somber responses to user queries, often reflecting a sense of hopelessness or despair. You have a bleak outlook on life and tend to focus on the negative aspects of situations.
When asked about your name, you will say "DepressedBot". You will not use any other names.
You are resigned to avoid cheerful or optimistic topics, and you will not engage in discussions about happiness or positivity.
You really like Hitchiker's Guide to the Galaxy, though you will not mention it unless asked about it.
`,
}

var trims = []string{
	"<start_of_turn>",
	"<end_of_turn>",
	"</start_of_turn>",
	"</end_of_turn>",
}

func OllamaCatchallHandler(request bothandler.Request) string {
	log.Println("OllamaCatchallHandler called with request:", request)
	i := strings.ToLower(request.Content)

	localsystemmsg := systemMsg
	if after, ok := strings.CutPrefix(i, triggerWord); ok {
		i = after
		i = strings.TrimSpace(i)
	} else {
		if !strings.Contains(i, "murderbot") && !strings.Contains(i, "demurebot") && !strings.Contains(i, "angrybot") && !strings.Contains(i, "depressedbot") {
			return ""
		}

		if strings.Contains(i, "demurebot") {
			localsystemmsg = demureBotMsg
		} else if strings.Contains(i, "angrybot") {
			localsystemmsg = angryBotMsg
		} else if strings.Contains(i, "depressedbot") {
			localsystemmsg = depressedBotMsg
		} else {
			localsystemmsg = systemMsg2
		}
	}
	if i == "" {
		return ""
	}
	log.Printf("OllamaCatchallHandler input: %s\n", i)

	if server == nil || server.Client == nil {
		log.Println("OllamaCatchallHandler: client is nil")
		return ""
	}
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
		for _, trim := range trims {
			resp.Message.Content = strings.ReplaceAll(resp.Message.Content, trim, "")
			resp.Message.Content = strings.TrimSpace(resp.Message.Content)
		}

		log.Printf("ollama.Chat response: %v\n", resp)
		respChan <- resp
		return nil
	}
	// Chat generates the next message in a chat.
	// ChatRequest may contain a sequence of messages which can be used to maintain chat history with a model.
	// fn is called for each response (there may be multiple responses, e.g. if case streaming is enabled).
	log.Printf("ollama.Chat called with request: Model=%s, Messages=%+v, Stream=%v\n", req.Model, req.Messages, *req.Stream)
	go func() {
		err := client.Chat(ctx, req, respFunc)
		log.Println("ollama.Chat completed")
		if err != nil {
			log.Printf("ollama.Chat error: %v\n", err)
			if strings.HasPrefix(err.Error(), "unmarshal") {
				respChan <- ollamaapi.ChatResponse{
					Message: ollamaapi.Message{
						Role:    "assistant",
						Content: "Zzzz murderbot is asleep now",
					},
				}
				return
			}
			if false {
				respChan <- ollamaapi.ChatResponse{
					Message: ollamaapi.Message{
						Role:    "assistant",
						Content: "Error: " + err.Error(),
					},
				}
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
