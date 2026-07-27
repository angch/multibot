package xkcd

import (
	"testing"

	"github.com/angch/multibot/pkg/bothandler"
)

func TestGetXKCD(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"303", "https://www.xkcd.com/303/"},
		{" 303 ", "https://www.xkcd.com/303/"},
		{"invalid", ""},
		{"-1", ""},
	}

	for _, tt := range tests {
		got := GetXKCD(bothandler.Request{Content: tt.input})
		if got != tt.want {
			t.Errorf("GetXKCD(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGetXKCDExplained(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"303", "https://www.explainxkcd.com/wiki/index.php?title=303"},
		{" 303 ", "https://www.explainxkcd.com/wiki/index.php?title=303"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		got := GetXKCDExplained(bothandler.Request{Content: tt.input})
		if got != tt.want {
			t.Errorf("GetXKCDExplained(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
