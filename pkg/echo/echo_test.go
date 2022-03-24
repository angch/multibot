package echo

import (
	"testing"

	"github.com/angch/discordbot/pkg/bothandler"
)

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"1", "uwu", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EchoHandler(bothandler.Request{Content: tt.args, Platform: "", Channel: "", From: ""}); got != tt.want {
				t.Errorf("EchoHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
