package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
)

func queryProcessState(ctx context.Context) (catch bool) {
	catch = true
	// 继续处理
	if service.Process().IsBotProcess() {
		service.Bot().SendPlainMsg(ctx, "处于正常处理所有信息的状态")
	} else {
		service.Bot().SendPlainMsg(ctx, "处于暂停处理所有信息的状态")
	}
	return
}

func pauseProcess(ctx context.Context) (catch bool) {
	catch = true
	// 继续处理
	if !service.Process().IsBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已处于暂停处理所有信息的状态")
		return
	}
	if service.Process().PauseBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已暂停处理所有信息")
		g.Log().Info(ctx, "Pause process")
	} else {
		service.Bot().SendPlainMsg(ctx, "暂停处理所有信息失败")
	}
	return
}

func continueProcess(ctx context.Context) (catch bool) {
	catch = true
	// 继续处理
	if service.Process().IsBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已处于正常处理所有信息的状态")
		return
	}
	if service.Process().ContinueBotProcess() {
		service.Bot().SendPlainMsg(ctx, "已恢复处理所有信息")
		g.Log().Info(ctx, "Continue process")
	} else {
		service.Bot().SendPlainMsg(ctx, "恢复处理所有信息失败")
	}
	return
}
