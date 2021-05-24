package meme

import (
	"fmt"
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

var responseMap = map[string]string{
	"omae wa mou shindeiru"	: "Nani?!",
	"omae wa mo shindeiru"	: "Nani?!",
	"omaewamoushindeiru"	: "Nani?!",
	"omaewamoshindeiru"	: "Nani?!",
	"お前はもう死んでいる"	: "なに?!",
}

var asciiExplosion = "```" +`.
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

func getKeys(mapItem map[string]string) []string {
    keys := make([]string, 0, len(mapItem))
    for k := range mapItem {
        keys = append(keys, k)
    }
	return keys
}

func ReplyNani(input string) string {
	i := strings.ToLower(input)
	response := "";
	mapKey := ""

	// Check if input string contains any of the map keys
    for _, key := range getKeys(responseMap) {
		if (strings.Contains(i, key)) {
			mapKey = key
		}
    }

	value, exists := responseMap[mapKey]

	if exists {
		response = fmt.Sprintf("%s %s", value, asciiExplosion)
	}

	return response
}
