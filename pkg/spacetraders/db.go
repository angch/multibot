package spacetraders

import (
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

const savefile string = "spacetraders.sqlite"

var lock sync.Mutex

func load() {
	lock.Lock()
	defer lock.Unlock()

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

func (a *SpaceTraders) LoadKnown() {
	traits := make([]Trait, 0)
	factions := make([]Faction, 0)

	a.GormDB.Find(&traits)   // FIXME
	a.GormDB.Find(&factions) // FIXME
	for _, v := range traits {
		a.KnownTraits[v.Trait] = v
	}
	for _, v := range factions {
		v.Fix(a)
		a.KnownFactions[v.Faction] = v
	}
}
