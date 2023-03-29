package cfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

func (s *sCfg) GetRetryIntervalMilliseconds(ctx context.Context) time.Duration {
	retryIntervalMilliseconds, err := g.Cfg().Get(ctx, "bot.retryIntervalMilliseconds")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if retryIntervalMilliseconds == nil {
		retryIntervalMilliseconds = gvar.New(3000)
	}
	return retryIntervalMilliseconds.Duration()
}
