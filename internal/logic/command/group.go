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
		case "log":
			// /group log <>
			catch = tryGroupLog(ctx, next[2])
		case "export":
			// /group export <>
			catch = tryGroupExport(ctx, next[2])
		case "bind":
			// /group bind <namespace>
			service.Group().BindNamespace(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
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
		case consts.WhitelistCmd, consts.BlacklistCmd:
			// /group approval add whitelist <list_name>
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
		switch cmd {
		case consts.WhitelistCmd, consts.BlacklistCmd, consts.RegexpCmd, consts.AutoPassCmd, consts.McCmd:
			// /group approval enable whitelist
			// /group approval enable blacklist
			// /group approval enable regexp
			// /group approval enable autopass
			// /group approval enable mc
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), cmd)
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
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
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
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), cmd)
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
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
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
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
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
			service.Group().AddKeywordProcess(ctx, service.Bot().GetGroupId(ctx), cmd)
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
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case consts.ReplyCmd:
			// /group keyword rm reply
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), cmd)
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
			service.Group().RemoveKeywordProcess(ctx, service.Bot().GetGroupId(ctx), cmd)
			catch = true
		}
	}
	return
}

func tryGroupLog(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "leave":
			// /group log leave <>
			catch = tryGroupLogLeave(ctx, next[2])
		}
	}
	return
}

func tryGroupLogLeave(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "set":
			// /group log leave set <list_name>
			service.Group().SetLogLeaveList(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "rm":
			// /group log leave rm
			service.Group().RemoveLogLeaveList(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}

func tryGroupExport(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "member":
			// /group export member <list_name>
			service.Group().ExportGroupMemberList(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	}
	return
}
