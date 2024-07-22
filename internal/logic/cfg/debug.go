package cfg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sCfg) IsDebugEnabled(ctx context.Context) bool {
	const def = false
	enabled, err := g.Cfg().Get(ctx, "bot.debug.enabled", def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return enabled.Bool()
}

func (s *sCfg) GetDebugToken(ctx context.Context) string {
	const def = ""
	token, err := g.Cfg().Get(ctx, "bot.debug.token", def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return token.String()
}
