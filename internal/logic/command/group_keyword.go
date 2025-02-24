package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryGroupKeyword(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "add":
			// /group keyword add <>
			caught, retMsg = tryGroupKeywordAdd(ctx, args[1:])
		case "enable":
			// /group keyword enable <>
			caught, retMsg = tryGroupKeywordEnable(ctx, args[1:])
		case "rm":
			// /group keyword rm <>
			caught, retMsg = tryGroupKeywordRemove(ctx, args[1:])
		case "disable":
			// /group keyword disable <>
			caught, retMsg = tryGroupKeywordDisable(ctx, args[1:])
		}
	}
	return
}

func tryGroupKeywordAdd(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case consts.ReplyCmd, consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword add blacklist <list_name>
			// /group keyword add whitelist <list_name>
			// /group keyword add reply <list_name>
			retMsg = service.Group().AddKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), args[0], args[1])
			caught = true
		}
	}
	return
}

func tryGroupKeywordEnable(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword enable blacklist
			// /group keyword enable whitelist
			retMsg = service.Group().AddKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), args[0])
			caught = true
		}
	}
	return
}

func tryGroupKeywordRemove(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case consts.ReplyCmd, consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword rm blacklist <list_name>
			// /group keyword rm whitelist <list_name>
			// /group keyword rm reply <list_name>
			retMsg = service.Group().RemoveKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), args[0], args[1])
			caught = true
		}
	}
	return
}

func tryGroupKeywordDisable(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword disable blacklist
			// /group keyword disable whitelist
			retMsg = service.Group().RemoveKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), args[0])
			caught = true
		}
	}
	return
}
