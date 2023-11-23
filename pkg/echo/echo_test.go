package echo

import (
	"testing"

	"github.com/angch/multibot/pkg/bothandler"
)

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"1", "uwu", "(ꈍᴗꈍ)"},
		{"2", "Plus caffeine only delays the tiredness. It doesn't prevent it.", "Yes, we know."},
		{"3", `Taking caffeine is basically borrowing energy from your future self and you need to pay it back with interest (haram lol)…
		One thing that really motivated me to stop coffee intake is that coffee masks the root problem.
		For example, you feel tired everyday and you solved it with drinking coffee while actually, the real problem is that you don’t exercise enough or constantly watching youtube before bed or don’t eat healthily.`, "Yes, we know."},
		{"4a", "Any Java experts around?", "https://dontasktoask.com/"},
		{"4b", "Any Java experts around who are willing to commit into looking into my problem, whatever that may turn out to be, even if it's not actually related to Java or if someone who doesn't know anything about Java could actually answer my question?", ""},
		{"5", "can i advice you something?", "No!"},
		{"6", "Anyone can helps me on something's? Very urgent and Idk how to solve", "https://dontasktoask.com/"},
		{"7a", "agi", ""},
		{"7b", "AGI", "*Feel* the AGI!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EchoHandler(bothandler.Request{Content: tt.args, Platform: "", Channel: "", From: ""}); got != tt.want {
				t.Errorf("EchoHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
