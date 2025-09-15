package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupCard(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "check":
			// /group card check <>
			caught, retMsg = tryGroupCardCheckout(ctx, next[2])
		case "set":
			// /group card set <>
			caught, retMsg = tryGroupCardSet(ctx, next[2])
		case "rm":
			// /group card rm <>
			caught, retMsg = tryGroupCardRemove(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "lock":
			// /group card lock
			retMsg = service.Group().LockCardReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = caughtNeedOkay
		case "unlock":
			// /group card unlock
			retMsg = service.Group().UnlockCardReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = caughtNeedOkay
		}
	}
	return
}

func tryGroupCardSet(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "auto-set":
			// /group card set auto-set <list_name>
			retMsg = service.Group().SetAutoSetListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			caught = caughtNeedOkay
		}
	}
	return
}

func tryGroupCardRemove(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "auto-set":
			// /group card rm auto-set
			retMsg = service.Group().RemoveAutoSetListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = caughtNeedOkay
		}
	}
	return
}

func tryGroupCardCheckout(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		// 提取 listName
		listName := next[1]
		// 以防命令不完整
		if !nextBranchRe.MatchString(next[2]) {
			break
		}
		next = nextBranchRe.FindStringSubmatch(next[2])
		switch next[1] {
		case "with":
			// /group card check <list_name> with <regexp>
			retMsg = service.Group().CheckCardWithRegexpReturnRes(ctx, service.Bot().GetGroupId(ctx), listName, next[2])
			caught = caughtNeedOkay
		case "by":
			// /group card check <to_list_name> by <from_list_name>
			retMsg = service.Group().CheckCardByListReturnRes(ctx, service.Bot().GetGroupId(ctx), listName, next[2])
			caught = caughtNeedOkay
		}
	}
	return
}
