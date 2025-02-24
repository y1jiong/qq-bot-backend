package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryGroupMessage(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "enable":
			// /group message enable <>
			caught, retMsg = tryGroupMessageEnable(ctx, args[1:])
		case "disable":
			// /group message disable <>
			caught, retMsg = tryGroupMessageDisable(ctx, args[1:])
		case "set":
			// /group message set <>
			caught, retMsg = tryGroupMessageSet(ctx, args[1:])
		case "rm":
			// /group message rm <>
			caught, retMsg = tryGroupMessageRemove(ctx, args[1:])
		}
	}
	return
}

func tryGroupMessageEnable(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "anti-recall":
			// /group message enable anti-recall
			retMsg = service.Group().SetAntiRecallReturnRes(ctx, service.Bot().GetGroupId(ctx), true)
			caught = true
		}
	}
	return
}

func tryGroupMessageDisable(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "anti-recall":
			// /group message disable anti-recall
			retMsg = service.Group().SetAntiRecallReturnRes(ctx, service.Bot().GetGroupId(ctx), false)
			caught = true
		}
	}
	return
}

func tryGroupMessageSet(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "notification":
			// /group message set notification <group_id>
			retMsg = service.Group().SetMessageNotificationReturnRes(ctx, service.Bot().GetGroupId(ctx), gconv.Int64(args[1]))
			caught = true
		}
	case len(args) == 1:
		switch args[0] {
		case "only-anti-recall-member":
			// /group message set only-anti-recall-member
			retMsg = service.Group().SetOnlyAntiRecallMemberReturnRes(ctx, service.Bot().GetGroupId(ctx), true)
			caught = true
		}
	}
	return
}

func tryGroupMessageRemove(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "notification":
			// /group message rm notification
			retMsg = service.Group().RemoveMessageNotificationReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "only-anti-recall-member":
			// /group message rm only-anti-recall-member
			retMsg = service.Group().SetOnlyAntiRecallMemberReturnRes(ctx, service.Bot().GetGroupId(ctx), false)
			caught = true
		}
	}
	return
}
