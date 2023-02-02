package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func trySys(ctx context.Context, cmd string) (catch bool) {
	// 权限校验
	if !service.User().IsSystemTrustUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "grant":
			// /sys grant <>
			catch = trySysGrant(ctx, next[2])
		case "revoke":
			// /sys revoke <>
			catch = trySysRevoke(ctx, next[2])
		case "trust":
			// /sys trust <user_id>
			service.User().SystemTrustUser(ctx, gconv.Int64(next[2]))
			catch = true
		case "distrust":
			// /sys distrust <user_id>
			service.User().SystemDistrustUser(ctx, gconv.Int64(next[2]))
			catch = true
		}
	}
	return
}

func trySysGrant(ctx context.Context, cmd string) (catch bool) {
	switch {
	case doubleValueCmdEndRe.MatchString(cmd):
		dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
		switch dv[2] {
		case "raw":
			// /sys grant <user_id> raw
			service.User().GrantGetRawMsg(ctx, gconv.Int64(dv[1]))
			catch = true
		case "namespace":
			// /sys grant <user_id> namespace
			service.User().GrantOpNamespace(ctx, gconv.Int64(dv[1]))
			catch = true
		}
	}
	return
}

func trySysRevoke(ctx context.Context, cmd string) (catch bool) {
	switch {
	case doubleValueCmdEndRe.MatchString(cmd):
		dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
		switch dv[2] {
		case "raw":
			// /sys revoke <user_id> raw
			service.User().RevokeGetRawMsg(ctx, gconv.Int64(dv[1]))
			catch = true
		case "namespace":
			// /sys revoke <user_id> namespace
			service.User().RevokeOpNamespace(ctx, gconv.Int64(dv[1]))
			catch = true
		}
	}
	return
}
