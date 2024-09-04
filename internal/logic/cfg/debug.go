package cfg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

func (s *sCfg) IsDebugEnabled(ctx context.Context) bool {
	const def = false
	enabled, err := gcfg.Instance().Get(ctx, "bot.debug.enabled", def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return enabled.Bool()
}

func (s *sCfg) GetDebugToken(ctx context.Context) string {
	const def = ""
	token, err := gcfg.Instance().Get(ctx, "bot.debug.token", def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return token.String()
}
