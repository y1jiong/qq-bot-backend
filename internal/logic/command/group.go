package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryGroup(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "approval":
			// /group approval <>
			catch = tryGroupApproval(ctx, next[2])
		case "keyword":
			// /group keyword <>
			catch = tryGroupKeyword(ctx, next[2])
		case "bind":
			// /group bind <namespace>
			service.Group().BindNamespace(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch endBranchRe.FindString(cmd) {
		case "query":
			// /group query
			service.Group().QueryGroup(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		case "unbind":
			// /group unbind
			service.Group().Unbind(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}

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
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
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
		case consts.WhitelistCmd:
			// /group approval add whitelist <list_name>
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case consts.BlacklistCmd:
			// /group approval add blacklist <list_name>
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalEnable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		end := endBranchRe.FindString(cmd)
		switch end {
		case consts.WhitelistCmd:
			// /group approval enable whitelist
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.BlacklistCmd:
			// /group approval enable blacklist
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.RegexpCmd:
			// /group approval enable regexp
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.McCmd:
			// /group approval enable mc
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
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
		case consts.WhitelistCmd:
			// /group approval rm whitelist <list_name>
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case consts.BlacklistCmd:
			// /group approval rm blacklist <list_name>
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupApprovalDisable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		end := endBranchRe.FindString(cmd)
		switch end {
		case consts.WhitelistCmd:
			// /group approval disable whitelist
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.BlacklistCmd:
			// /group approval disable blacklist
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.RegexpCmd:
			// /group approval disable regexp
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.McCmd:
			// /group approval disable mc
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		}
	}
	return
}

func tryGroupKeyword(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /group keyword add <>
			catch = tryGroupKeywordAdd(ctx, next[2])
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
		case consts.BlacklistCmd:
			// /group keyword add blacklist <list_name>
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case consts.WhitelistCmd:
			// /group keyword add whitelist <list_name>
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupKeywordEnable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		end := endBranchRe.FindString(cmd)
		switch end {
		case consts.BlacklistCmd:
			// /group keyword enable blacklist
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.WhitelistCmd:
			// /group keyword enable whitelist
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), end)
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
		case consts.BlacklistCmd:
			// /group keyword rm blacklist <list_name>
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case consts.WhitelistCmd:
			// /group keyword rm whitelist <list_name>
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	}
	return
}

func tryGroupKeywordDisable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		end := endBranchRe.FindString(cmd)
		switch end {
		case consts.WhitelistCmd:
			// /group keyword disable whitelist
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case consts.BlacklistCmd:
			// /group keyword disable blacklist
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		}
	}
	return
}
