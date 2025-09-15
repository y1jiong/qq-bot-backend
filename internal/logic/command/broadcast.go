package command

import (
	"context"
	"qq-bot-backend/internal/service"

	"github.com/gogf/gf/v2/net/gtrace"
)

func tryBroadcast(ctx context.Context, cmd string) (caught catch, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.broadcast")
	defer span.End()

	groupId := service.Bot().GetGroupId(ctx)
	namespace := service.Group().GetNamespace(ctx, groupId)
	if namespace == "" ||
		!service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}

	caught = caughtNeedOkay

	// /broadcast <message>
	if err := service.Namespace().Broadcast(ctx, namespace, cmd, groupId); err != nil {
		retMsg = "广播失败"
	}
	return
}
