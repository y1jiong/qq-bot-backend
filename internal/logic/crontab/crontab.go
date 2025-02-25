package crontab

import (
	"context"
	"qq-bot-backend/internal/service"
)

type sCrontab struct{}

func New() *sCrontab {
	return &sCrontab{}
}

func init() {
	service.RegisterCrontab(New())
}

func (s *sCrontab) Run(ctx context.Context) {

}
