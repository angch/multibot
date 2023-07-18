package spacetraders

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var defaultGormConfig = gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
	// Logger: &GormLogger{},
}

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

type Trait struct {
	gorm.Model
	Trait       string `gorm:"uniqueIndex"`
	Description string
}

type Faction struct {
	gorm.Model
	Faction      string `gorm:"uniqueIndex"`
	Name         string
	Description  string
	Headquarters string
	Traits       []Trait `gorm:"-"`
	TraitsString string  `gorm:"column:traits;type:text"`
	IsRecruiting bool    `gorm:"column:is_recruiting"`
}

type Ship struct {
	gorm.Model
	Ship  string `gorm:"uniqueIndex" json:"symbol"`
	Owner string `gorm:"column:owner" json:"owner"`

	Data string `gorm:"column:data;type:text" json:"-"`

	Nav          NavData              `gorm:"-" json:"nav"`
	Crew         CrewData             `gorm:"-" json:"crew"`
	Fuel         FuelData             `gorm:"-" json:"fuel"`
	Frame        ShipFrameData        `gorm:"-" json:"frame"`
	Reactor      ShipReactorData      `gorm:"-" json:"reactor"`
	Engine       ShipEngineData       `gorm:"-" json:"engine"`
	Modules      []ShipModuleData     `gorm:"-" json:"modules"`
	Mounts       []ShipMountData      `gorm:"-" json:"mounts"`
	Registration ShipRegistrationData `gorm:"-" json:"registration"`
	Cargo        ShipCargoData        `gorm:"-" json:"cargo"`
}

func (s *Ship) Fix(a *SpaceTraders) {
	if s.Data != "" {
		d := s.Data
		o := s.Owner
		err := json.Unmarshal([]byte(d), s)
		if err != nil {
			log.Println(err)
		}
		s.Data = d
		s.Owner = o
	} else {
		d, err := json.Marshal(s)
		if err != nil {
			log.Println(err)
		}
		s.Data = string(d)
	}
}

func (f *Faction) Fix(a *SpaceTraders) {
	f.Traits = make([]Trait, 0)
	for _, v := range strings.Split(f.TraitsString, ",") {
		t := a.GetTrait(v)
		f.Traits = append(f.Traits, t)
	}
}

func (a *SpaceTraders) GetTrait(traitCode string) Trait {
	a.lock.RLock()
	defer a.lock.RUnlock()
	t, ok := a.KnownTraits[traitCode]
	if !ok {
		return Trait{Trait: traitCode}
	}
	return t
}

func (f *Faction) PrettyPrint() string {
	a := fmt.Sprintf("%s (%s): %s\n", f.Name, f.Faction, f.Description)
	for _, v := range f.Traits {
		a += "* *" + v.Trait + "* : " + v.Description + "\n"
	}
	if f.IsRecruiting {
		a += "\nThis faction is recruiting.\n"
	}
	return a
}

func (s *Ship) PrettyPrint() string {
	// a := fmt.Sprintf("%s\n  Nav: %+v\n Crew: %+v\n Fuel: %+v\n Frame: %+v\n Reactor: %+v\n Engine: %+v\n Modules: %+v\n Mounts: %+v\n Registration: %+v\n Cargo: %+v\n",
	// s.Ship, s.Nav, s.Crew, s.Fuel, s.Frame, s.Reactor, s.Engine, s.Modules, s.Mounts, s.Registration, s.Cargo)
	a := fmt.Sprintf("%s\n  Nav: %+v\n Crew: %+v\n Fuel: %+v\n Registration: %+v\n Cargo: %+v\n",
		s.Ship, s.Nav, s.Crew, s.Fuel, s.Registration, s.Cargo)
	return a
}

const savefile string = "spacetraders.sqlite"

var lock sync.Mutex

func load() {
	lock.Lock()

	gormdb, err := gorm.Open(sqlite.Open(savefile), &defaultGormConfig)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	gormdb.AutoMigrate(
		&ChannelAgents{},
		&Agent{},
		&RequestLog{},
		&Trait{},
		&Faction{},
		&Ship{},
	)
	this.GormDB = gormdb

	// FIXME: slurp everything into globalState
	ca := make([]ChannelAgents, 0)
	gormdb.Find(&ca)
	for _, c := range ca {
		ag := &Agent{}
		gormdb.First(ag, "agent = ?", c.AgentSymbol)
		globalState[PlatformChannel{c.Platform, c.Channel}] = ag
	}
	lock.Unlock()

	this.LoadKnown()
}

func save() {
	lock.Lock()
	// FIXME: dump everything from globalState into spacetraders.sqlite
	defer lock.Unlock()
}

func (a *SpaceTraders) SetAgent(agent Agent) {
	// gormdb := a.GormDB
}

func (a *SpaceTraders) SetTrait(trait Trait) {
	gormdb := a.GormDB

	a.lock.RLock()
	knownTrait, ok := a.KnownTraits[trait.Trait]
	a.lock.RUnlock()
	if !ok {
		knownTrait = trait
	} else {
		knownTrait.Trait = trait.Trait
		knownTrait.Description = trait.Description
	}

	err := gormdb.Save(&knownTrait).Error
	if err != nil {
		log.Println(err)
	}
	a.lock.Lock()
	a.KnownTraits[trait.Trait] = knownTrait
	a.lock.Unlock()
}

func (a *SpaceTraders) SetFaction(faction Faction) {
	gormdb := a.GormDB

	a.lock.RLock()
	knownFaction, ok := a.KnownFactions[faction.Faction]
	a.lock.RUnlock()
	if !ok {
		knownFaction = faction
	} else {
		knownFaction.Name = faction.Name
		knownFaction.Description = faction.Description
		knownFaction.Headquarters = faction.Headquarters
		knownFaction.IsRecruiting = faction.IsRecruiting
		knownFaction.TraitsString = faction.TraitsString
	}
	err := gormdb.Save(&knownFaction).Error
	if err != nil {
		log.Println(err)
	}
	a.lock.Lock()
	a.KnownFactions[faction.Faction] = knownFaction
	a.lock.Unlock()
}

func (a *SpaceTraders) SetShip(ship Ship) {
	gormdb := a.GormDB
	ship.Fix(a)
	log.Printf("SetShip owner %s\n", ship.Owner)

	a.lock.RLock()
	knownShip, ok := a.KnownShips[ship.Ship]
	a.lock.RUnlock()
	if !ok {
		knownShip = ship
	} else {
		knownShip.Crew = ship.Crew
		knownShip.Fuel = ship.Fuel
		knownShip.Frame = ship.Frame
		knownShip.Reactor = ship.Reactor
		knownShip.Engine = ship.Engine
		knownShip.Modules = ship.Modules
		knownShip.Mounts = ship.Mounts
		knownShip.Registration = ship.Registration
		knownShip.Cargo = ship.Cargo
		knownShip.Data = ship.Data
		knownShip.Owner = ship.Owner
	}
	err := gormdb.Debug().Save(&knownShip).Error
	if err != nil {
		log.Println(err)
	}
	a.lock.Lock()
	a.KnownShips[ship.Ship] = knownShip
	ships, ok := a.AgentShips[ship.Owner]
	if !ok {
		a.AgentShips[ship.Owner] = make(map[string]bool)
		ships = a.AgentShips[ship.Owner]
	}
	ships[ship.Ship] = true
	a.lock.Unlock()
}

func (a *SpaceTraders) LoadKnown() {
	traits := make([]Trait, 0)
	a.GormDB.Find(&traits) // FIXME
	for _, v := range traits {
		a.lock.Lock()
		a.KnownTraits[v.Trait] = v
		a.lock.Unlock()
	}

	factions := make([]Faction, 0)
	a.GormDB.Find(&factions) // FIXME

	for _, v := range factions {
		v.Fix(a)
		a.lock.Lock()
		a.KnownFactions[v.Faction] = v
		a.lock.Unlock()
	}

	ships := make([]Ship, 0)
	a.GormDB.Find(&ships)
	for _, v := range ships {
		v.Fix(a)

		a.lock.Lock()
		a.KnownShips[v.Ship] = v
		ships, ok := a.AgentShips[v.Owner]
		if !ok {
			a.AgentShips[v.Owner] = make(map[string]bool)
			ships = a.AgentShips[v.Owner]
		}
		ships[v.Ship] = true
		a.lock.Unlock()
	}
}
