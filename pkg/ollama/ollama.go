package ollama

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

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
var model = "qwen3.6:27b"

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

func OllamaImageListHandler(filenames []string, request bothandler.Request) string {
	lower := strings.ToLower(request.Content)
	sysMsg, query, personality := selectSystemMsg(lower, request.Content)
	if personality == "" {
		return ""
	}

	var images []ollamaapi.ImageData
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("OllamaImageListHandler: read %s: %v", filename, err)
			continue
		}
		images = append(images, ollamaapi.ImageData(data))
	}
	if len(images) == 0 {
		return ""
	}

	if server == nil || server.Client == nil {
		return ""
	}

	key := historyKey(request.Platform, request.Channel, personality)

	effectiveSysMsg := sysMsg
	if summary := readSummary(key); summary != "" {
		effectiveSysMsg.Content += "\n\n[Memory of past conversations:\n" + summary + "]"
	}

	prior := convHistory.get(key)
	if len(prior) > maxInjectMessages {
		prior = prior[len(prior)-maxInjectMessages:]
	}

	messages := make([]ollamaapi.Message, 0, 1+len(prior)+1)
	messages = append(messages, effectiveSysMsg)
	messages = append(messages, prior...)
	messages = append(messages, ollamaapi.Message{
		Role:    "user",
		Content: query,
		Images:  images,
	})

	stream := false
	req := &ollamaapi.ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   &stream,
	}
	log.Printf("OllamaImageListHandler: calling Chat with %d image(s), query=%q", len(images), query)

	var result string
	err := server.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		result = cleanContent(resp.Message.Content)
		return nil
	})
	if err != nil {
		log.Printf("OllamaImageListHandler: Chat error: %v", err)
		return ""
	}

	if result != "" {
		convHistory.add(key,
			ollamaapi.Message{Role: "user", Content: query},
			ollamaapi.Message{Role: "assistant", Content: result},
		)
		go persistExchange(key, request.From, query, result, "")
		notifyActivity()
	}

	return result
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

// selectSystemMsg maps input to the appropriate system message and personality name.
// Returns (systemMsg, userQuery, personalityName). Returns empty personality when unrecognised.
func selectSystemMsg(lowerInput, originalInput string) (ollamaapi.Message, string, string) {
	if _, ok := strings.CutPrefix(lowerInput, triggerWord); ok {
		originalAfter := originalInput[len(triggerWord):]
		return systemMsg, strings.TrimSpace(originalAfter), "default"
	}
	switch {
	case strings.Contains(lowerInput, "demurebot"):
		return demureBotMsg, originalInput, "demurebot"
	case strings.Contains(lowerInput, "angrybot"):
		return angryBotMsg, originalInput, "angrybot"
	case strings.Contains(lowerInput, "depressedbot"):
		return depressedBotMsg, originalInput, "depressedbot"
	case strings.Contains(lowerInput, "murderbot"):
		return systemMsg2, originalInput, "murderbot"
	default:
		return ollamaapi.Message{}, "", ""
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

// ---------------------------------------------------------------------------
// Memory system (MEM-01 through MEM-08)
// ---------------------------------------------------------------------------

const (
	maxHistoryMessages = 40   // ring buffer cap (20 exchanges)
	maxInjectMessages  = 20   // messages injected into context (10 exchanges)
	maxMessageLen      = 2048 // MEM-08: per-message truncation limit
	memoryDir          = "data/memory"
	summaryStartMarker = "<!-- SUMMARY_START -->"
	summaryEndMarker   = "<!-- SUMMARY_END -->"
)

// MEM-01: per-key in-memory conversation ring buffer

type memHistory struct {
	mu      sync.RWMutex
	buffers map[string][]ollamaapi.Message
}

var convHistory = &memHistory{buffers: make(map[string][]ollamaapi.Message)}

func historyKey(platform, channel, personality string) string {
	key := platform + "_" + channel + "_" + personality
	r := strings.NewReplacer("/", "_", ":", "_", " ", "_", ".", "_")
	return r.Replace(key)
}

func (h *memHistory) add(key string, msgs ...ollamaapi.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	buf := append(h.buffers[key], msgs...)
	if len(buf) > maxHistoryMessages {
		buf = buf[len(buf)-maxHistoryMessages:]
	}
	h.buffers[key] = buf
}

func (h *memHistory) get(key string) []ollamaapi.Message {
	h.mu.RLock()
	defer h.mu.RUnlock()
	src := h.buffers[key]
	if len(src) == 0 {
		return nil
	}
	out := make([]ollamaapi.Message, len(src))
	copy(out, src)
	return out
}

func (h *memHistory) clear(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.buffers, key)
}

// MEM-02: append exchange to disk

func memoryFilePath(key string) string {
	return fmt.Sprintf("%s/%s.md", memoryDir, key)
}

func persistExchange(key, from, userMsg, botMsg, thinking string) {
	mu := keyMu(key)
	mu.Lock()
	defer mu.Unlock()

	path := memoryFilePath(key)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("memory: open %s: %v", path, err)
		return
	}
	defer f.Close()
	ts := time.Now().UTC().Format(time.RFC3339)
	if thinking != "" {
		fmt.Fprintf(f, "\n## %s — %s\n**User:** %s\n**Thinking:** %s\n**Bot:** %s\n", ts, from, userMsg, thinking, botMsg)
	} else {
		fmt.Fprintf(f, "\n## %s — %s\n**User:** %s\n**Bot:** %s\n", ts, from, userMsg, botMsg)
	}
}

// MEM-05: read compacted summary from disk with mtime-based cache

type summaryEntry struct {
	mtime   time.Time
	summary string
}

var (
	summaryMu      sync.RWMutex
	summaryEntries = map[string]*summaryEntry{}
)

func readSummary(key string) string {
	path := memoryFilePath(key)
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	mtime := info.ModTime()

	summaryMu.RLock()
	cached := summaryEntries[key]
	summaryMu.RUnlock()
	if cached != nil && cached.mtime.Equal(mtime) {
		return cached.summary
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	content := string(data)
	start := strings.Index(content, summaryStartMarker)
	end := strings.Index(content, summaryEndMarker)
	var summary string
	if start >= 0 && end > start {
		summary = strings.TrimSpace(content[start+len(summaryStartMarker) : end])
	}

	summaryMu.Lock()
	summaryEntries[key] = &summaryEntry{mtime: mtime, summary: summary}
	summaryMu.Unlock()
	return summary
}

// MEM-06: forget confirmation messages stay in character

func forgetConfirmation(personality string) string {
	switch personality {
	case "murderbot":
		return "Fine. Wiped. Like it never happened."
	case "demurebot":
		return "Of course. Your conversation history has been cleared."
	case "angrybot":
		return "FINE. Gone. Happy now?!"
	case "depressedbot":
		return "Erased. Not that it matters. Nothing does."
	default:
		return "Memory cleared."
	}
}

func isForgetCommand(query, personality string) bool {
	q := strings.TrimSpace(query)
	if q == "forget" {
		return true
	}
	// named-bot variant: "murderbot forget" -> strip personality prefix
	stripped := strings.TrimSpace(strings.TrimPrefix(q, personality))
	return stripped == "forget"
}

// ---------------------------------------------------------------------------

func OllamaCatchallHandler(request bothandler.Request) string {
	log.Println("OllamaCatchallHandler called with request:", request)

	lower := strings.ToLower(request.Content)
	sysMsg, query, personality := selectSystemMsg(lower, request.Content)
	if query == "" || personality == "" {
		return ""
	}

	// MEM-08: compact oversized messages via LLM before processing or storing
	if len(query) > maxMessageLen {
		log.Printf("memory: compacting long message from %s (%d chars)", request.From, len(query))
		query = llmCompactMessage(query)
	}

	// MEM-06: forget command — clear buffer and delete file
	if isForgetCommand(strings.ToLower(query), personality) {
		key := historyKey(request.Platform, request.Channel, personality)
		convHistory.clear(key)
		_ = os.Remove(memoryFilePath(key))
		summaryMu.Lock()
		delete(summaryEntries, key)
		summaryMu.Unlock()
		return forgetConfirmation(personality)
	}

	log.Printf("OllamaCatchallHandler input: %s\n", query)

	if server == nil || server.Client == nil {
		log.Println("OllamaCatchallHandler: client is nil")
		return ""
	}

	key := historyKey(request.Platform, request.Channel, personality)

	// MEM-05: augment system prompt with compacted summary if one exists
	effectiveSysMsg := sysMsg
	if summary := readSummary(key); summary != "" {
		effectiveSysMsg.Content += "\n\n[Memory of past conversations:\n" + summary + "]"
	}

	// MEM-03: prepend recent history between system prompt and new user message
	prior := convHistory.get(key)
	if len(prior) > maxInjectMessages {
		prior = prior[len(prior)-maxInjectMessages:]
	}
	messages := make([]ollamaapi.Message, 0, 1+len(prior)+1)
	messages = append(messages, effectiveSysMsg)
	messages = append(messages, prior...)
	messages = append(messages, ollamaapi.Message{Role: "user", Content: query})

	stream := false
	enableThink := &ollamaapi.ThinkValue{Value: true}
	req := &ollamaapi.ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   &stream,
		Think:    enableThink,
	}
	log.Printf("ollama.Chat called with request: Model=%s, Messages=%+v, Stream=%v\n", req.Model, req.Messages, *req.Stream)

	var result, thinking string
	err := server.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		result = cleanContent(resp.Message.Content)
		thinking = resp.Message.Thinking
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

	if result != "" {
		// MEM-01: append to in-memory ring buffer
		convHistory.add(key,
			ollamaapi.Message{Role: "user", Content: query},
			ollamaapi.Message{Role: "assistant", Content: result},
		)
		// MEM-02: persist exchange to disk asynchronously (thinking logged, not shown in chat)
		go persistExchange(key, request.From, query, result, thinking)
		// MEM-04: signal activity so the idle timer resets
		notifyActivity()
	}

	return result
}

func init() {
	if err := os.MkdirAll(memoryDir, 0700); err != nil {
		log.Printf("memory: mkdir %s: %v", memoryDir, err)
	}

	server = NewOllamaServer("")

	go idleLoop()

	bothandler.RegisterImageListHandler(OllamaImageListHandler)
	bothandler.RegisterCatchallHandler(OllamaCatchallHandler)
}
