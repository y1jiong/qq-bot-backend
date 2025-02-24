package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupCard(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "check":
			// /group card check <>
			caught, retMsg = tryGroupCardCheck(ctx, args[1:])
		case "set":
			// /group card set <>
			caught, retMsg = tryGroupCardSet(ctx, args[1:])
		case "rm":
			// /group card rm <>
			caught, retMsg = tryGroupCardRemove(ctx, args[1:])
		}
	case len(args) == 1:
		switch args[0] {
		case "lock":
			// /group card lock
			retMsg = service.Group().LockCardReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "unlock":
			// /group card unlock
			retMsg = service.Group().UnlockCardReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		}
	}
	return
}

func tryGroupCardSet(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "auto-set":
			// /group card set auto-set <list_name>
			retMsg = service.Group().SetAutoSetListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		}
	}
	return
}

func tryGroupCardRemove(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "auto-set":
			// /group card rm auto-set
			retMsg = service.Group().RemoveAutoSetListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		}
	}
	return
}

func tryGroupCardCheck(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		// 以防命令不完整
		if len(args) < 3 {
			break
		}
		switch args[1] {
		case "with":
			// /group card check <list_name> with <regexp>
			retMsg = service.Group().CheckCardWithRegexpReturnRes(ctx, service.Bot().GetGroupId(ctx), args[0], args[2])
			caught = true
		case "by":
			// /group card check <to_list_name> by <from_list_name>
			retMsg = service.Group().CheckCardByListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[0], args[2])
			caught = true
		}
	}
	return
}
