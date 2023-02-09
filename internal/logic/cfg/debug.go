package cfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sCfg) IsDebugEnabled(ctx context.Context) bool {
	debug, err := g.Cfg().Get(ctx, "bot.debug")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if debug == nil {
		debug = gvar.New(false)
	}
	return debug.Bool()
}

func (s *sCfg) GetDebugToken(ctx context.Context) string {
	debugToken, err := g.Cfg().Get(ctx, "bot.debugToken")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if debugToken == nil {
		debugToken = gvar.New("")
	}
	return debugToken.String()
}
