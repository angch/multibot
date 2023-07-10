package spacetraders

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/angch/discordbot/pkg/bothandler"
)

// FIXME: This is a placeholder work in progress.

// ChannelAgent maps chat channels and platforms to Agent.Symbol
type ChannelAgent struct {
	Platform string
	Channel  string
}

var globalState = map[ChannelAgent]Agent{}

type Agent struct {
	Symbol    string
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
var myrand = rand.New(rand.NewSource(time.Now().UnixNano()))

var activeDev = true

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

}

func save() {
	lock.Lock()
	defer lock.Unlock()
}

func SpaceTradersHandler(request bothandler.Request) string {
	if activeDev {
		log.Printf("pkg/spacetraders/SpaceTradersHandler %+v\n", request)
	}
	input := request.Content
	// Jan 2 15:04:05 2006 MST
	today := time.Now().Local().Format("20060102")
	key := fmt.Sprintf("%s/%s/%s", request.Platform, request.Channel, today)
	lock.Lock()

	_ = input
	_ = key

	return ""
}
