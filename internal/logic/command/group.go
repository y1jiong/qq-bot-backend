package command

import (
	"context"
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
		case "card":
			// /group card <>
			catch = tryGroupCard(ctx, next[2])
		case "kick":
			// /group kick <list_name>
			service.Group().KickFromListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		case "keep":
			// /group keep <list_name>
			service.Group().KeepFromListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		case "log":
			// /group log <>
			catch = tryGroupLog(ctx, next[2])
		case "export":
			// /group export <>
			catch = tryGroupExport(ctx, next[2])
		case "message":
			// /group message <>
			catch = tryGroupMessage(ctx, next[2])
		case "bind":
			// /group bind <namespace>
			service.Group().BindNamespaceReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// /group query
			service.Group().QueryGroupReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		case "unbind":
			// /group unbind
			service.Group().UnbindReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}
