package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryUser(ctx context.Context, args []string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryUser")
	defer span.End()

	switch {
	case len(args) > 1:
		switch args[0] {
		case "join":
			// /user join <>
			caught, retMsg = tryUserJoin(ctx, args[1:])
		case "leave":
			// /user leave <>
			caught, retMsg = tryUserLeave(ctx, args[1:])
		}
	}
	return
}

func tryUserJoin(ctx context.Context, args []string) (caught bool, retMsg string) {
	if len(args) < 2 {
		return
	}
	// /user join <namespace> <user_id>
	retMsg = service.Namespace().AddNamespaceAdminReturnRes(ctx, args[0], gconv.Int64(args[1]))
	caught = true
	return
}

func tryUserLeave(ctx context.Context, args []string) (caught bool, retMsg string) {
	if len(args) < 2 {
		return
	}
	// /user leave <namespace> <user_id>
	retMsg = service.Namespace().RemoveNamespaceAdminReturnRes(ctx, args[0], gconv.Int64(args[1]))
	caught = true
	return
}
