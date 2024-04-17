package cfg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
	"time"
)

type sCfg struct{}

func New() *sCfg {
	return &sCfg{}
}

func init() {
	service.RegisterCfg(New())
}

func (s *sCfg) GetRetryIntervalSeconds(ctx context.Context) time.Duration {
	seconds, err := g.Cfg().Get(ctx, "bot.retryIntervalSeconds")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if seconds == nil {
		return 3
	}
	return seconds.Duration()
}
