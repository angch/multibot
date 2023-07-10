package engineersmy

// FIXME: Put this outside as a file

var KnownDiscordChannels = map[string]string{
	"offtopic":     "811472319876562991",
	"general":      "811472319876562989",
	"":             "811472319876562989",
	"spacetraders": "1127471366501834763",
	"sandbox":      "846297823624298517",
}

var KnownTelegramChannels = map[string]int{
	"offtopic": -1001430213215,
	"general":  -1001430213215,
	// "":         -1001430213215,
}

func IsKnownDiscordChannel(channelname string, channel string) bool {
	id, ok := KnownDiscordChannels[channelname]
	if ok && id == channel {
		return true
	}
	return false
}
