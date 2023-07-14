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
		{"2", `Taking caffeine is basically borrowing energy from your future self and you need to pay it back with interest (haram lol)…
		One thing that really motivated me to stop coffee intake is that coffee masks the root problem.
		For example, you feel tired everyday and you solved it with drinking coffee while actually, the real problem is that you don’t exercise enough or constantly watching youtube before bed or don’t eat healthily.`, "Yes, we know."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EchoHandler(bothandler.Request{Content: tt.args, Platform: "", Channel: "", From: ""}); got != tt.want {
				t.Errorf("EchoHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
