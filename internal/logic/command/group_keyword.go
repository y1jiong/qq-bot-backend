package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryGroupKeyword(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /group keyword add <>
			caught, retMsg = tryGroupKeywordAdd(ctx, next[2])
		case "enable":
			// /group keyword enable <>
			caught, retMsg = tryGroupKeywordEnable(ctx, next[2])
		case "rm":
			// /group keyword rm <>
			caught, retMsg = tryGroupKeywordRemove(ctx, next[2])
		case "disable":
			// /group keyword disable <>
			caught, retMsg = tryGroupKeywordDisable(ctx, next[2])
		}
	}
	return
}

func tryGroupKeywordAdd(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.ReplyCmd, consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword add blacklist <list_name>
			// /group keyword add whitelist <list_name>
			// /group keyword add reply <list_name>
			retMsg = service.Group().AddKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			caught = true
		}
	}
	return
}

func tryGroupKeywordEnable(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword enable blacklist
			// /group keyword enable whitelist
			retMsg = service.Group().AddKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			caught = true
		}
	}
	return
}

func tryGroupKeywordRemove(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.ReplyCmd, consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword rm blacklist <list_name>
			// /group keyword rm whitelist <list_name>
			// /group keyword rm reply <list_name>
			retMsg = service.Group().RemoveKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			caught = true
		}
	}
	return
}

func tryGroupKeywordDisable(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword disable blacklist
			// /group keyword disable whitelist
			retMsg = service.Group().RemoveKeywordPolicyReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			caught = true
		}
	}
	return
}
