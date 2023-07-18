package spacetraders

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strings"
)

func (a *SpaceTraders) ProcessRegisterAgentResponse(ctx context.Context, pc PlatformChannel, resp *RegisterAgentResponse) {
	data := resp.Data
	agent := data.Agent
	if agent == nil {
		log.Println("No agent")
		return
	}
	faction := data.Faction
	if faction == nil {
		log.Println("No faction")
		return
	}

	if false {
		log.Printf("RegisterAgentResponse:\n Agent: %+v\nContract: %+v\nFaction: %+v\nShip: %+v\n", agent, data.Contract, faction, data.Ship)
	}

	agentState := &Agent{}

	gormdb := a.GormDB

	// We can't use gormdb.First or gormdb.Take because we don't have the primary key in ID
	err := gormdb.Find(agentState, "agent = ?", agent.Symbol).Error
	if err != nil {
		log.Println(err)
	}
	// If we found it, agentState.ID will be filled in.
	// Either way, we cram it with the "latest" values
	agentState.Agent = agent.Symbol
	agentState.Faction = faction.Symbol
	agentState.AuthToken = data.Token

	err = gormdb.Save(agentState).Error
	if err != nil {
		log.Println(err)
	}

	channelagent := &ChannelAgents{
		Platform:    pc.Platform,
		Channel:     pc.Channel,
		AgentSymbol: agent.Symbol,
	}
	err = gormdb.Save(channelagent).Error
	if err != nil {
		log.Println(err)
	}

	ag := &Agent{
		Agent:     agent.Symbol,
		Faction:   faction.Symbol,
		AuthToken: data.Token,
	}

	dbfaction := &Faction{
		Faction:      faction.Symbol,
		Name:         faction.Name,
		Description:  faction.Description,
		Headquarters: faction.Headquarters,
	}
	if faction.IsRecruiting != nil {
		dbfaction.IsRecruiting = *faction.IsRecruiting
	}
	tn := []string{}
	for _, t := range faction.Traits {
		dbtrait := &Trait{
			Trait:       t.Name,
			Description: t.Description,
		}
		a.SetTrait(*dbtrait)
		tn = append(tn, t.Name)
		dbfaction.Traits = append(dbfaction.Traits, *dbtrait)
	}
	sort.StringSlice(tn).Sort()
	dbfaction.TraitsString = strings.Join(tn, ",")
	a.SetFaction(*dbfaction)

	ag.Faction = faction.Symbol

	globalState[PlatformChannel{pc.Platform, pc.Channel}] = ag

	dbship := &Ship{}
	s, err := json.Marshal(data.Ship)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(s, dbship)
	if err != nil {
		log.Println(err)
	}
	dbship.Owner = agent.Symbol
	// log.Printf("ship is %+v\n", data.Ship)
	// log.Printf("s is %s\n", string(s))
	// log.Printf("dbship is %v\n", dbship)

	a.SetShip(*dbship)
}
