package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
)

func (s *sBot) GetData(ctx context.Context) *sj.Json {
	return s.reqJsonFromCtx(ctx).Get("data")
}

func (s *sBot) GetFile(ctx context.Context) string {
	return s.GetData(ctx).Get("file").MustString()
}
