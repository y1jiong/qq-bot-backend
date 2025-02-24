package event

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryChainRecall(ctx context.Context) (caught bool) {
	messageIds, err := service.Bot().GetCachedMessageContext(ctx)
	if err != nil || len(messageIds) == 0 {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.TryChainRecall")
	defer span.End()

	for _, messageId := range messageIds {
		service.Bot().RecallMessage(ctx, messageId)
	}

	caught = true
	return
}

func (s *sEvent) TryEmojiRecall(ctx context.Context) (caught bool) {
	likes := service.Bot().GetLikes(ctx)
	if len(likes) == 0 {
		return
	}
	if !service.User().CanRecallMessage(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.TryEmojiRecall")
	defer span.End()

	for _, like := range likes {
		if gconv.String(like["emoji_id"]) == "326" { // 326: 机器人生气
			service.Bot().RecallMessage(ctx, service.Bot().GetMsgId(ctx))
			caught = true
			break
		}
	}

	return
}
