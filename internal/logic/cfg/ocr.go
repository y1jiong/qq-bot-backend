package cfg

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

const pathOcrURL = "ocr.url"

func (s *sCfg) GetOcrURL(ctx context.Context) string {
	const def = ""
	url, err := gcfg.Instance().Get(ctx, pathOcrURL, def)
	if err != nil {
		g.Log().Warning(ctx, err)
		return def
	}
	return url.String()
}
