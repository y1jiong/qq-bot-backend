package command

import (
	"context"
	"qq-bot-backend/internal/service"

	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
)

func tryUser(ctx context.Context, cmd string) (caught catch, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.user")
	defer span.End()

	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "join":
			// /user join <>
			caught, retMsg = tryUserJoin(ctx, next[2])
		case "leave":
			// /user leave <>
			caught, retMsg = tryUserLeave(ctx, next[2])
		}
	}
	return
}

func tryUserJoin(ctx context.Context, cmd string) (caught catch, retMsg string) {
	if !dualValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /user join <namespace> <user_id>
	dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Namespace().AddNamespaceAdminReturnRes(ctx, dv[1], gconv.Int64(dv[2]))
	caught = caughtNeedOkay
	return
}

func tryUserLeave(ctx context.Context, cmd string) (caught catch, retMsg string) {
	if !dualValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /user leave <namespace> <user_id>
	dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Namespace().RemoveNamespaceAdminReturnRes(ctx, dv[1], gconv.Int64(dv[2]))
	caught = caughtNeedOkay
	return
}
