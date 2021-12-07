package askfaz

import (
	"fmt"
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
	"credit",
	"serverless",
}

var handles = []string{
	"faz",
	"tech_tarik",
}

func init() {
	bothandler.RegisterCatchallHandler(AskFazHandler)
}

func AskFazHandler(r bothandler.Request) string {
	input := r.Content
	i := strings.ToLower(input)

	count := 0
	for _, v := range triggers {
		if strings.Contains(i, v) {
			count++
		}
	}

	uncount := 0
	for _, v := range handles {
		if strings.Contains(i, v) {
			uncount++
		}
	}

	if count >= 3 && uncount == 0 {
		return fmt.Sprintf("Good question, %s, what do you think?", strings.Join(handles[0:1], "/"))
	}

	return ""
}
