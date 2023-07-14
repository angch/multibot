package ymca

import (
	"math/rand"
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
)

var triggers = []string{
	"Y-M-C-A",
	"ymca",
	"YMCA!",
	"Y-M-C-A!",
}

var handles = []string{
	// "faz",
	// "tech_tarik",
}

var response_ascii_art = []string{
	"```" + `
      \   / /V\   r _  / \
       \O/  \O/  /O/   \O/
        #    #    #     #
       / \  / \  / \   / \
      ^   ^^   ^^   ^ ^   ^
` + "```",
	"```" + `
888  88888888b.d88b.  .d8888b 8888b.  
888  888888 "888 "88bd88P"       "88b 
888  888888  888  888888     .d888888 
Y88b 888888  888  888Y88b.   888  888 
 "Y88888888  888  888 "Y8888P"Y888888 
     888                              
Y8b d88P                              
 "Y88P" 
` + "```",
	"```" + `
	_   _ _ __ ___   ___ __ _ 
	| | | | '_ ' _ \ / __/ _' |
	| |_| | | | | | | (_| (_| |
	 \__, |_| |_| |_|\___\__,_|
	  __/ |                    
	 |___/                     ` + "```",
	"```" + `
8b       d8 88,dPYba,,adPYba,   ,adPPYba, ,adPPYYba,  
'8b     d8' 88P'   "88"    "8a a8"     "" ""     'Y8  
 '8b   d8'  88      88      88 8b         ,adPPPPP88  
  '8b,d8'   88      88      88 "8a,   ,aa 88,    ,88  
    Y88'    88      88      88  '"Ybbd8"' '"8bbdP"Y8  
    d8'                                               
   d8'                                                ` + "```",
}

func init() {
	bothandler.RegisterCatchallHandler(YMCAHandler)
}

func YMCAHandler(request bothandler.Request) string {
	input := request.Content

	i := strings.ToLower(input)

	count := 0
	for _, v := range triggers {
		if strings.Contains(i, strings.ToLower(v)) {
			count++
		}
	}
	uncount := 0
	for _, v := range handles {
		// Trigger is also the response, so some of them are in caps.
		if strings.Contains(i, strings.ToLower(v)) {
			uncount++
		}
	}

	if count >= 1 && uncount == 0 {
		pick := rand.Intn(len(response_ascii_art))
		return response_ascii_art[pick]
	}

	return ""
}
