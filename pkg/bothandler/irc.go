package bothandler

import (
	"crypto/tls"
	"log"
	"os"
	"strings"

	irc "gopkg.in/irc.v3"
)

// Implements MessagePlatform
type IrcMessagePlatform struct {
	Conn           *tls.Conn
	Signal         chan os.Signal
	ClientConfig   *irc.ClientConfig
	Client         *irc.Client
	DefaultChannel string
	CloseMe        bool
}

func NewMessagePlatformFromIrc(serveraddr string, clientconfig *irc.ClientConfig, signal chan os.Signal) (*IrcMessagePlatform, error) {
	// if clientconfig == nil {
	// 	return nil, fmt.Errorf("No client config supplied")
	// }
	if clientconfig == nil {
		clientconfig = &irc.ClientConfig{}
	}
	conn, err := tls.Dial("tcp", serveraddr, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &IrcMessagePlatform{
		Conn:         conn,
		Signal:       signal,
		ClientConfig: clientconfig,
	}, nil
}

func (s *IrcMessagePlatform) ProcessMessages() {
	s.ClientConfig.Handler = irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
		// log.Printf("irchandler %+v\n", *m)
		if m.Command == "001" {
			// 001 is a welcome event, so we join channels there
			err := c.Write("JOIN #" + s.DefaultChannel)
			if err != nil {
				log.Println(err)
			}
		} else if m.Command == "PRIVMSG" {
			// log.Printf("params are: %v\n", m.Params)
			if len(m.Params) > 1 {
				channel := m.Params[0]
				content := strings.Join(m.Params[1:], " ")

				h, ok := Handlers[content]
				if ok {
					response := h()
					err := c.WriteMessage(&irc.Message{
						Command: "PRIVMSG",
						Params: []string{
							channel,
							response,
						},
					})
					if err != nil {
						log.Println(err)
					}
				}

				// Can be better to decouple 1 to 1 of message : response
				for _, v := range CatchallHandlers {
					r := v(content)
					if r != "" {
						err := c.WriteMessage(&irc.Message{
							Command: "PRIVMSG",
							Params: []string{
								channel,
								r,
							},
						})
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	})
	// Create the client
	client := irc.NewClient(s.Conn, *s.ClientConfig)
	s.Client = client

	for {
		err := client.Run()
		if err != nil {
			log.Println(err)
			if s.CloseMe {
				break
			}
			if err.Error() == "EOF" {
				continue
			}
		}
	}
}

func (s *IrcMessagePlatform) Close() {
	if s != nil && s.Conn != nil {
		s.CloseMe = true
		s.Conn.Close()
	}
}

func (s *IrcMessagePlatform) Send(text string) {
	if s == nil {
		return
	}
	err := s.ChannelMessageSend("", text)
	if err != nil {
		log.Println(err)
	}
}

func (s *IrcMessagePlatform) ChannelMessageSend(channelId, message string) error {
	if channelId == "" {
		channelId = s.DefaultChannel
	}
	err := s.Client.WriteMessage(&irc.Message{
		Command: "PRIVMSG",
		Params: []string{
			channelId,
			message,
		},
	})
	if err != nil {
		log.Println(err)
	}
	return nil
}
