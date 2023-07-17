package spacetraders

import (
	"log"

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

const savefile string = "spacetraders.sqlite"

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
	)
	this.GormDB = gormdb

	// FIXME: slurp everything into globalState
	ca := make([]ChannelAgents, 0)
	gormdb.Find(&ca)
	for _, c := range ca {
		ag := &Agent{}
		gormdb.First(ag, "symbol = ?", c.AgentSymbol)
		globalState[PlatformChannel{c.Platform, c.Channel}] = ag
	}
}

func save() {
	lock.Lock()
	// FIXME: dump everything from globalState into spacetraders.sqlite
	defer lock.Unlock()
}
