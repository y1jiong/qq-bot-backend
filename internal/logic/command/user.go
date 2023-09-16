package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryUser(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "join":
			// /user join <>
			catch, retMsg = tryUserJoin(ctx, next[2])
		case "leave":
			// /user leave <>
			catch, retMsg = tryUserLeave(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
	}
	return
}

func tryUserJoin(ctx context.Context, cmd string) (catch bool, retMsg string) {
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /user join <namespace> <user_id>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Namespace().AddNamespaceAdminReturnRes(ctx, dv[1], gconv.Int64(dv[2]))
	catch = true
	return
}

func tryUserLeave(ctx context.Context, cmd string) (catch bool, retMsg string) {
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /user leave <namespace> <user_id>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	retMsg = service.Namespace().RemoveNamespaceAdminReturnRes(ctx, dv[1], gconv.Int64(dv[2]))
	catch = true
	return
}
