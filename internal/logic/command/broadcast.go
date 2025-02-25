package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func tryBroadcast(ctx context.Context, cmd string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryBroadcast")
	defer span.End()

	groupId := service.Bot().GetGroupId(ctx)
	namespace := service.Group().GetNamespace(ctx, groupId)
	if namespace == "" ||
		!service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}

	caught = true

	if err := service.Namespace().Broadcast(ctx, namespace, cmd, groupId); err != nil {
		retMsg = "广播失败"
	}
	return
}
