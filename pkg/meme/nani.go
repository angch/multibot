package meme

import (
	"fmt"
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

var responseMap = map[string]string{
	"omae wa mou shindeiru": "Nani?!",
	"omae wa mo shindeiru":  "Nani?!",
	"omaewamoushindeiru":    "Nani?!",
	"omaewamoshindeiru":     "Nani?!",
	"お前はもう死んでいる":            "なに?!",
}

var asciiExplosion = "```" + `.
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

func init() {
	bothandler.RegisterCatchallHandler(ReplyNani)
}

func ReplyNani(input string) string {
	i := strings.ToLower(input)

	// Check if input string contains any of the map keys
	// Note: order of iteration is random, because map.
	for key, value := range responseMap {
		if strings.Contains(i, key) {
			return fmt.Sprintf("%s %s", value, asciiExplosion)
		}
	}
	return ""
}
