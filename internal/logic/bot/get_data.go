package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
)

func (s *sBot) getData(ctx context.Context) *sj.Json {
	return s.reqJsonFromCtx(ctx).Get("data")
}

func (s *sBot) getFileFromData(ctx context.Context) string {
	return s.getData(ctx).Get("file").MustString()
}
