package qrcode

import (
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
	"github.com/skip2/go-qrcode"
)

func init() {
	bothandler.RegisterCatchallExtendeHandler(GetMessage)
}

func GetMessage(input bothandler.ExtendedMessage) *bothandler.ExtendedMessage {
	i := strings.ToLower(input.Text)

	if strings.HasPrefix(i, "/qrcode ") || strings.HasPrefix(i, "!qrcode ") {
		i = i[8:]
	} else {
		return nil
	}
	if i == "" {
		return nil
	}

	var err error
	_ = err
	png, err := qrcode.Encode(i, qrcode.Medium, 256)

	if err != nil {
		// It's likely an error

		return &bothandler.ExtendedMessage{
			Text:  "Zzzz server is sleeping",
			Image: nil,
		}
	}

	return &bothandler.ExtendedMessage{
		Text:  i,
		Image: png,
	}
}
