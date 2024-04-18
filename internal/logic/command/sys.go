package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func trySys(ctx context.Context, cmd string) (catch bool, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "check":
			// /sys check <>
			catch, retMsg = trySysCheck(ctx, next[2])
		case "query":
			// /sys query <user_id>
			retMsg = service.User().QueryUserReturnRes(ctx, gconv.Int64(next[2]))
			catch = true
		case "grant":
			// /sys grant <>
			catch, retMsg = trySysGrant(ctx, next[2])
		case "revoke":
			// /sys revoke <>
			catch, retMsg = trySysRevoke(ctx, next[2])
		case "trust":
			// /sys trust <user_id>
			retMsg = service.User().SystemTrustUserReturnRes(ctx, gconv.Int64(next[2]))
			catch = true
		case "distrust":
			// /sys distrust <user_id>
			retMsg = service.User().SystemDistrustUserReturnRes(ctx, gconv.Int64(next[2]))
			catch = true
		}
	}
	return
}

func trySysCheck(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "group":
			// /sys check group
			retMsg = service.Group().CheckExistReturnRes(ctx)
			catch = true
		}
	}
	return
}

func trySysGrant(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case doubleValueCmdEndRe.MatchString(cmd):
		dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
		switch dv[1] {
		case "raw":
			// /sys grant raw <user_id>
			retMsg = service.User().GrantGetRawMsgReturnRes(ctx, gconv.Int64(dv[2]))
			catch = true
		case "namespace":
			// /sys grant namespace <user_id>
			retMsg = service.User().GrantOpNamespaceReturnRes(ctx, gconv.Int64(dv[2]))
			catch = true
		case "token":
			// /sys grant token <user_id>
			retMsg = service.User().GrantOpTokenReturnRes(ctx, gconv.Int64(dv[2]))
			catch = true
		}
	}
	return
}

func trySysRevoke(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case doubleValueCmdEndRe.MatchString(cmd):
		dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
		switch dv[1] {
		case "raw":
			// /sys revoke raw <user_id>
			retMsg = service.User().RevokeGetRawMsgReturnRes(ctx, gconv.Int64(dv[2]))
			catch = true
		case "namespace":
			// /sys revoke namespace <user_id>
			retMsg = service.User().RevokeOpNamespaceReturnRes(ctx, gconv.Int64(dv[2]))
			catch = true
		case "token":
			// /sys revoke token <user_id>
			retMsg = service.User().RevokeOpTokenReturnRes(ctx, gconv.Int64(dv[2]))
			catch = true
		}
	}
	return
}
