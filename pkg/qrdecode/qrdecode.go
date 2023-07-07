package qrdecode

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"strings"

	"github.com/liyue201/goqr"

	"github.com/angch/discordbot/pkg/bothandler"
)

func QrdecodeHandler(filename string, request bothandler.Request) string {
	imgdata, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return ""
	}

	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		log.Printf("image.Decode error: %v\n", err)
		return ""
	}
	// // Set the expected size that you want:
	// ratio := src.Bounds().Max.X / src.Bounds().Max.Y

	// img2 := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

	// // Resize:
	// draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Printf("Recognize failed: %v\n", err)
		return ""
	}
	out := ""
	for _, qrCode := range qrCodes {
		log.Printf("qrCode text: %s\n", qrCode.Payload)
		out += string(qrCode.Payload) + "\n"
	}

	out2 := strings.TrimSpace(string(out))
	return "QRCode: " + out2
}

func init() {
	bothandler.RegisterImageHandler(QrdecodeHandler)
}
