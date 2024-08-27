package event

import (
	"context"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryChainRecall(ctx context.Context) (catch bool) {
	msgIds, exist, err := service.Bot().GetCachedMessageContext(ctx,
		service.Bot().GetUserId(ctx),
		service.Bot().GetMsgId(ctx),
	)
	if err != nil || !exist {
		return
	}
	for _, msgId := range msgIds {
		service.Bot().RecallMessage(ctx, msgId)
	}
	catch = true
	return
}
