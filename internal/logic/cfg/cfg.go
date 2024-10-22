package cfg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
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

const (
	pathBotRetryIntervalSeconds = "bot.retryIntervalSeconds"
)

func (s *sCfg) GetRetryIntervalSeconds(ctx context.Context) time.Duration {
	const def = 3
	seconds, err := gcfg.Instance().Get(ctx, pathBotRetryIntervalSeconds, def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return seconds.Duration()
}
