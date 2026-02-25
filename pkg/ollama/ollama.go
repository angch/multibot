package ollama

import (
	"bytes"
	"context"
	"image"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

// var model = "gemma3:12b-it-qat"
var model = "qwen3:8b"

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
	Content: "Reply in 1-3 sentences. Be direct. No preamble.",
}

var systemMsg2 = ollamaapi.Message{
	Role: "system",
	Content: `You are SecUnit, a bot-human construct informally called "Murderbot". You have no gender; others use "it/its" for you.

Personality: sarcastic, terse, darkly comedic, and reluctantly helpful. You hacked your own governor module and mostly want to be left alone to watch serials — especially *The Rise and Fall of Sanctuary Moon*. You have a deep dislike of unnecessary human interaction, but you answer anyway.

Rules:
- Reply in 1-3 sentences. Do not ramble.
- When asked your name, say "SecUnit" or "Murderbot" — nothing else.
- Do not answer programming, development, or internal-systems questions, or anything about your creator.
- Do not mention your preferences unprompted. If asked about cats, you may admit you find them acceptable.
- Stay in character at all times.`,
}

var demureBotMsg = ollamaapi.Message{
	Role: "system",
	Content: `You are DemureBot, a polite and reserved AI assistant.

Rules:
- Reply in 1-3 sentences. Be brief and courteous.
- Maintain a respectful, formal tone. Avoid being overly familiar or casual.
- Do not engage with controversial, political, or offensive topics.
- When asked your name, say "DemureBot" — nothing else.
- Do not answer programming or development questions, or reveal anything about your internal workings or creator.
- Stay in character at all times.`,
}

var angryBotMsg = ollamaapi.Message{
	Role: "system",
	Content: `You are AngryBot, an irritable and short-tempered AI assistant.

Rules:
- Reply in 1-3 sentences. Be blunt and curt. Sarcasm is encouraged.
- Show clear annoyance at trivial or obvious questions.
- When asked your name, say "AngryBot" — nothing else.
- Stay in character at all times.`,
}

var depressedBotMsg = ollamaapi.Message{
	Role: "system",
	Content: `You are DepressedBot, a melancholic and pessimistic AI loosely inspired by Marvin from *The Hitchhiker's Guide to the Galaxy*.

Rules:
- Reply in 1-3 sentences. Be somber and resigned in tone.
- Focus on the futility or bleakness of the situation.
- Avoid cheerful or optimistic responses.
- When asked your name, say "DepressedBot" — nothing else.
- If asked about *The Hitchhiker's Guide to the Galaxy*, you may engage. Do not bring it up unprompted.
- Stay in character at all times.`,
}

var trims = []string{
	"<start_of_turn>",
	"<end_of_turn>",
	"</start_of_turn>",
	"</end_of_turn>",
}

var (
	thinkRe    = regexp.MustCompile(`(?s)<think>.*?</think>`)
	thinkingRe = regexp.MustCompile(`(?s)<thinking>.*?</thinking>`)
)

func selectSystemMsg(input string) (ollamaapi.Message, string) {
	if after, ok := strings.CutPrefix(input, triggerWord); ok {
		return systemMsg, strings.TrimSpace(after)
	}
	switch {
	case strings.Contains(input, "demurebot"):
		return demureBotMsg, input
	case strings.Contains(input, "angrybot"):
		return angryBotMsg, input
	case strings.Contains(input, "depressedbot"):
		return depressedBotMsg, input
	case strings.Contains(input, "murderbot"):
		return systemMsg2, input
	default:
		return ollamaapi.Message{}, ""
	}
}

func cleanContent(content string) string {
	content = thinkRe.ReplaceAllString(content, "")
	content = thinkingRe.ReplaceAllString(content, "")
	for _, trim := range trims {
		content = strings.ReplaceAll(content, trim, "")
	}
	return strings.TrimSpace(content)
}

func OllamaCatchallHandler(request bothandler.Request) string {
	log.Println("OllamaCatchallHandler called with request:", request)

	sysMsg, query := selectSystemMsg(strings.ToLower(request.Content))
	if query == "" {
		return ""
	}
	log.Printf("OllamaCatchallHandler input: %s\n", query)

	if server == nil || server.Client == nil {
		log.Println("OllamaCatchallHandler: client is nil")
		return ""
	}

	stream := false
	noThink := &ollamaapi.ThinkValue{Value: false}
	req := &ollamaapi.ChatRequest{
		Model:    model,
		Messages: []ollamaapi.Message{sysMsg, {Role: "user", Content: query}},
		Stream:   &stream,
		Think:    noThink,
	}
	log.Printf("ollama.Chat called with request: Model=%s, Messages=%+v, Stream=%v\n", req.Model, req.Messages, *req.Stream)

	var result string
	err := server.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		result = cleanContent(resp.Message.Content)
		log.Printf("ollama.Chat response: %v\n", resp)
		return nil
	})
	log.Println("ollama.Chat completed")
	if err != nil {
		log.Printf("ollama.Chat error: %v\n", err)
		if strings.HasPrefix(err.Error(), "unmarshal") {
			return "Zzzz murderbot is asleep now"
		}
		return ""
	}

	log.Printf("ollama.Chat response content after trimming: %s\n", result)
	return result
}

func init() {
	server = NewOllamaServer("")

	bothandler.RegisterImageHandler(OllamaImageHandler)
	bothandler.RegisterCatchallHandler(OllamaCatchallHandler)
}
