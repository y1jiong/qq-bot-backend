package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryForward(ctx context.Context) (catch bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryForward")
	defer span.End()

	groupId := service.Bot().GetGroupId(ctx)
	userId := service.Bot().GetUserId(ctx)
	if groupId != 0 {
		if !service.Namespace().IsForwardingMatchGroupId(ctx, gconv.String(groupId)) {
			return
		}
	} else if userId != 0 && !service.Namespace().IsForwardingMatchUserId(ctx, gconv.String(userId)) {
		return
	}
	aliasList := service.Namespace().GetForwardingToAliasList(ctx)
	for alias := range aliasList {
		url, key := service.Namespace().GetForwardingTo(ctx, alias)
		if url == "" {
			continue
		}
		err := service.Bot().Forward(ctx, url, key)
		if err != nil {
			g.Log().Notice(ctx, "forward", url, err)
		}

		catch = true
	}
	return
}
