package cfg

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

const (
	pathOTelEndpoint   = "otel.endpoint"
	pathOTelTraceToken = "otel.traceToken"
)

func (s *sCfg) GetOTel(ctx context.Context) (endpoint, traceToken string) {
	const def = ""
	e, err := gcfg.Instance().Get(ctx, pathOTelEndpoint, def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def, def
	}
	t, err := gcfg.Instance().Get(ctx, pathOTelTraceToken, def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return e.String(), def
	}
	return e.String(), t.String()
}
