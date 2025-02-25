package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func trySys(ctx context.Context, cmd string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.trySys")
	defer span.End()

	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "query":
			// /sys query <user_id>
			retMsg = service.User().QueryUserReturnRes(ctx, gconv.Int64(next[2]))
			caught = true
		case "grant":
			// /sys grant <>
			caught, retMsg = trySysGrant(ctx, next[2])
		case "revoke":
			// /sys revoke <>
			caught, retMsg = trySysRevoke(ctx, next[2])
		case "check":
			// /sys check <>
			caught, retMsg = trySysCheck(ctx, next[2])
		case "forward":
			// /sys forward <>
			caught, retMsg = trySysForward(ctx, next[2])
		case "trust":
			// /sys trust <user_id>
			retMsg = service.User().SystemTrustUserReturnRes(ctx, gconv.Int64(next[2]))
			caught = true
		case "distrust":
			// /sys distrust <user_id>
			retMsg = service.User().SystemDistrustUserReturnRes(ctx, gconv.Int64(next[2]))
			caught = true
		}
	}
	return
}

func trySysForward(ctx context.Context, cmd string) (caught bool, retMsg string) {
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
				caught = true
			case "group":
				// /sys forward join group <group_id>
				retMsg = service.Namespace().AddForwardingMatchGroupIdReturnRes(ctx, dv[2])
				caught = true
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
				caught = true
			case "group":
				// /sys forward leave group <group_id>
				retMsg = service.Namespace().RemoveForwardingMatchGroupIdReturnRes(ctx, dv[2])
				caught = true
			}
		case "reset":
			switch next[2] {
			case "user":
				// /sys forward reset user
				retMsg = service.Namespace().ResetForwardingMatchUserIdReturnRes(ctx)
				caught = true
			case "group":
				// /sys forward reset group
				retMsg = service.Namespace().ResetForwardingMatchGroupIdReturnRes(ctx)
				caught = true
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
				// /sys forward add <alias> <url> <key>
				retMsg = service.Namespace().AddForwardingToReturnRes(ctx, args[0], args[1], args[2])
				caught = true
			}
			if endBranchRe.MatchString(ne[2]) {
				// /sys forward add <alias> <url>
				retMsg = service.Namespace().AddForwardingToReturnRes(ctx, ne[1], ne[2], "")
				caught = true
			}
		case "rm":
			if !endBranchRe.MatchString(next[2]) {
				break
			}
			// /sys forward rm <alias>
			retMsg = service.Namespace().RemoveForwardingToReturnRes(ctx, next[2])
			caught = true
		}
	}
	return
}

func trySysCheck(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "group":
			// /sys check group
			retMsg = service.Group().CheckExistReturnRes(ctx)
			caught = true
		}
	}
	return
}

func trySysGrant(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case dualValueCmdEndRe.MatchString(cmd):
		dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
		switch dv[1] {
		case "raw":
			// /sys grant raw <user_id>
			retMsg = service.User().GrantGetRawMsgReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "recall":
			// /sys grant recall <user_id>
			retMsg = service.User().GrantRecallReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "crontab":
			// /sys grant crontab <user_id>
			retMsg = service.User().GrantOpCrontabReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "namespace":
			// /sys grant namespace <user_id>
			retMsg = service.User().GrantOpNamespaceReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "token":
			// /sys grant token <user_id>
			retMsg = service.User().GrantOpTokenReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		}
	}
	return
}

func trySysRevoke(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case dualValueCmdEndRe.MatchString(cmd):
		dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
		switch dv[1] {
		case "raw":
			// /sys revoke raw <user_id>
			retMsg = service.User().RevokeGetRawMsgReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "recall":
			// /sys revoke recall <user_id>
			retMsg = service.User().RevokeRecallReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "crontab":
			// /sys revoke crontab <user_id>
			retMsg = service.User().RevokeOpCrontabReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "namespace":
			// /sys revoke namespace <user_id>
			retMsg = service.User().RevokeOpNamespaceReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		case "token":
			// /sys revoke token <user_id>
			retMsg = service.User().RevokeOpTokenReturnRes(ctx, gconv.Int64(dv[2]))
			caught = true
		}
	}
	return
}
