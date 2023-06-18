package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryGroupKeyword(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /group keyword add <>
			catch = tryGroupKeywordAdd(ctx, next[2])
		case "set":
			// /group keyword set <>
			catch = tryGroupKeywordSet(ctx, next[2])
		case "enable":
			// /group keyword enable <>
			catch = tryGroupKeywordEnable(ctx, next[2])
		case "rm":
			// /group keyword rm <>
			catch = tryGroupKeywordRemove(ctx, next[2])
		case "disable":
			// /group keyword disable <>
			catch = tryGroupKeywordDisable(ctx, next[2])
		}
	}
	return
}

func tryGroupKeywordAdd(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword add blacklist <list_name>
			// /group keyword add whitelist <list_name>
			service.Group().AddKeywordProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupKeywordSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.ReplyCmd:
			// /group keyword set reply <list_name>
			service.Group().AddKeywordProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupKeywordEnable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword enable blacklist
			// /group keyword enable whitelist
			service.Group().AddKeywordProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}

func tryGroupKeywordRemove(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword rm blacklist <list_name>
			// /group keyword rm whitelist <list_name>
			service.Group().RemoveKeywordProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.ReplyCmd:
			// /group keyword rm reply
			service.Group().RemoveKeywordProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}

func tryGroupKeywordDisable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.BlacklistCmd, consts.WhitelistCmd:
			// /group keyword disable blacklist
			// /group keyword disable whitelist
			service.Group().RemoveKeywordProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}
