package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryGroup(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "approval":
			// /group approval <>
			catch, retMsg = tryGroupApproval(ctx, next[2])
		case "keyword":
			// /group keyword <>
			catch, retMsg = tryGroupKeyword(ctx, next[2])
		case "card":
			// /group card <>
			catch, retMsg = tryGroupCard(ctx, next[2])
		case "kick":
			// /group kick <list_name>
			retMsg = service.Group().KickFromListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		case "keep":
			// /group keep <list_name>
			retMsg = service.Group().KeepFromListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		case "log":
			// /group log <>
			catch, retMsg = tryGroupLog(ctx, next[2])
		case "export":
			// /group export <>
			catch, retMsg = tryGroupExport(ctx, next[2])
		case "message":
			// /group message <>
			catch, retMsg = tryGroupMessage(ctx, next[2])
		case "bind":
			// /group bind <namespace>
			retMsg = service.Group().BindNamespaceReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		case "clone":
			// /group clone <group_id>
			retMsg = service.Group().CloneReturnRes(ctx, service.Bot().GetGroupId(ctx), gconv.Int64(next[2]))
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// /group query
			retMsg = service.Group().QueryGroupReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		case "unbind":
			// /group unbind
			retMsg = service.Group().UnbindReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}
