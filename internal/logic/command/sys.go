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
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "forward":
			// /sys forward <>
			catch, retMsg = trySysForward(ctx, next[2])
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

func trySysForward(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "join":
			if !dualValueCmdEndRe.MatchString(next[2]) {
				break
			}
			dv := dualValueCmdEndRe.FindStringSubmatch(next[2])
			switch dv[1] {
			case "user":
				// /sys forward join user <user_id>
				retMsg = service.Namespace().AddForwardingMatchUserIdReturnRes(ctx, dv[2])
				catch = true
			case "group":
				// /sys forward join group <group_id>
				retMsg = service.Namespace().AddForwardingMatchGroupIdReturnRes(ctx, dv[2])
				catch = true
			}
		case "leave":
			if !dualValueCmdEndRe.MatchString(next[2]) {
				break
			}
			dv := dualValueCmdEndRe.FindStringSubmatch(next[2])
			switch dv[1] {
			case "user":
				// /sys forward leave user <user_id>
				retMsg = service.Namespace().RemoveForwardingMatchUserIdReturnRes(ctx, dv[2])
				catch = true
			case "group":
				// /sys forward leave group <group_id>
				retMsg = service.Namespace().RemoveForwardingMatchGroupIdReturnRes(ctx, dv[2])
				catch = true
			}
		case "reset":
			if !endBranchRe.MatchString(next[2]) {
				break
			}
			switch next[2] {
			case "user":
				// /sys forward reset user
				retMsg = service.Namespace().ResetForwardingMatchUserIdReturnRes(ctx)
				catch = true
			case "group":
				// /sys forward reset group
				retMsg = service.Namespace().ResetForwardingMatchGroupIdReturnRes(ctx)
				catch = true
			}
		case "add":
			if !nextBranchRe.MatchString(next[2]) {
				break
			}
			ne := nextBranchRe.FindStringSubmatch(next[2])
			if dualValueCmdEndRe.MatchString(ne[2]) {
				dv := dualValueCmdEndRe.FindStringSubmatch(ne[2])
				args := make([]string, 3)
				args[0] = ne[1]
				args[1] = dv[1]
				args[2] = dv[2]
				// /sys forward add <alias> <url> <authorization>
				retMsg = service.Namespace().AddForwardingToReturnRes(ctx, args[0], args[1], args[2])
				catch = true
			}
			if endBranchRe.MatchString(ne[2]) {
				// /sys forward add <alias> <url>
				retMsg = service.Namespace().AddForwardingToReturnRes(ctx, ne[1], ne[2], "")
				catch = true
			}
		case "rm":
			if !endBranchRe.MatchString(next[2]) {
				break
			}
			// /sys forward rm <alias>
			retMsg = service.Namespace().RemoveForwardingToReturnRes(ctx, next[2])
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
	case dualValueCmdEndRe.MatchString(cmd):
		dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
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
	case dualValueCmdEndRe.MatchString(cmd):
		dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
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
