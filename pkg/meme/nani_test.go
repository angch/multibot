package meme

import (
	"testing"
	"strings"
)

func TestReplyNani(t *testing.T) {
    testCases := []struct {
        input  string
        want string
    }{
        {"omaewamoushindeiru", "Nani?! https://tenor.com/"},
        {"omaewamoshindeiru", "Nani?! https://tenor.com/"},
        {"omae wa mou shindeiru", "Nani?! https://tenor.com/"},
        {"omae wa mo shindeiru", "Nani?! https://tenor.com/"},
        {"OMAEWAMOSHINDEIRU", "Nani?! https://tenor.com/"},
        {"OmAe Wa MoU sHiNdEiRu", "Nani?! https://tenor.com/"},
        {"omaewamoushindeiru...", "Nani?! https://tenor.com/"},
        {"ooooomaewamoushindeiruuuuuu", "Nani?! https://tenor.com/"},
        {"お前はもう死んでいる", "なに?! https://tenor.com/"},
        {"なに?!お前はもう死んでいる!!!!", "なに?! https://tenor.com/"},
    }
    for _, tc := range testCases {
        if got := ReplyNani(tc.input); !strings.Contains(got, tc.want) {
			t.Errorf("ReplyNani(\"%s\") = \"%s\"; should include string \"%s\"", tc.input, got, tc.want)
        }
    }
}

func TestReplyNaniNoMatches(t *testing.T) {
    testCases := []struct {
        input  string
        want string
    }{
        {"omae wa mou shin deiru", ""},
        {"onigiri", ""},
        {"お???", ""},
    }
    for _, tc := range testCases {
        if got := ReplyNani(tc.input); got != tc.want {
			t.Errorf("ReplyNani(\"%s\") = \"%s\"; want \"%s\"", tc.input, got, tc.want)
        }
    }
}