package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryGroup(ctx context.Context, args []string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryGroup")
	defer span.End()

	switch {
	case len(args) > 1:
		switch args[0] {
		case "query":
			// /group query <group_id>
			retMsg = service.Group().QueryGroupReturnRes(ctx, gconv.Int64(args[1]))
			caught = true
		case "approval":
			// /group approval <>
			caught, retMsg = tryGroupApproval(ctx, args[1:])
		case "keyword":
			// /group keyword <>
			caught, retMsg = tryGroupKeyword(ctx, args[1:])
		case "message":
			// /group message <>
			caught, retMsg = tryGroupMessage(ctx, args[1:])
		case "card":
			// /group card <>
			caught, retMsg = tryGroupCard(ctx, args[1:])
		case "export":
			// /group export <>
			caught, retMsg = tryGroupExport(ctx, args[1:])
		case "log":
			// /group log <>
			caught, retMsg = tryGroupLog(ctx, args[1:])
		case "kick":
			// /group kick <list_name>
			retMsg = service.Group().KickFromListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		case "keep":
			// /group keep <list_name>
			retMsg = service.Group().KeepFromListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		case "clone":
			// /group clone <group_id>
			retMsg = service.Group().CloneReturnRes(ctx, service.Bot().GetGroupId(ctx), gconv.Int64(args[1]))
			caught = true
		case "bind":
			// /group bind <namespace>
			retMsg = service.Group().BindNamespaceReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		}
	case len(args) == 1:
		switch args[0] {
		case "query":
			// /group query
			retMsg = service.Group().QueryGroupReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "unbind":
			// /group unbind
			retMsg = service.Group().UnbindReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		}
	}
	return
}
