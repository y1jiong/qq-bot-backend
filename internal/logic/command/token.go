package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryToken(ctx context.Context, cmd string) (catch bool, retMsg string) {
	// 权限校验
	if !service.User().CouldOpToken(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /token add <>
			catch, retMsg = tryTokenAdd(ctx, next[2])
		case "rm":
			// /token rm <name>
			retMsg = service.Token().RemoveTokenReturnRes(ctx, next[2])
			catch = true
		case "chown":
			// /token chown <>
			catch, retMsg = tryTokenChown(ctx, next[2])
		case "bind":
			// /token bind <>
			catch, retMsg = tryTokenBind(ctx, next[2])
		case "unbind":
			// /token unbind <name>
			retMsg = service.Token().UnbindTokenBotId(ctx, next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// /token query
			retMsg = service.Token().QueryTokenReturnRes(ctx)
			catch = true
		}
	}
	return
}

func tryTokenAdd(ctx context.Context, cmd string) (catch bool, retMsg string) {
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token add <name> <token>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Token().AddNewTokenReturnRes(ctx, dv[1], dv[2])
	catch = true
	return
}

func tryTokenChown(ctx context.Context, cmd string) (catch bool, retMsg string) {
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token chown <owner_id> <name>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Token().ChangeTokenOwnerReturnRes(ctx, dv[2], dv[1])
	catch = true
	return
}

func tryTokenBind(ctx context.Context, cmd string) (catch bool, retMsg string) {
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token bind <bot_id> <name>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Token().BindTokenBotId(ctx, dv[2], dv[1])
	catch = true
	return
}
