package event

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryChainRecall(ctx context.Context) (catch bool) {
	messageIds, err := service.Bot().GetCachedMessageContext(ctx)
	if err != nil || len(messageIds) == 0 {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.TryChainRecall")
	defer span.End()

	for _, messageId := range messageIds {
		service.Bot().RecallMessage(ctx, messageId)
	}

	catch = true
	return
}
