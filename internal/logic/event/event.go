package event

import "qq-bot-backend/internal/service"

type sEvent struct{}

func New() *sEvent {
	return &sEvent{}
}

func init() {
	service.RegisterEvent(New())
}
