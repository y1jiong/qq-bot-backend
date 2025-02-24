package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupLog(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "set":
			// /group log set <>
			caught, retMsg = tryGroupLogSet(ctx, args[1:])
		case "rm":
			// /group log rm <>
			caught, retMsg = tryGroupLogRemove(ctx, args[1:])
		}
	}
	return
}

func tryGroupLogSet(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "approval":
			// /group log set approval <list_name>
			retMsg = service.Group().SetLogApprovalListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		case "leave":
			// /group log set leave <list_name>
			retMsg = service.Group().SetLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		}
	}
	return
}

func tryGroupLogRemove(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "approval":
			// /group log rm approval
			retMsg = service.Group().RemoveLogApprovalListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "leave":
			// /group log rm leave
			retMsg = service.Group().RemoveLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		}
	}
	return
}
