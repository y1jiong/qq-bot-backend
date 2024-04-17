package cfg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sCfg) IsEnabledDebug(ctx context.Context) bool {
	enabled, err := g.Cfg().Get(ctx, "bot.debug.enabled")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if enabled == nil {
		return false
	}
	return enabled.Bool()
}

func (s *sCfg) GetDebugToken(ctx context.Context) string {
	token, err := g.Cfg().Get(ctx, "bot.debug.token")
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	if token == nil {
		return ""
	}
	return token.String()
}
