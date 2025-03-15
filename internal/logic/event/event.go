package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
)

type sEvent struct{}

func New() *sEvent {
	return &sEvent{}
}

func init() {
	service.RegisterEvent(New())
}

func (s *sEvent) TryCacheMessageAstNode(ctx context.Context) {
	if err := service.Bot().CacheMessageAstNode(ctx); err != nil {
		g.Log().Error(ctx, err)
	}
}
