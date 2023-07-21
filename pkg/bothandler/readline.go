package bothandler

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

// Implements MessagePlatform
type ReadlineMessagePlatform struct {
	Instance *readline.Instance
	Signal   chan os.Signal
}

func NewMessagePlatformFromReadline(historyfile string, signal chan os.Signal) (*ReadlineMessagePlatform, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     historyfile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &ReadlineMessagePlatform{
		Instance: l,
		Signal:   signal,
	}, nil
}

func (s *ReadlineMessagePlatform) Send(text string) {
	log.Println(text)
}

func (s *ReadlineMessagePlatform) SendWithOptions(text string, options SendOptions) {
	if s == nil {
		return
	}
	s.Send(text)
}

func (s *ReadlineMessagePlatform) ProcessMessages() {
	l := s.Instance
outer:
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		content := line
		h, ok := Handlers[content]
		if ok {
			response := h()
			fmt.Println("Bot says", response)
		}

		sliced_content := strings.SplitN(content, " ", 2)
		if len(sliced_content) > 1 {
			command := sliced_content[0]
			actual_content := sliced_content[1]

			ih, ok := MsgInputHandlers[command]
			if ok {
				response := ih(Request{actual_content, "readline", "", ""})
				if response != "" {
					fmt.Println("Bot says", response)
				}
			}
		}

		// Can be better to decouple 1 to 1 of message : response
		for _, v := range CatchallHandlers {
			r := v(Request{content, "readline", "", ""})
			if r != "" {
				fmt.Println(">", r)
			}
		}

		switch {
		case line == "bye", line == "quit":
			s.Signal <- os.Interrupt
			break outer
		}
	}
}

func (s *ReadlineMessagePlatform) Close() {
	if s != nil && s.Instance != nil {
		s.Instance.Close()
	}
}

func (s *ReadlineMessagePlatform) ChannelMessageSend(channelId, message string) error {
	return nil
}
