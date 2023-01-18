package state

import (
	"qq-bot-backend/internal/service"
)

type sState struct{}

func New() *sState {
	return &sState{}
}

func init() {
	service.RegisterState(New())
}
