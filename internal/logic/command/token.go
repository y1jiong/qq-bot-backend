package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func tryToken(ctx context.Context, args []string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryToken")
	defer span.End()

	switch {
	case len(args) > 1:
		// 权限校验
		if !service.User().CanOpToken(ctx, service.Bot().GetUserId(ctx)) {
			return
		}
		switch args[0] {
		case "add":
			// /token add <>
			caught, retMsg = tryTokenAdd(ctx, args[1:])
		case "rm":
			// /token rm <name>
			retMsg = service.Token().RemoveTokenReturnRes(ctx, args[1])
			caught = true
		case "chown":
			// /token chown <>
			caught, retMsg = tryTokenChown(ctx, args[1:])
		case "bind":
			// /token bind <>
			caught, retMsg = tryTokenBind(ctx, args[1:])
		case "unbind":
			// /token unbind <name>
			retMsg = service.Token().UnbindTokenBotId(ctx, args[1])
			caught = true
		case "query":
			// /token query <name>
			retMsg = service.Token().QueryTokenReturnRes(ctx, args[1])
			caught = true
		}
	case len(args) == 1:
		switch args[0] {
		case "query":
			// /token query
			retMsg = service.Token().QueryOwnTokenReturnRes(ctx)
			caught = true
		}
	}
	return
}

func tryTokenAdd(ctx context.Context, args []string) (caught bool, retMsg string) {
	if len(args) < 2 {
		return
	}
	// /token add <name> <token>
	retMsg = service.Token().AddNewTokenReturnRes(ctx, args[0], args[1])
	caught = true
	return
}

func tryTokenChown(ctx context.Context, args []string) (caught bool, retMsg string) {
	if len(args) < 2 {
		return
	}
	// /token chown <owner_id> <name>
	retMsg = service.Token().ChangeTokenOwnerReturnRes(ctx, args[0], args[1])
	caught = true
	return
}

func tryTokenBind(ctx context.Context, args []string) (caught bool, retMsg string) {
	if len(args) < 2 {
		return
	}
	// /token bind <bot_id> <name>
	retMsg = service.Token().BindTokenBotId(ctx, args[0], args[1])
	caught = true
	return
}
