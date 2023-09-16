package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryRaw(ctx context.Context, cmd string) (catch bool, retMsg string) {
	// 权限校验
	if !service.User().CouldGetRawMsg(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	retMsg = cmd
	catch = true
	return
}
