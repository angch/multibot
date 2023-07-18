package spacetraders

import (
	"encoding/json"
	"time"
)

type RegisterAgentRequest struct {
	Symbol  string `json:"symbol"`
	Faction string `json:"faction"`
}

type RegisterAgentResponse struct {
	Data  *RegisterAgentResponseData `json:"data"`
	Error *ErrorResponse             `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    json.Number `json:"code"`
}

type RegisterAgentResponseData struct {
	Token    string        `json:"token,omitempty"`
	Agent    *AgentData    `json:"agent,omitempty"`
	Contract *ContractData `json:"contract,omitempty"`
	Faction  *FactionData  `json:"faction,omitempty"`
	Ship     *ShipData     `json:"ship,omitempty"`
}

type AgentData struct {
	AccountID       string `json:"accountId"`
	Symbol          string `json:"symbol"`
	Headquarters    string `json:"headquarters"`
	Credits         int    `json:"credits"`
	StartingFaction string `json:"startingFaction"`
}

type PaymentData struct {
	OnAccepted  int `json:"onAccepted"`
	OnFulfilled int `json:"onFulfilled"`
}

type DeliverData struct {
	TradeSymbol       string `json:"tradeSymbol"`
	DestinationSymbol string `json:"destinationSymbol"`
	UnitsRequired     int    `json:"unitsRequired"`
	UnitsFulfilled    int    `json:"unitsFulfilled"`
}

type TermsData struct {
	Deadline time.Time     `json:"deadline"`
	Payment  PaymentData   `json:"payment"`
	Deliver  []DeliverData `json:"deliver"`
}

type ContractData struct {
	ID               string    `json:"id"`
	FactionSymbol    string    `json:"factionSymbol"`
	Type             string    `json:"type"`
	Terms            TermsData `json:"terms"`
	Accepted         bool      `json:"accepted"`
	Fulfilled        bool      `json:"fulfilled"`
	Expiration       time.Time `json:"expiration"`
	DeadlineToAccept time.Time `json:"deadlineToAccept"`
}

type TraitsData struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FactionData struct {
	Symbol       string       `json:"symbol"`
	Name         string       `json:"name,omitempty"`
	Description  string       `json:"description,omitempty"`
	Headquarters string       `json:"headquarters,omitempty"`
	Traits       []TraitsData `json:"traits,omitempty"`
	IsRecruiting *bool        `json:"isRecruiting,omitempty"`
}

type WaypointData struct {
	Symbol       string `json:"symbol"`
	Type         string `json:"type"`
	SystemSymbol string `json:"systemSymbol"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

type RouteData struct {
	Departure     WaypointData `json:"departure"`
	Destination   WaypointData `json:"destination"`
	Arrival       time.Time    `json:"arrival"`
	DepartureTime time.Time    `json:"departureTime"`
}

type ShipFrameData struct {
	Symbol         string `json:"symbol"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ModuleSlots    int    `json:"moduleSlots"`
	MountingPoints int    `json:"mountingPoints"`
	FuelCapacity   int    `json:"fuelCapacity"`
	Condition      int    `json:"condition"`
	Requirements   struct {
		Power int `json:"power"`
		Crew  int `json:"crew"`
	} `json:"requirements"`
}

type ShipModuleData struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Capacity     int    `json:"capacity,omitempty"`
	Requirements struct {
		Crew  int `json:"crew"`
		Power int `json:"power"`
		Slots int `json:"slots"`
	} `json:"requirements"`
	Range int `json:"range,omitempty"`
}

type ShipMountData struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Strength     int    `json:"strength"`
	Requirements struct {
		Crew  int `json:"crew"`
		Power int `json:"power"`
	} `json:"requirements"`
	Deposits []string `json:"deposits,omitempty"`
}

type NavData struct {
	SystemSymbol   string    `json:"systemSymbol"`
	WaypointSymbol string    `json:"waypointSymbol"`
	Route          RouteData `json:"route"`
	Status         string    `json:"status"`
	FlightMode     string    `json:"flightMode"`
}

type CrewData struct {
	Current  int    `json:"current"`
	Capacity int    `json:"capacity"`
	Required int    `json:"required"`
	Rotation string `json:"rotation"`
	Morale   int    `json:"morale"`
	Wages    int    `json:"wages"`
}

type FuelData struct {
	Current  int `json:"current"`
	Capacity int `json:"capacity"`
	Consumed struct {
		Amount    int       `json:"amount"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"consumed"`
}

type ShipReactorData struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Condition    int    `json:"condition"`
	PowerOutput  int    `json:"powerOutput"`
	Requirements struct {
		Crew int `json:"crew"`
	} `json:"requirements"`
}

type ShipEngineData struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Condition    int    `json:"condition"`
	Speed        int    `json:"speed"`
	Requirements struct {
		Power int `json:"power"`
		Crew  int `json:"crew"`
	} `json:"requirements"`
}

type ShipRegistrationData struct {
	Name          string `json:"name"`
	FactionSymbol string `json:"factionSymbol"`
	Role          string `json:"role"`
}

type ShipCargoData struct {
	Capacity  int             `json:"capacity"`
	Units     int             `json:"units"`
	Inventory []InventoryData `json:"inventory"`
}

type ShipData struct {
	Symbol       string               `json:"symbol"`
	Nav          NavData              `json:"nav"`
	Crew         CrewData             `json:"crew"`
	Fuel         FuelData             `json:"fuel"`
	Frame        ShipFrameData        `json:"frame"`
	Reactor      ShipReactorData      `json:"reactor"`
	Engine       ShipEngineData       `json:"engine"`
	Modules      []ShipModuleData     `json:"modules"`
	Mounts       []ShipMountData      `json:"mounts"`
	Registration ShipRegistrationData `json:"registration"`
	Cargo        ShipCargoData        `json:"cargo"`
}

type InventoryData any

type SystemData struct {
	SystemSymbol string `json:"systemSymbol"`
	Symbol       string `json:"symbol"`
	Type         string `json:"type"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
	Orbitals     []struct {
		Symbol string `json:"symbol"`
	} `json:"orbitals"`
	Traits []TraitsData `json:"traits"`
	Chart  struct {
		SubmittedBy string    `json:"submittedBy"`
		SubmittedOn time.Time `json:"submittedOn"`
	} `json:"chart"`
	Faction FactionData `json:"faction"` // Note: only Symbol given
}

type SystemResponse struct {
	Data SystemData `json:"data"`
}
