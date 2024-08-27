package event

import (
	"context"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryChainRecall(ctx context.Context) (catch bool) {
	msgId, exist, err := service.Bot().GetCachedMessageContext(ctx,
		service.Bot().GetUserId(ctx),
		service.Bot().GetMsgId(ctx),
	)
	if err != nil || !exist {
		return
	}
	service.Bot().RecallMessage(ctx, msgId)
	catch = true
	return
}
