package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func queryProcessStatus(ctx context.Context) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.queryProcessStatus")
	defer span.End()

	caught = true
	if service.Process().IsBotProcessEnabled() {
		retMsg = "正常状态"
	} else {
		retMsg = "暂停状态"
	}
	return
}

func pauseProcess(ctx context.Context) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.pauseProcess")
	defer span.End()

	caught = true
	if !service.Process().IsBotProcessEnabled() {
		retMsg = "已处于暂停状态"
		return
	}
	if service.Process().PauseBotProcess() {
		retMsg = "已调至暂停状态"
		g.Log().Info(ctx, "Pause process")
	} else {
		retMsg = "调至暂停状态失败"
	}
	return
}

func continueProcess(ctx context.Context) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.continueProcess")
	defer span.End()

	caught = true
	if service.Process().IsBotProcessEnabled() {
		retMsg = "已处于正常状态"
		return
	}
	if service.Process().ContinueBotProcess() {
		retMsg = "已调至正常状态"
		g.Log().Info(ctx, "Continue process")
	} else {
		retMsg = "调至正常状态失败"
	}
	return
}
