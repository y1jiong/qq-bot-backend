package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
)

func queryProcessStatus(ctx context.Context) (catch bool) {
	catch = true
	// 继续处理
	if service.Process().IsBotProcess() {
		service.Bot().SendPlainMsg(ctx, "正常状态")
	} else {
		service.Bot().SendPlainMsg(ctx, "暂停状态")
	}
	return
}

func pauseProcess(ctx context.Context) (catch bool) {
	catch = true
	// 继续处理
	if !service.Process().IsBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已处于暂停状态")
		return
	}
	if service.Process().PauseBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已调至暂停状态")
		g.Log().Info(ctx, "Pause process")
	} else {
		service.Bot().SendPlainMsg(ctx, "调至暂停状态失败")
	}
	return
}

func continueProcess(ctx context.Context) (catch bool) {
	catch = true
	// 继续处理
	if service.Process().IsBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已处于正常状态")
		return
	}
	if service.Process().ContinueBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已调至正常状态")
		g.Log().Info(ctx, "Continue process")
	} else {
		service.Bot().SendPlainMsg(ctx, "调至正常状态失败")
	}
	return
}
