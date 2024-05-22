package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryGroupApproval(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /group approval add <>
			catch, retMsg = tryGroupApprovalAdd(ctx, next[2])
		case "set":
			// /group approval set <>
			catch, retMsg = tryGroupApprovalSet(ctx, next[2])
		case "enable":
			// /group approval enable <>
			catch, retMsg = tryGroupApprovalEnable(ctx, next[2])
		case "rm":
			// /group approval rm <>
			catch, retMsg = tryGroupApprovalRemove(ctx, next[2])
		case "disable":
			// /group approval disable <>
			catch, retMsg = tryGroupApprovalDisable(ctx, next[2])
		}
	}
	return
}

func tryGroupApprovalSet(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.RegexpCmd, consts.NotificationCmd:
			// /group approval set regexp <regexp>
			// /group approval set notification <group_id>
			retMsg = service.Group().AddApprovalProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalAdd(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.WhitelistCmd, consts.BlacklistCmd:
			// /group approval add whitelist <list_name>
			// /group approval add blacklist <list_name>
			retMsg = service.Group().AddApprovalProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalEnable(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.WhitelistCmd, consts.BlacklistCmd, consts.RegexpCmd, consts.McCmd,
			consts.NotifyOnlyCmd, consts.AutoPassCmd, consts.AutoRejectCmd:
			// /group approval enable whitelist
			// /group approval enable blacklist
			// /group approval enable regexp
			// /group approval enable mc
			// /group approval enable notify-only
			// /group approval enable auto-pass
			// /group approval enable auto-reject
			retMsg = service.Group().AddApprovalProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}

func tryGroupApprovalRemove(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.WhitelistCmd, consts.BlacklistCmd:
			// /group approval rm whitelist <list_name>
			// /group approval rm blacklist <list_name>
			retMsg = service.Group().RemoveApprovalProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.NotificationCmd:
			// /group approval rm notification
			retMsg = service.Group().RemoveApprovalProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}

func tryGroupApprovalDisable(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.WhitelistCmd, consts.BlacklistCmd, consts.RegexpCmd, consts.McCmd,
			consts.NotifyOnlyCmd, consts.AutoPassCmd, consts.AutoRejectCmd:
			// /group approval disable whitelist
			// /group approval disable blacklist
			// /group approval disable regexp
			// /group approval disable mc
			// /group approval disable notify-only
			// /group approval disable auto-pass
			// /group approval disable auto-reject
			retMsg = service.Group().RemoveApprovalProcessReturnRes(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}
