package bot

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"time"
)

const (
	messageContextPrefix = "bot_msg_ctx_"
	messageContextTTL    = 2*time.Minute - 5*time.Second
)

func (s *sBot) generateEchoSignWithTrace(ctx context.Context) string {
	header := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(header))
	header["uid"] = guid.S()
	echoSign, err := sonic.MarshalString(header)
	if err != nil {
		return header["uid"]
	}
	return echoSign
}

func (s *sBot) extractEchoSign(ctx context.Context, echoSign string) context.Context {
	header := make(map[string]string)
	if err := sonic.UnmarshalString(echoSign, &header); err != nil {
		return ctx
	}
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(header))
}

func (s *sBot) getMessageContextCacheKey(ctx context.Context) string {
	return messageContextPrefix + gconv.String(s.GetSelfId(ctx)) + "_" + gconv.String(s.GetMsgId(ctx))
}

func (s *sBot) CacheMessageContext(ctx context.Context, nextMessageId int64) error {
	cacheKey := s.getMessageContextCacheKey(ctx)
	v, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}

	arr := v.Int64s()
	if arr == nil {
		arr = make([]int64, 0, 1)
	}
	arr = append(arr, nextMessageId)

	return gcache.Set(ctx, cacheKey, arr, messageContextTTL)
}

func (s *sBot) GetCachedMessageContext(ctx context.Context) (nextMessageIds []int64, err error) {
	v, err := gcache.Get(ctx, s.getMessageContextCacheKey(ctx))
	if err != nil || v == nil {
		return
	}

	return v.Int64s(), nil
}
