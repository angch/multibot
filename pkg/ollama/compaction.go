package ollama

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	ollamaapi "github.com/ollama/ollama/api"
)

const compactionLineThreshold = 200

// Per-key mutex serialises file appends and compaction for the same history file.
var fileKeyMu sync.Map // map[string]*sync.Mutex

func keyMu(key string) *sync.Mutex {
	v, _ := fileKeyMu.LoadOrStore(key, new(sync.Mutex))
	return v.(*sync.Mutex)
}

// activityCh receives a signal on each successful LLM response.
// Buffered so rapid callers don't block.
var activityCh = make(chan struct{}, 1)

func notifyActivity() {
	select {
	case activityCh <- struct{}{}:
	default:
	}
}

const idleDuration = time.Hour

// idleLoop runs compaction once per idle period (1 h of no activity),
// then waits for the next activity signal before arming again.
func idleLoop() {
	timer := time.NewTimer(idleDuration)
	defer timer.Stop()
	for {
		select {
		case <-activityCh:
			// Activity: reset the countdown.
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(idleDuration)

		case <-timer.C:
			// Idle for a full hour — compact all eligible files.
			runCompaction()
			// Block until there is new activity before re-arming,
			// so we don't spin compacting files that haven't changed.
			<-activityCh
			timer.Reset(idleDuration)
		}
	}
}

func personalityFromKey(key string) string {
	for _, p := range []string{"murderbot", "demurebot", "angrybot", "depressedbot"} {
		if strings.HasSuffix(key, "_"+p) {
			return p
		}
	}
	return "default"
}

func runCompaction() {
	entries, err := os.ReadDir(memoryDir)
	if err != nil {
		log.Printf("memory: readdir %s: %v", memoryDir, err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		compactFile(strings.TrimSuffix(entry.Name(), ".md"))
	}
}

func compactFile(key string) {
	mu := keyMu(key)
	mu.Lock()
	defer mu.Unlock()

	path := memoryFilePath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Locate existing summary block and the start of raw log entries.
	var existingSummary string
	rawStart := 0
	startIdx := strings.Index(content, summaryStartMarker)
	endIdx := strings.Index(content, summaryEndMarker)
	if startIdx >= 0 && endIdx > startIdx {
		existingSummary = strings.TrimSpace(content[startIdx+len(summaryStartMarker) : endIdx])
		for i, line := range lines {
			if strings.Contains(line, summaryEndMarker) {
				rawStart = i + 1
				break
			}
		}
	}

	rawLines := lines[rawStart:]
	if len(rawLines) <= compactionLineThreshold {
		return
	}

	// Summarise the oldest half; keep the rest verbatim.
	half := len(rawLines) / 2
	toSummarize := strings.Join(rawLines[:half], "\n")
	toKeep := rawLines[half:]

	personality := personalityFromKey(key)
	newSummary := llmSummarize(toSummarize, personality)
	if newSummary == "" {
		log.Printf("memory: compaction of %s skipped (empty LLM summary)", key)
		return
	}

	combined := newSummary
	if existingSummary != "" {
		combined = existingSummary + "\n\n" + newSummary
	}

	var sb strings.Builder
	sb.WriteString(summaryStartMarker)
	sb.WriteString("\n")
	sb.WriteString(combined)
	sb.WriteString("\n")
	sb.WriteString(summaryEndMarker)
	sb.WriteString("\n")
	for _, line := range toKeep {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(sb.String()), 0600); err != nil {
		log.Printf("memory: write tmp %s: %v", tmpPath, err)
		return
	}
	if err := os.Rename(tmpPath, path); err != nil {
		log.Printf("memory: rename %s: %v", tmpPath, err)
		_ = os.Remove(tmpPath)
		return
	}

	// Invalidate the mtime-based summary cache so the new summary is picked up.
	summaryMu.Lock()
	delete(summaryEntries, key)
	summaryMu.Unlock()

	log.Printf("memory: compacted %s (%d raw lines -> %d kept)", key, len(rawLines), len(toKeep))
}

func llmCompactMessage(msg string) string {
	if server == nil || server.Client == nil {
		return msg[:maxMessageLen] + "… [truncated]"
	}
	prompt := "The following message is too long to process directly. Summarise it concisely, preserving the core question or intent, in 200 words or fewer:\n---\n" + msg
	stream := false
	noThink := &ollamaapi.ThinkValue{Value: true}
	req := &ollamaapi.ChatRequest{
		Model:  model,
		Stream: &stream,
		Think:  noThink,
		Messages: []ollamaapi.Message{
			{Role: "user", Content: prompt},
		},
	}
	var result string
	if err := server.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		result = cleanContent(resp.Message.Content)
		return nil
	}); err != nil {
		log.Printf("memory: llmCompactMessage error: %v, falling back to truncation", err)
		return msg[:maxMessageLen] + "… [truncated]"
	}
	if result == "" {
		return msg[:maxMessageLen] + "… [truncated]"
	}
	log.Printf("memory: compacted long message (%d chars -> %d chars)", len(msg), len(result))
	return result
}

func llmSummarize(rawLog, personality string) string {
	if server == nil || server.Client == nil {
		return ""
	}
	prompt := fmt.Sprintf(
		"Summarise the following conversation log from %s's perspective in 3-5 sentences.\n"+
			"Capture: recurring topics, user names and what they tend to ask about, any promises or conclusions reached.\n"+
			"Keep the summary in character.\n---\n%s",
		personality, rawLog,
	)
	stream := false
	noThink := &ollamaapi.ThinkValue{Value: true}
	req := &ollamaapi.ChatRequest{
		Model:  model,
		Stream: &stream,
		Think:  noThink,
		Messages: []ollamaapi.Message{
			{Role: "user", Content: prompt},
		},
	}
	var result string
	if err := server.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		result = cleanContent(resp.Message.Content)
		return nil
	}); err != nil {
		log.Printf("memory: llmSummarize error: %v", err)
		return ""
	}
	return result
}
