package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
)

func queryProcessState(ctx context.Context) (catch bool) {
	catch = true
	
	if service.State().IsBotProcess() {
		service.Bot().SendMsg(ctx, "处于正常处理所有信息的状态")
	} else {
		service.Bot().SendMsg(ctx, "处于暂停处理所有信息的状态")
	}
	return
}

func pauseProcess(ctx context.Context) (catch bool) {
	catch = true

	if !service.State().IsBotProcess() {
		service.Bot().SendMsg(ctx, "已处于暂停处理所有信息的状态")
		return
	}
	if service.State().PauseBotProcess() {
		service.Bot().SendMsg(ctx, "已暂停处理所有信息")
		g.Log().Info(ctx, "Pause parse")
	} else {
		service.Bot().SendMsg(ctx, "暂停处理所有信息失败")
	}
	return
}

func continueProcess(ctx context.Context) (catch bool) {
	catch = true

	if service.State().IsBotProcess() {
		service.Bot().SendMsg(ctx, "已处于正常处理所有信息的状态")
		return
	}
	if service.State().ContinueBotProcess() {
		service.Bot().SendMsg(ctx, "已恢复处理所有信息")
		g.Log().Info(ctx, "Continue parse")
	} else {
		service.Bot().SendMsg(ctx, "恢复处理所有信息失败")
	}
	return
}
