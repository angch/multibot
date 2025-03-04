package compreface

import (
	"log"
	"os"
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
)

func ComprefaceHandler(filename string, request bothandler.Request) string {
	if !strings.HasPrefix(request.Content, "!addface") {
		return ""
	}

	words := strings.Split(request.Content, " ")
	if len(words) < 2 {
		return "compreface: No subject"
	}
	subject := strings.Join(words[1:], " ")

	filedata, err := os.ReadFile(filename)
	if err != nil {
		return "compreface: Error opening file"
	}
	log.Println("File size is", len(filedata))
	output := botFaceRecognition.AddFace(subject, 1.0, filedata)

	return "Added " + string(output)
}

var botCompreface *Compreface
var botFaceRecognition *RecognitionService

func init() {
	botCompreface = New("http://localhost:8000")
	botFaceRecognition = botCompreface.InitFaceRecognition("92a551ee-3bf6-447d-97dd-824c61846192")

	bothandler.RegisterImageHandler(ComprefaceHandler)
}
