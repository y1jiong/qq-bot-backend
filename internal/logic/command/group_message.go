package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryGroupMessage(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "enable":
			// /group message enable <>
			catch = tryGroupMessageEnable(ctx, next[2])
		case "disable":
			// /group message disable <>
			catch = tryGroupMessageDisable(ctx, next[2])
		case "set":
			// /group message set <>
			catch = tryGroupMessageSet(ctx, next[2])
		case "rm":
			// /group message rm <>
			catch = tryGroupMessageRemove(ctx, next[2])
		}
	}
	return
}

func tryGroupMessageEnable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "anti-recall":
			// /group message enable anti-recall
			service.Group().SetAntiRecallReturnRes(ctx, service.Bot().GetGroupId(ctx), true)
			catch = true
		}
	}
	return
}

func tryGroupMessageDisable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "anti-recall":
			// /group message disable anti-recall
			service.Group().SetAntiRecallReturnRes(ctx, service.Bot().GetGroupId(ctx), false)
			catch = true
		}
	}
	return
}

func tryGroupMessageSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "notification":
			// /group message set notification <group_id>
			service.Group().SetMessageNotificationReturnRes(ctx, service.Bot().GetGroupId(ctx), gconv.Int64(next[2]))
			catch = true
		}
	}
	return
}

func tryGroupMessageRemove(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "notification":
			// /group message rm notification
			service.Group().RemoveMessageNotificationReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}
