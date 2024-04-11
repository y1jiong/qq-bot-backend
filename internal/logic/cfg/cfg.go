package cfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
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
	retryIntervalSeconds, err := g.Cfg().Get(ctx, "bot.retryIntervalSeconds")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if retryIntervalSeconds == nil {
		retryIntervalSeconds = gvar.New(3)
	}
	return retryIntervalSeconds.Duration()
}
