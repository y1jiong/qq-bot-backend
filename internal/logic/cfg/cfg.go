package cfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"net/http"
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

func (s *sCfg) GetUrlPrefix(ctx context.Context) (urlPrefix string, err error) {
	v, err := g.Cfg().Get(ctx, "bot.urlPrefix")
	if err != nil {
		return
	}
	if v == nil {
		err = gerror.NewCode(gcode.New(http.StatusNotFound, "", nil), "urlPrefix not found")
		return
	}
	urlPrefix = v.String()
	return
}
