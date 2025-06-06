package bothandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Mattermost API structures - simplified versions
type MattermostUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type MattermostTeam struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MattermostChannel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MattermostPost struct {
	Id        string   `json:"id"`
	ChannelId string   `json:"channel_id"`
	UserId    string   `json:"user_id"`
	Message   string   `json:"message"`
	FileIds   []string `json:"file_ids"`
	RootId    string   `json:"root_id"`
}

type MattermostWebSocketEvent struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// Implements MessagePlatform
type MattermostMessagePlatform struct {
	Client         *http.Client
	WebSocketConn  *websocket.Conn
	User           *MattermostUser
	ChannelId      map[string]string
	KnownUsers     map[string]*MattermostUser
	KnownUsersLock sync.RWMutex
	DefaultChannel string
	TeamId         string
	BotToken       string
	ServerURL      string
	stopChan       chan bool
}

func NewMessagePlatformFromMattermost(mattermostBotToken, mattermostURL string) (*MattermostMessagePlatform, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Test the connection and get bot user info
	user, err := getMattermostUser(client, mattermostURL, mattermostBotToken)
	if err != nil {
		log.Printf("Failed to get bot user info: %v", err)
		return nil, err
	}

	log.Printf("Connected to Mattermost on account %s", user.Username)

	// Get teams
	teams, err := getMattermostTeams(client, mattermostURL, mattermostBotToken, user.Id)
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
		KnownUsers:     map[string]*MattermostUser{},
		DefaultChannel: "",
		TeamId:         teamId,
		BotToken:       mattermostBotToken,
		ServerURL:      mattermostURL,
		stopChan:       make(chan bool),
	}, nil
}

func getMattermostUser(client *http.Client, serverURL, token string) (*MattermostUser, error) {
	req, err := http.NewRequest("GET", serverURL+"/api/v4/users/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user MattermostUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getMattermostTeams(client *http.Client, serverURL, token, userId string) ([]MattermostTeam, error) {
	req, err := http.NewRequest("GET", serverURL+"/api/v4/users/"+userId+"/teams", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get teams with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var teams []MattermostTeam
	err = json.Unmarshal(body, &teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
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
	authMsg := map[string]interface{}{
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

	var post MattermostPost
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
		s.sendReply(post.ChannelId, response, post.Id)
	}

	// Handle catchall handlers
	for _, v := range CatchallHandlers {
		r := v(Request{content, "mattermost", post.ChannelId, post.UserId})
		if r != "" {
			s.sendReply(post.ChannelId, r, post.Id)
		}
	}

	// Handle extended catchall handlers
	for _, v := range CatchallExtendedHandlers {
		r := v(ExtendedMessage{Text: content})
		if r != nil {
			if r.Text != "" && r.Image == nil {
				s.sendReply(post.ChannelId, r.Text, post.Id)
			}
			if r.Image != nil {
				// Handle image uploads
				s.sendImageReply(post.ChannelId, r.Text, r.Image, post.Id, content)
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
				s.sendReply(post.ChannelId, response, post.Id)
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
					s.sendReply(post.ChannelId, r, post.Id)
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
	postData := map[string]interface{}{
		"channel_id": channelId,
		"message":    message,
	}
	if rootId != "" {
		postData["root_id"] = rootId
	}

	jsonData, err := json.Marshal(postData)
	if err != nil {
		log.Printf("Failed to marshal post data: %v", err)
		return
	}

	req, err := http.NewRequest("POST", s.ServerURL+"/api/v4/posts", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+s.BotToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to send message, status %d: %s", resp.StatusCode, string(body))
	}
}

func (s *MattermostMessagePlatform) sendImageReply(channelId, message string, imageData []byte, rootId, originalContent string) {
	// For simplicity, we'll just send the text message for now
	// Implementing file upload would require multipart form data
	s.sendReply(channelId, message, rootId)
}

func (s *MattermostMessagePlatform) downloadFile(fileId, localFilename string) error {
	req, err := http.NewRequest("GET", s.ServerURL+"/api/v4/files/"+fileId, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.BotToken)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download file, status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
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

	postData := map[string]interface{}{
		"channel_id": channelId,
		"message":    message,
	}

	jsonData, err := json.Marshal(postData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.ServerURL+"/api/v4/posts", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.BotToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message to channel %s, status %d: %s", channel, resp.StatusCode, string(body))
	}

	return nil
}
