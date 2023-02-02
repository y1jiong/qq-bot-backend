package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryRaw(ctx context.Context, cmd string) (catch bool) {
	// 权限校验
	if !service.User().CouldGetRawMsg(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	service.Bot().SendPlainMsg(ctx, cmd)
	catch = true
	return
}
