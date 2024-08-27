package bot

import (
	"context"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

const (
	messageContextPrefix = "msg_ctx_"
	messageContextExpire = time.Minute*2 - time.Second*10
)

func getMessageContextCacheKey(userId, lastMessageId int64) string {
	return messageContextPrefix + gconv.String(userId) + "_" + gconv.String(lastMessageId)
}

func (s *sBot) CacheMessageContext(ctx context.Context, userId, lastMessageId, currentMessageId int64) error {
	return gcache.Set(ctx, getMessageContextCacheKey(userId, lastMessageId), currentMessageId, messageContextExpire)
}

func (s *sBot) GetCachedMessageContext(ctx context.Context, userId, lastMessageId int64,
) (currentMessageId int64, exist bool, err error) {
	v, err := gcache.Get(ctx, getMessageContextCacheKey(userId, lastMessageId))
	if err != nil || v == nil {
		return
	}

	exist = true
	currentMessageId = v.Int64()
	return
}
