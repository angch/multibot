package meme

import (
	"testing"
	"strings"
)

var expectedAsciiExplosion = "```" +`.
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
        input  string
        want string
    }{
        {"omaewamoushindeiru", "Nani?!"},
        {"omaewamoshindeiru", "Nani?!"},
        {"omae wa mou shindeiru", "Nani?!"},
        {"omae wa mo shindeiru", "Nani?!"},
        {"OMAEWAMOSHINDEIRU", "Nani?!"},
        {"OmAe Wa MoU sHiNdEiRu", "Nani?!"},
        {"omaewamoushindeiru...", "Nani?!"},
        {"ooooomaewamoushindeiruuuuuu", "Nani?!"},
        {"お前はもう死んでいる", "なに?!"},
        {"なに?!お前はもう死んでいる!!!!", "なに?!"},
    }
    for _, tc := range testCases {
        if got := ReplyNani(tc.input); !strings.Contains(got, want) {
			t.Errorf("ReplyNani(\"%s\") = \"%s\"; want \"%s\"", tc.input, got, want)
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