package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func trySys(ctx context.Context, args []string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.trySys")
	defer span.End()

	switch {
	case len(args) > 1:
		switch args[0] {
		case "forward":
			// /sys forward <>
			caught, retMsg = trySysForward(ctx, args[1:])
		case "check":
			// /sys check <>
			caught, retMsg = trySysCheck(ctx, args[1:])
		case "query":
			// /sys query <user_id>
			retMsg = service.User().QueryUserReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "grant":
			// /sys grant <>
			caught, retMsg = trySysGrant(ctx, args[1:])
		case "revoke":
			// /sys revoke <>
			caught, retMsg = trySysRevoke(ctx, args[1:])
		case "trust":
			// /sys trust <user_id>
			retMsg = service.User().SystemTrustUserReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "distrust":
			// /sys distrust <user_id>
			retMsg = service.User().SystemDistrustUserReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		}
	}
	return
}

func trySysForward(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "join":
			if len(args) < 3 {
				break
			}
			switch args[1] {
			case "user":
				// /sys forward join user <user_id>
				retMsg = service.Namespace().AddForwardingMatchUserIdReturnRes(ctx, args[2])
				caught = true
			case "group":
				// /sys forward join group <group_id>
				retMsg = service.Namespace().AddForwardingMatchGroupIdReturnRes(ctx, args[2])
				caught = true
			}
		case "leave":
			if len(args) < 3 {
				break
			}
			switch args[1] {
			case "user":
				// /sys forward leave user <user_id>
				retMsg = service.Namespace().RemoveForwardingMatchUserIdReturnRes(ctx, args[2])
				caught = true
			case "group":
				// /sys forward leave group <group_id>
				retMsg = service.Namespace().RemoveForwardingMatchGroupIdReturnRes(ctx, args[2])
				caught = true
			}
		case "reset":
			switch args[1] {
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
			switch {
			case len(args) > 3:
				// /sys forward add <alias> <url> <key>
				retMsg = service.Namespace().AddForwardingToReturnRes(ctx, args[1], args[2], args[3])
				caught = true
			case len(args) == 3:
				// /sys forward add <alias> <url>
				retMsg = service.Namespace().AddForwardingToReturnRes(ctx, args[1], args[2], "")
				caught = true
			}
		case "rm":
			// /sys forward rm <alias>
			retMsg = service.Namespace().RemoveForwardingToReturnRes(ctx, args[1])
			caught = true
		}
	}
	return
}

func trySysCheck(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "group":
			// /sys check group
			retMsg = service.Group().CheckExistReturnRes(ctx)
			caught = true
		}
	}
	return
}

func trySysGrant(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "raw":
			// /sys grant raw <user_id>
			retMsg = service.User().GrantGetRawMsgReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "recall":
			// /sys grant recall <user_id>
			retMsg = service.User().GrantRecallReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "namespace":
			// /sys grant namespace <user_id>
			retMsg = service.User().GrantOpNamespaceReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "token":
			// /sys grant token <user_id>
			retMsg = service.User().GrantOpTokenReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		}
	}
	return
}

func trySysRevoke(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "raw":
			// /sys revoke raw <user_id>
			retMsg = service.User().RevokeGetRawMsgReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "recall":
			// /sys revoke recall <user_id>
			retMsg = service.User().RevokeRecallReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "namespace":
			// /sys revoke namespace <user_id>
			retMsg = service.User().RevokeOpNamespaceReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "token":
			// /sys revoke token <user_id>
			retMsg = service.User().RevokeOpTokenReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		}
	}
	return
}
