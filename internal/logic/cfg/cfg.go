package cfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
	"time"
)

var (
	middlewareAccessIntervalMilliseconds *gvar.Var
	retryIntervalMilliseconds            *gvar.Var
	debugToken                           *gvar.Var
	debug                                *gvar.Var
)

type sCfg struct{}

func New() *sCfg {
	return &sCfg{}
}

func init() {
	service.RegisterCfg(New())
}

func (s *sCfg) IsDebugEnabled(ctx context.Context) bool {
	var err error
	if debug == nil {
		debug, err = g.Cfg().Get(ctx, "bot.debug")
		if err != nil {
			g.Log().Warning(ctx, err)
		}
		if debug == nil {
			debug = gvar.New(false)
		}
	}
	return debug.Bool()
}

func (s *sCfg) GetDebugToken(ctx context.Context) string {
	var err error
	if debugToken == nil {
		debugToken, err = g.Cfg().Get(ctx, "bot.debugToken")
		if err != nil {
			g.Log().Warning(ctx, err)
		}
		if debugToken == nil {
			debugToken = gvar.New("")
		}
	}
	return debugToken.String()
}

func (s *sCfg) GetRetryIntervalMilliseconds(ctx context.Context) time.Duration {
	var err error
	if retryIntervalMilliseconds == nil {
		retryIntervalMilliseconds, err = g.Cfg().Get(ctx, "bot.retryIntervalMilliseconds")
		if err != nil {
			g.Log().Warning(ctx, err)
		}
		if retryIntervalMilliseconds == nil {
			retryIntervalMilliseconds = gvar.New(3000)
		}
	}
	return retryIntervalMilliseconds.Duration()
}

func (s *sCfg) GetMiddlewareAccessIntervalMilliseconds(ctx context.Context) time.Duration {
	var err error
	if middlewareAccessIntervalMilliseconds == nil {
		middlewareAccessIntervalMilliseconds, err = g.Cfg().Get(ctx, "api.middlewareAccessIntervalMilliseconds")
		if err != nil {
			g.Log().Warning(ctx, err)
		}
		if middlewareAccessIntervalMilliseconds == nil {
			middlewareAccessIntervalMilliseconds = gvar.New(2000)
		}
	}
	return middlewareAccessIntervalMilliseconds.Duration()
}
