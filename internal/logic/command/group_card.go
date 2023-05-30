package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupCard(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "check":
			// /group card check <>
			catch = tryGroupCardCheckout(ctx, next[2])
		case "set":
			// /group card set <>
			catch = tryGroupCardSet(ctx, next[2])
		case "rm":
			// /group card rm <>
			catch = tryGroupCardRemove(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "lock":
			// /group card lock
			service.Group().LockCardWithRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		case "unlock":
			// /group card unlock
			service.Group().UnlockCardWithRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}

func tryGroupCardSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "auto_set":
			// /group card set auto_set <list_name>
			service.Group().SetAutoSetListWithRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	}
	return
}

func tryGroupCardRemove(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "auto_set":
			// /group card rm auto_set
			service.Group().RemoveAutoSetListWithRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}

func tryGroupCardCheckout(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		// 提取 listName
		listName := next[1]
		// 以防命令不完整
		if !nextBranchRe.MatchString(next[2]) {
			return
		}
		next = nextBranchRe.FindStringSubmatch(next[2])
		switch next[1] {
		case "with":
			// /group card check <list_name> with <regexp>
			service.Group().CheckCardWithRegexpWithRes(ctx, service.Bot().GetGroupId(ctx), listName, next[2])
			catch = true
		case "by":
			// /group card check <to_list_name> by <from_list_name>
			service.Group().CheckCardByListWithRes(ctx, service.Bot().GetGroupId(ctx), listName, next[2])
			catch = true
		}
	}
	return
}
