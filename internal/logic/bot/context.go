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
	messageContextPrefix = "msg_ctx_"
	messageContextExpire = time.Minute*2 - time.Second*10
)

func (s *sBot) generateEchoSignWithTrace(ctx context.Context) string {
	header := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(header))
	header["uid"] = guid.S()
	echoSign, err := sonic.ConfigDefault.MarshalToString(header)
	if err != nil {
		return header["uid"]
	}
	return echoSign
}

func (s *sBot) extractEchoSign(ctx context.Context, echoSign string) context.Context {
	header := make(map[string]string)
	if err := sonic.ConfigDefault.UnmarshalFromString(echoSign, &header); err != nil {
		return ctx
	}
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(header))
	return ctx
}

func getMessageContextCacheKey(userId, lastMessageId int64) string {
	return messageContextPrefix + gconv.String(userId) + "_" + gconv.String(lastMessageId)
}

func (s *sBot) CacheMessageContext(ctx context.Context, userId, lastMessageId, currentMessageId int64) error {
	cacheKey := getMessageContextCacheKey(userId, lastMessageId)
	v, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}

	arr := v.Int64s()
	if arr == nil {
		arr = make([]int64, 0, 1)
	}
	arr = append(arr, currentMessageId)

	return gcache.Set(ctx, cacheKey, arr, messageContextExpire)
}

func (s *sBot) GetCachedMessageContext(ctx context.Context, userId, lastMessageId int64,
) (currentMessageIds []int64, exist bool, err error) {
	v, err := gcache.Get(ctx, getMessageContextCacheKey(userId, lastMessageId))
	if err != nil || v == nil {
		return
	}

	exist = true
	currentMessageIds = v.Int64s()
	return
}
