package bothandler

import (
	"crypto/tls"
	"log"
	"os"
	"strings"
	"time"

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
	serveraddr     string // in case we need to reconnect
}

func NewMessagePlatformFromIrc(serveraddr string, clientconfig *irc.ClientConfig, signal chan os.Signal) (*IrcMessagePlatform, error) {
	if clientconfig == nil {
		clientconfig = &irc.ClientConfig{}
	}
	platform := IrcMessagePlatform{
		Signal:       signal,
		ClientConfig: clientconfig,
		serveraddr:   serveraddr,
	}
	err := platform.connect()
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

func (s *IrcMessagePlatform) connect() error {
	conn, err := tls.Dial("tcp", s.serveraddr, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	s.Conn = conn
	return nil
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
					// FIXME
					r := v(Request{content, "IRC", "", ""})
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

				sliced_content := strings.SplitN(content, " ", 2)
				if len(sliced_content) > 1 {
					command := sliced_content[0]
					actual_content := sliced_content[1]

					ih, ok := MsgInputHandlers[command]
					if ok {
						response := ih(Request{actual_content, "IRC", "", ""})
						if response != "" {
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
					}
				}

			}
		}
	})

	// Create the client
	for {
		client := irc.NewClient(s.Conn, *s.ClientConfig)
		s.Client = client
		err := client.Run()
		if err != nil {
			log.Println(err)
			if s.CloseMe {
				break
			}
			x := err.Error()
			if x == "EOF" || strings.HasPrefix(x, "use of closed network connection") {
				time.Sleep(5 * time.Second)
				err := s.connect()
				if err != nil {
					// FIXME
					break
				}
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
