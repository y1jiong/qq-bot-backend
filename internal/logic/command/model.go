package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func tryModelSet(ctx context.Context, args []string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.tryModelSet")
	defer span.End()

	// /model <op> <model>
	if len(args) < 2 {
		return
	}
	if args[0] != "set" {
		return
	}

	caught = true

	// /model set <model>
	if err := service.Bot().SetModel(ctx, args[1]); err != nil {
		retMsg = err.Error()
		return
	}
	retMsg = "已更改机型为 '" + args[1] + "'"
	return
}
