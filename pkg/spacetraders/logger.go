package spacetraders

import (
	"context"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type RequestLog struct {
	gorm.Model
	Platform string `gorm:"column:platform" json:"platform"`
	Channel  string `gorm:"column:channel" json:"channel"`

	Type string `gorm:"column:type" json:"type"`
	URL  string `gorm:"column:url" json:"url"`

	Data string `gorm:"column:data" json:"data"`

	ResponseStatusCode int    `gorm:"column:response_status_code" json:"response_status_code"`
	Response           string `gorm:"column:response" json:"response"`
}

func (m *SpaceTraders) LogRequest(ctx context.Context, rl *RequestLog, r *http.Request) {
	rl.URL = r.URL.String()
	_ = ctx
	err := m.GormDB.Create(rl).Error
	if err != nil {
		log.Println(err)
	}
}
