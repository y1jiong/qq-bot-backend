package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func tryToken(ctx context.Context, cmd string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.token")
	defer span.End()

	switch {
	case nextBranchRe.MatchString(cmd):
		// 权限校验
		if !service.User().CanOpToken(ctx, service.Bot().GetUserId(ctx)) {
			return
		}
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /token add <>
			caught, retMsg = tryTokenAdd(ctx, next[2])
		case "rm":
			// /token rm <name>
			retMsg = service.Token().RemoveTokenReturnRes(ctx, next[2])
			caught = true
		case "chown":
			// /token chown <>
			caught, retMsg = tryTokenChown(ctx, next[2])
		case "bind":
			// /token bind <>
			caught, retMsg = tryTokenBind(ctx, next[2])
		case "unbind":
			// /token unbind <name>
			retMsg = service.Token().UnbindTokenBotId(ctx, next[2])
			caught = true
		case "query":
			// /token query <name>
			retMsg = service.Token().QueryTokenReturnRes(ctx, next[2])
			caught = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// /token query
			retMsg = service.Token().QueryOwnTokenReturnRes(ctx)
			caught = true
		}
	}
	return
}

func tryTokenAdd(ctx context.Context, cmd string) (caught bool, retMsg string) {
	if !dualValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token add <name> <token>
	dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Token().AddNewTokenReturnRes(ctx, dv[1], dv[2])
	caught = true
	return
}

func tryTokenChown(ctx context.Context, cmd string) (caught bool, retMsg string) {
	if !dualValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token chown <owner_id> <name>
	dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Token().ChangeTokenOwnerReturnRes(ctx, dv[1], dv[2])
	caught = true
	return
}

func tryTokenBind(ctx context.Context, cmd string) (caught bool, retMsg string) {
	if !dualValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token bind <bot_id> <name>
	dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Token().BindTokenBotId(ctx, dv[1], dv[2])
	caught = true
	return
}
