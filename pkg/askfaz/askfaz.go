package askfaz

import (
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

var triggers = []string{
	"cloud",
	"?",
	"recommend",
	"architecture",
	"how to",
	"best practice",
	"aws",
	"amazon",
}

func init() {
	bothandler.RegisterCatchallHandler(AskFazHandler)
}

func AskFazHandler(input string) string {
	i := strings.ToLower(input)

	count := 0
	for _, v := range triggers {
		if strings.Contains(i, v) {
			count++
		}
	}
	if count > 3 {
		return "Good question, @tech_tarik, what do you think?"
	}

	return ""
}
