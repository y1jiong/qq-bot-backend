package cfg

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

const (
	pathBotDebugEnabled = "bot.debug.enabled"
	pathBotDebugToken   = "bot.debug.token"
)

func (s *sCfg) IsDebugEnabled(ctx context.Context) bool {
	const def = false
	enabled, err := gcfg.Instance().Get(ctx, pathBotDebugEnabled, def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return enabled.Bool()
}

func (s *sCfg) GetDebugToken(ctx context.Context) string {
	const def = ""
	token, err := gcfg.Instance().Get(ctx, pathBotDebugToken, def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return token.String()
}
