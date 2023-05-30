package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryGroupApproval(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /group approval add <>
			catch = tryGroupApprovalAdd(ctx, next[2])
		case "set":
			// /group approval set <>
			catch = tryGroupApprovalSet(ctx, next[2])
		case "enable":
			// /group approval enable <>
			catch = tryGroupApprovalEnable(ctx, next[2])
		case "rm":
			// /group approval rm <>
			catch = tryGroupApprovalRemove(ctx, next[2])
		case "disable":
			// /group approval disable <>
			catch = tryGroupApprovalDisable(ctx, next[2])
		}
	}
	return
}

func tryGroupApprovalSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.RegexpCmd:
			// /group approval set regexp <regexp>
			service.Group().AddApprovalProcessWithRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalAdd(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.WhitelistCmd, consts.BlacklistCmd:
			// /group approval add whitelist <list_name>
			// /group approval add blacklist <list_name>
			service.Group().AddApprovalProcessWithRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalEnable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.WhitelistCmd, consts.BlacklistCmd, consts.RegexpCmd, consts.AutoPassCmd, consts.McCmd:
			// /group approval enable whitelist
			// /group approval enable blacklist
			// /group approval enable regexp
			// /group approval enable autopass
			// /group approval enable mc
			service.Group().AddApprovalProcessWithRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}

func tryGroupApprovalRemove(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.WhitelistCmd, consts.BlacklistCmd:
			// /group approval rm whitelist <list_name>
			// /group approval rm blacklist <list_name>
			service.Group().RemoveApprovalProcessWithRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalDisable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.WhitelistCmd, consts.BlacklistCmd, consts.RegexpCmd, consts.AutoPassCmd, consts.McCmd:
			// /group approval disable whitelist
			// /group approval disable blacklist
			// /group approval disable regexp
			// /group approval disable autopass
			// /group approval disable mc
			service.Group().RemoveApprovalProcessWithRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}
