package bothandler

import "log"

// Implements MessagePlatform
type DevMessagePlatform struct {
}

func NewMessagePlatformFromDev() *DevMessagePlatform {
	return &DevMessagePlatform{}
}

func (s *DevMessagePlatform) Send(text string) {
	log.Println(text)
}
