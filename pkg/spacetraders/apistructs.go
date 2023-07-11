package spacetraders

type RegisterAgentRequest struct {
	Symbol  string `json:"username"`
	Faction string `json:"faction"`
}
