package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
)

func (s *sBot) GetData(ctx context.Context) *sj.Json {
	return s.reqJsonFromCtx(ctx).Get("data")
}

func (s *sBot) GetFileFromData(ctx context.Context) string {
	return s.GetData(ctx).Get("file").MustString()
}

func (s *sBot) GetSenderFromData(ctx context.Context) (nickname string, userId int64) {
	nickname = s.GetData(ctx).Get("sender").Get("nickname").MustString()
	userId = s.GetData(ctx).Get("sender").Get("user_id").MustInt64()
	return
}

func (s *sBot) GetMessageFromData(ctx context.Context) string {
	return s.GetData(ctx).Get("message").MustString()
}
