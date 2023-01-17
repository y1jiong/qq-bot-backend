package cfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"he3-bot/internal/service"
	"time"
)

var (
	middlewareAccessIntervalMilliseconds *gvar.Var
	retryIntervalMilliseconds            *gvar.Var
	authToken                            *gvar.Var
)

type sCfg struct{}

func init() {
	service.RegisterCfg(New())
}

func New() *sCfg {
	return &sCfg{}
}

func (s *sCfg) logError(ctx context.Context, err error) {
	glog.Warningf(ctx, "an error occurred while get config from file. %v", err)
}

func (s *sCfg) GetAuthToken(ctx context.Context) string {
	var err error
	if authToken == nil {
		authToken, err = g.Cfg().Get(ctx, "bot.authToken")
		if err != nil {
			s.logError(ctx, err)
		}
		if authToken == nil {
			authToken = gvar.New("")
		}
	}
	return authToken.String()
}

func (s *sCfg) GetRetryIntervalMilliseconds(ctx context.Context) time.Duration {
	var err error
	if retryIntervalMilliseconds == nil {
		retryIntervalMilliseconds, err = g.Cfg().Get(ctx, "bot.retryIntervalMilliseconds")
		if err != nil {
			s.logError(ctx, err)
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
			s.logError(ctx, err)
		}
		if middlewareAccessIntervalMilliseconds == nil {
			middlewareAccessIntervalMilliseconds = gvar.New(2000)
		}
	}
	return middlewareAccessIntervalMilliseconds.Duration()
}
