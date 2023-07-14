package meme

import (
	"fmt"
	"strings"
	"testing"

	"github.com/angch/multibot/pkg/bothandler"
)

var expectedAsciiExplosion = "```" + `.
      _.-^^---....,,--
  _--                  --_  
 <                        >)
 |                         | 
  \._                   _./  
     '''--. . , ; .--'''       
           | |   |             
        .-=||  | |=-.   
        '-=#$%&%$#=-'   
           | ;   |     
  _____.,-#%&$@%#&#~,._____
` + "```"

func TestReplyNani(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"omaewamoushindeiru", "Nani?!"},
		{"omaewamoshindeiru", "Nani?!"},
		{"omae wa mou shindeiru", "Nani?!"},
		{"omae wa mo shindeiru", "Nani?!"},
		{"OMAEWAMOSHINDEIRU", "Nani?!"},
		{"OmAe Wa MoU sHiNdEiRu", "Nani?!"},
		{"omaewamoushindeiru...", "Nani?!"},
		{"ooooomaewamoushindeiruuuuuu", "Nani?!"},
		{"お前はもう死んでいる", "なに？！"},
		{"なに?!お前はもう死んでいる!!!!", "なに？！"},
	}
	for _, tc := range testCases {
		want := fmt.Sprintf("%s %s", tc.want, expectedAsciiExplosion)
		r := bothandler.Request{Content: tc.input, Platform: "", Channel: "", From: ""}
		if got := ReplyNani(r); !strings.Contains(got, want) {
			t.Errorf("ReplyNani(\"%s\") = \"%s\"; want \"%s\"", tc.input, got, want)
		}
	}
}

func TestReplyNaniNoMatches(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"omae wa mou shin deiru", ""},
		{"onigiri", ""},
		{"お???", ""},
	}
	for _, tc := range testCases {
		r := bothandler.Request{Content: tc.input, Platform: "", Channel: "", From: ""}
		if got := ReplyNani(r); got != tc.want {
			t.Errorf("ReplyNani(\"%s\") = \"%s\"; want \"%s\"", tc.input, got, tc.want)
		}
	}
}
