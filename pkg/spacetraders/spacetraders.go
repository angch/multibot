package spacetraders

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/angch/multibot/pkg/bothandler"
	"github.com/angch/multibot/pkg/engineersmy"
	"gorm.io/gorm"
)

// FIXME: This is a placeholder work in progress.

// PlatformChannel maps chat channels and platforms to Agent.Symbol
type PlatformChannel struct {
	Platform string
	Channel  string
}

var globalState = map[PlatformChannel]*Agent{}

type ChannelAgents struct {
	gorm.Model
	Platform    string
	Channel     string
	AgentSymbol string
}

type Agent struct {
	gorm.Model
	Agent     string `gorm:"uniqueIndex"`
	Faction   string
	AuthToken string
}

type AgentState struct {
	Systems map[string]System
	Ships   map[string]Ship
}

type Ship struct {
	LastUpdate time.Time
	Fuel       int
	InOrbit    int
}
type System struct {
	LastUpdate time.Time
}
type Waypoint struct {
}

var lock = sync.Mutex{}
var activeDev = true

type SpaceTraders struct {
	GormDB     *gorm.DB
	Rand       *rand.Rand
	HttpClient http.Client
}

var this = SpaceTraders{
	Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	HttpClient: http.Client{
		Timeout: time.Second * 10,
	},
}

func init() {
	if activeDev {
		log.Println("pkg/spacetraders/init")
	}
	// Singleton pattern, to fit in with the rest of the bot architecture
	bothandler.RegisterCatchallHandler(SpaceTradersHandler)
	load()
}

func isValidPlatformChannel(platform, channel string) bool {
	switch platform {
	case "discord":
		ok := engineersmy.IsKnownDiscordChannel("spacetraders", channel)
		// ok = ok || engineersmy.IsKnownDiscordChannel("sandbox", channel)
		return ok
	case "readline":
		return true
	default:
		return false
	}
}

func SpaceTradersHandler(request bothandler.Request) string {
	if activeDev {
		log.Printf("pkg/spacetraders/SpaceTradersHandler %+v\n", request)
	}

	if !isValidPlatformChannel(request.Platform, request.Channel) {
		return ""
	}

	input := request.Content

	// Instead of a command, we route off the channel
	words := removeEmptyStrings(strings.Split(input, " "))
	if len(words) < 1 {
		return ""
	}

	channelAgent := PlatformChannel{request.Platform, request.Channel}
	agentState, ok := globalState[channelAgent]
	if !ok && words[0] != "init" {
		return "This agent is not initialized"
	}

	switch words[0] {
	case "status":
		return fmt.Sprintf("%+v", agentState)
	case "init":
		if agentState != nil {
			return "This agent is already initialized as" + agentState.Agent
		}
		if len(words) < 2 {
			return "Need a callsign (faction is always COSMIC)"
		}
		req := RegisterAgentRequest{
			Symbol:  words[1],
			Faction: "COSMIC",
		}
		resp := this.RegisterAgent(req)
		if resp == nil {
			return "Failed to register agent"
		}

		data := resp.Data
		agent := data.Agent
		if agent == nil {
			log.Println(agent)
		}
		faction := data.Faction
		if faction == nil {
			log.Println(faction)
		}
		agentState = &Agent{
			Agent:     agent.Symbol,
			Faction:   faction.Symbol,
			AuthToken: data.Token,
		}

		log.Printf("%+v\n", resp)
		return fmt.Sprintf("Registering callsign %s faction %s", words[1], "COSMIC")
	case "agent":
		return "agent detaisl is work in progress"
	default:
		return ""
	}
}
