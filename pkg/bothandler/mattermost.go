package bothandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mattermost/mattermost/server/public/model"
)

type MattermostWebSocketEvent struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data"`
}

// Implements MessagePlatform
type MattermostMessagePlatform struct {
	Client         *model.Client4
	WebSocketConn  *websocket.Conn
	User           *model.User
	ChannelId      map[string]string
	KnownUsers     map[string]*model.User
	KnownUsersLock sync.RWMutex
	DefaultChannel string
	TeamId         string
	BotToken       string
	ServerURL      string
	stopChan       chan bool
}

func NewMessagePlatformFromMattermost(mattermostBotToken, mattermostURL string) (*MattermostMessagePlatform, error) {
	client := model.NewAPIv4Client(mattermostURL)
	client.SetToken(mattermostBotToken)

	ctx := context.Background()

	// Test the connection and get bot user info
	user, _, err := client.GetMe(ctx, "")
	if err != nil {
		log.Printf("Failed to get bot user info: %v", err)
		return nil, err
	}

	log.Printf("Connected to Mattermost on account %s", user.Username)

	// Get teams
	teams, _, err := client.GetTeamsForUser(ctx, user.Id, "")
	if err != nil {
		log.Printf("Failed to get teams: %v", err)
		return nil, err
	}

	var teamId string
	if len(teams) > 0 {
		teamId = teams[0].Id
		log.Printf("Using team: %s", teams[0].Name)
	}

	return &MattermostMessagePlatform{
		Client:         client,
		User:           user,
		ChannelId:      map[string]string{},
		KnownUsers:     map[string]*model.User{},
		DefaultChannel: "",
		TeamId:         teamId,
		BotToken:       mattermostBotToken,
		ServerURL:      mattermostURL,
		stopChan:       make(chan bool),
	}, nil
}

func (s *MattermostMessagePlatform) ProcessMessages() {
	// Connect to WebSocket for real-time messaging
	wsURL := strings.Replace(s.ServerURL, "http", "ws", 1) + "/api/v4/websocket"

	dialer := websocket.Dialer{}
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+s.BotToken)

	conn, _, err := dialer.Dial(wsURL, headers)
	if err != nil {
		log.Printf("Failed to connect to WebSocket: %v", err)
		return
	}
	s.WebSocketConn = conn
	defer conn.Close()

	// Send authentication message
	authMsg := map[string]any{
		"seq":    1,
		"action": "authentication_challenge",
		"data": map[string]string{
			"token": s.BotToken,
		},
	}
	if err := conn.WriteJSON(authMsg); err != nil {
		log.Printf("Failed to send auth message: %v", err)
		return
	}

	for {
		select {
		case <-s.stopChan:
			return
		default:
			var event MattermostWebSocketEvent
			err := conn.ReadJSON(&event)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("WebSocket connection closed normally")
					return
				}
				log.Printf("Failed to read WebSocket message: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			s.handleWebSocketEvent(&event)
		}
	}
}

func (s *MattermostMessagePlatform) handleWebSocketEvent(event *MattermostWebSocketEvent) {
	if event.Event != "posted" {
		return
	}

	postData, ok := event.Data["post"].(string)
	if !ok {
		return
	}

	var post model.Post
	err := json.Unmarshal([]byte(postData), &post)
	if err != nil {
		log.Printf("Failed to unmarshal post: %v", err)
		return
	}

	// Ignore messages from the bot itself
	if post.UserId == s.User.Id {
		return
	}

	content := post.Message

	// Update known users - simplified, we'll skip this for now to avoid extra API calls

	// Handle exact match handlers
	h, ok := Handlers[content]
	if ok {
		response := h()
		s.sendReply(post.ChannelId, response, post.RootId)
	}

	// Handle catchall handlers
	for _, v := range CatchallHandlers {
		r := v(Request{content, "mattermost", post.ChannelId, post.UserId})
		if r != "" {
			s.sendReply(post.ChannelId, r, post.RootId)
		}
	}

	// Handle extended catchall handlers
	for _, v := range CatchallExtendedHandlers {
		r := v(ExtendedMessage{Text: content})
		if r != nil {
			if r.Text != "" && r.Image == nil {
				s.sendReply(post.ChannelId, r.Text, post.RootId)
			}
			if r.Image != nil {
				// Handle image uploads
				s.sendImageReply(post.ChannelId, r.Text, r.Image, post.RootId, content)
			}
		}
	}

	// Handle input commands (messages starting with a command)
	sliced_content := strings.SplitN(content, " ", 2)
	if len(sliced_content) > 1 {
		command := sliced_content[0]
		actual_content := sliced_content[1]

		ih, ok := MsgInputHandlers[command]
		if ok {
			response := ih(Request{actual_content, "mattermost", post.ChannelId, post.UserId})
			if response != "" {
				s.sendReply(post.ChannelId, response, post.RootId)
			}
		}
	}

	// Handle file attachments
	if len(post.FileIds) > 0 {
		for _, fileId := range post.FileIds {
			filename := fmt.Sprintf("tmp/%s", fileId)
			err := s.downloadFile(fileId, filename)
			if err != nil {
				log.Printf("Failed to download file: %v", err)
				continue
			}

			for _, v := range ImageHandlers {
				r := v(filename, Request{content, "mattermost", post.ChannelId, post.UserId})
				if r != "" {
					s.sendReply(post.ChannelId, r, post.RootId)
				}
			}

			// Optionally clean up the file
			if false {
				os.Remove(filename)
			}
		}
	}
}

func (s *MattermostMessagePlatform) sendReply(channelId, message, rootId string) {
	post := &model.Post{
		ChannelId: channelId,
		Message:   message,
	}
	if rootId != "" {
		post.RootId = rootId
	}

	ctx := context.Background()
	_, _, err := s.Client.CreatePost(ctx, post)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (s *MattermostMessagePlatform) sendImageReply(channelId, message string, imageData []byte, rootId, originalContent string) {
	// For simplicity, we'll just send the text message for now
	// Implementing file upload would require multipart form data
	s.sendReply(channelId, message, rootId)
}

func (s *MattermostMessagePlatform) downloadFile(fileId, localFilename string) error {
	ctx := context.Background()
	data, _, err := s.Client.GetFile(ctx, fileId)
	if err != nil {
		return err
	}

	out, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(data)
	return err
}

func (s *MattermostMessagePlatform) Send(text string) {
	if s == nil {
		return
	}
	s.SendWithOptions(text, SendOptions{})
}

func (s *MattermostMessagePlatform) SendWithOptions(text string, options SendOptions) {
	if s == nil {
		return
	}
	err := s.ChannelMessageSend("", text)
	if err != nil {
		log.Println(err)
	}
}

func (s *MattermostMessagePlatform) Close() {
	if s.stopChan != nil {
		close(s.stopChan)
	}
	if s.WebSocketConn != nil {
		s.WebSocketConn.Close()
	}
}

func (s *MattermostMessagePlatform) ChannelMessageSend(channel, message string) error {
	if channel == "" {
		channel = s.DefaultChannel
	}

	// If channel is a name, try to resolve it to an ID
	channelId := channel
	if channelId == "" {
		return fmt.Errorf("no channel specified")
	}

	// For now, assume channel is already an ID or we can use it directly
	// In a full implementation, you might want to cache channel name->ID mappings

	post := &model.Post{
		ChannelId: channelId,
		Message:   message,
	}

	ctx := context.Background()
	_, _, err := s.Client.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to send message to channel %s: %v", channel, err)
	}

	return nil
}
