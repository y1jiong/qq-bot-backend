package cfg

import (
	"qq-bot-backend/internal/service"
)

type sCfg struct{}

func New() *sCfg {
	return &sCfg{}
}

func init() {
	service.RegisterCfg(New())
}
