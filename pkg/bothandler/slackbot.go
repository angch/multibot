package bothandler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// Implements MessagePlatform
// New way, using Slack Socket Mode
type SlackMessagePlatform struct {
	Client           *slack.Client
	SocketModeClient *socketmode.Client
	ChannelId        map[string]string
	AuthResponse     *slack.AuthTestResponse
	DefaultChannel   string
}

func NewMessagePlatformFromSlack(slackbottoken, slackapptoken string) (*SlackMessagePlatform, error) {
	if !strings.HasPrefix(slackapptoken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}
	if !strings.HasPrefix(slackbottoken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	client := slack.New(
		slackbottoken,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(slackapptoken),
	)
	authresp, err := client.AuthTest()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("AuthResponse: %+v\n", authresp)

	socketmodeclient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	params := slack.GetConversationsParameters{}
	conversations, next, err := client.GetConversations(&params)
	if err != nil {
		log.Println("Can't get conversation")
		return nil, err
	}
	// log.Println(conversations)
	channelid := make(map[string]string)
	for _, v := range conversations {
		// log.Println(v.ID, "is", v.IsChannel, v.Name)
		if v.IsChannel {
			channelid[v.Name] = v.ID
		}
	}
	_ = next
	// log.Println(next)

	return &SlackMessagePlatform{
		Client:           client,
		SocketModeClient: socketmodeclient,
		ChannelId:        channelid,
		AuthResponse:     authresp,
	}, nil
}

func (s *SlackMessagePlatform) ProcessMessages() {
	client := s.SocketModeClient
	go func() {
	eventloop:
		for evt := range client.Events {
			// log.Printf("Ping events %+v\n", evt)
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				log.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				log.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Ignored %+v\n", evt)

					continue eventloop
				}

				// log.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						// _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
						// if err != nil {
						// 	log.Printf("failed posting message: %v", err)
						// }
					case *slackevents.MemberJoinedChannelEvent:
						log.Printf("user %q joined to channel %q", ev.User, ev.Channel)
					case *slackevents.MessageEvent:
						if ev.BotID == s.AuthResponse.BotID {
							continue eventloop
						}
						// log.Println("xxx", ev.Text)

						content := ev.Text

						h, ok := Handlers[content]
						if ok {
							response := h()
							_, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(response, false))
							if err != nil {
								log.Println(err)
							}
						}

						// Can be better to decouple 1 to 1 of message : response
						for _, v := range CatchallHandlers {
							// FIXME
							r := v(Request{content, "slack", "", ""})
							if r != "" {
								_, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(r, false))
								if err != nil {
									log.Println(err)
								}
							}
						}
					default:
						log.Printf("Inner event %+v %T\n", ev, ev)
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					log.Printf("Ignored %+v\n", evt)
					continue
				}

				// log.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					// See https://api.slack.com/apis/connections/socket-implement#button

					client.Debugf("button clicked!")
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				client.Ack(*evt.Request, payload)
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					log.Printf("Ignored %+v\n", evt)
					continue
				}

				// client.Debugf("Slash command received: %+v", cmd)
				_ = cmd

				payload := map[string]interface{}{
					"blocks": []slack.Block{
						slack.NewSectionBlock(
							&slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: "foo",
							},
							nil,
							slack.NewAccessory(
								slack.NewButtonBlockElement(
									"",
									"somevalue",
									&slack.TextBlockObject{
										Type: slack.PlainTextType,
										Text: "bar",
									},
								),
							),
						),
					}}

				client.Ack(*evt.Request, payload)
			case socketmode.EventTypeHello:
				// Ignore me
			default:
				log.Printf("Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()
	err := client.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SlackMessagePlatform) Send(text string) {
	if s == nil || s.SocketModeClient == nil {
		return
	}
	err := s.ChannelMessageSend("", text)
	if err != nil {
		log.Println(err)
	}
}

func (s *SlackMessagePlatform) Close() {
}

func (s *SlackMessagePlatform) ChannelMessageSend(channel, message string) error {
	if channel == "" {
		channel = s.DefaultChannel
	}
	channelId, ok := s.ChannelId[channel]
	if !ok {
		log.Println("Unknown channel", channel)
		return fmt.Errorf("Unknown channel %s", channel)
	}
	log.Println("sending", message, "to", channelId)
	// m := sSock.NewOutgoingMessage(message, channelId)
	// s.Rtm.SendMessage(m)
	msg := slack.MsgOptionText(message, true)
	_, _, err := s.Client.PostMessage(channelId, msg)

	return err
}
