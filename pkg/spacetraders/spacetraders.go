package spacetraders

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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

// type AgentState struct {
// 	Systems map[string]System
// 	Ships   map[string]Ship
// }

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

var activeDev = true

type SpaceTraders struct {
	GormDB     *gorm.DB
	Rand       *rand.Rand
	HttpClient http.Client

	ActiveDev bool

	lock sync.RWMutex

	KnownFactions map[string]Faction
	KnownTraits   map[string]Trait
	KnownAgents   map[string]Agent
}

var this = SpaceTraders{
	Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	HttpClient: http.Client{
		Timeout: time.Second * 10,
	},

	KnownFactions: map[string]Faction{},
	KnownTraits:   map[string]Trait{},
	KnownAgents:   map[string]Agent{},
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

// SpaceTradersHandler is the catchall handler for SpaceTraders, to make it work within the bothandler framework
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
	ctx := context.Background()

	switch words[0] {
	case "status":
		return fmt.Sprintf("This channel's agent is called %+v", agentState.Agent)
		// return fmt.Sprintf("%+v", agentState)
	case "init":
		if agentState != nil {
			return "This agent is already initialized as " + agentState.Agent
		}
		if len(words) < 2 {
			return "Need a callsign (faction is always COSMIC)"
		}
		req := RegisterAgentRequest{
			Symbol:  words[1],
			Faction: "COSMIC",
		}
		pc := PlatformChannel{
			Platform: request.Platform,
			Channel:  request.Channel,
		}

		// We split the api calls and the code to handle the response,
		// in case for debugging we want to playback the response
		// to fix things.
		resp, err := this.RegisterAgent(ctx, pc, req)
		if err != nil {
			return fmt.Sprintf("Failed to register agent: %s", err)
		}
		if resp == nil {
			return "Failed to register agent"
		}

		this.ProcessRegisterAgentResponse(ctx, pc, resp)

		// log.Printf("%+v\n", resp)
		return fmt.Sprintf("Registering callsign %s faction %s", words[1], "COSMIC")
	case "agent":
		// agentCode := agentState.Agent
		// if len(words) < 2 {
		// 	agentCode = strings.ToUpper(words[1])
		// }
		return agentState.Agent + " is in the faction " + agentState.Faction

		// return "agent details is work in progress"
	case "faction":
		if len(words) < 2 {
			return "Need id for faction"
		}
		factionCode := strings.ToUpper(words[1])

		faction, ok := this.KnownFactions[factionCode]
		if !ok {
			return "No such faction " + factionCode
		}
		return faction.PrettyPrint()
	case "replay":
		if request.Platform == "readline" {
			if len(words) < 2 {
				return "Need id for replay"
			}
			gormdb := this.GormDB
			requestLog := RequestLog{}
			arg, err := strconv.Atoi(words[1])
			if err != nil {
				return "Need id for replay"
			}
			err = gormdb.Where("id = ?", arg).Find(&requestLog).Error
			if err != nil || requestLog.ID == 0 {
				return fmt.Sprintf("Failed to find request log: %d", arg)
			}

			// FIXME Yes, refac this.
			switch requestLog.Type {
			case "RegisterAgent":
				registerAgentResponse := RegisterAgentResponse{}
				json.Unmarshal([]byte(requestLog.Response), &registerAgentResponse)
				if registerAgentResponse.Error != nil {
					return fmt.Sprintf("Replaying RegisterAgent %d: error: %+v", arg, registerAgentResponse.Error)
				}
				this.ProcessRegisterAgentResponse(ctx, PlatformChannel{requestLog.Platform, requestLog.Channel}, &registerAgentResponse)
				return fmt.Sprintf("Replaying RegisterAgent %d: %+v", arg, registerAgentResponse.Data)
			default:
				log.Println("Unknown type", requestLog.Type)
			}

			return fmt.Sprintf("Replaying %d: %+v", arg, requestLog)
		} else {
			return "Replay is only supported on readline"
		}
	default:
		return ""
	}
}
