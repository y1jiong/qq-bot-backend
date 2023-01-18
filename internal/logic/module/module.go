package module

import "qq-bot-backend/internal/service"

type sModule struct{}

func New() *sModule {
	return &sModule{}
}

func init() {
	service.RegisterModule(New())
}
