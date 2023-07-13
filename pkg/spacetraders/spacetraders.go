package spacetraders

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/angch/discordbot/pkg/bothandler"
	"github.com/angch/discordbot/pkg/engineersmy"
	"github.com/glebarez/sqlite"
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
	Symbol    string `gorm:"uniqueIndex"`
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
	Db   *gorm.DB
	Rand *rand.Rand
}

var this = SpaceTraders{
	Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
}

func init() {
	if activeDev {
		log.Println("pkg/spacetraders/init")
	}
	// Singleton pattern, to fit in with the rest of the bot architecture
	bothandler.RegisterCatchallHandler(SpaceTradersHandler)
	load()
}

const savefile string = "spacetraders.sqlite"

func load() {
	lock.Lock()
	defer lock.Unlock()

	db, err := gorm.Open(sqlite.Open(savefile), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&ChannelAgents{}, &Agent{})
	this.Db = db

	// FIXME: slurp everything into globalState
	ca := make([]ChannelAgents, 0)
	db.Find(&ca)
	for _, c := range ca {
		ag := &Agent{}
		db.First(ag, "symbol = ?", c.AgentSymbol)
		globalState[PlatformChannel{c.Platform, c.Channel}] = ag
	}
}

func save() {
	lock.Lock()
	// FIXME: dump everything from globalState into spacetraders.sqlite
	defer lock.Unlock()
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

func removeEmptyStrings(words []string) []string {
	var ret []string
	for _, w := range words {
		if w != "" {
			ret = append(ret, w)
		}
	}
	return ret
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
			return "This agent is already initialized as" + agentState.Symbol
		}
		if len(words) < 3 {
			return "Need a callsign and faction"
		}
		return fmt.Sprintf("Registering callsign %s faction %s", words[1], words[2])
	case "agent":
		return "agent detaisl is work in progress"
	default:
		return ""
	}
}
