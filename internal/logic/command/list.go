package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryList(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "join":
			// /list join <>
			catch = tryListJoin(ctx, next[2])
		case "query":
			// /list query <list_name>
			service.List().QueryList(ctx, next[2])
			catch = true
		case "set":
			// /list set <>
			catch = tryListSet(ctx, next[2])
		case "leave":
			// /list leave <>
			catch = tryListLeave(ctx, next[2])
		case "reset":
			// /list reset <list_name>
			service.List().ResetListData(ctx, next[2])
			catch = true
		case "add":
			if !doubleValueCmdEndRe.MatchString(next[2]) {
				break
			}
			// /list add <list_name> <namespace>
			dv := doubleValueCmdEndRe.FindStringSubmatch(next[2])
			service.List().AddList(ctx, dv[1], dv[2])
			catch = true
		case "rm":
			// /list rm <list_name>
			service.List().RemoveList(ctx, next[2])
			catch = true
		}
	}
	return
}

func tryListJoin(ctx context.Context, cmd string) (catch bool) {
	if nextBranchRe.MatchString(cmd) {
		next := nextBranchRe.FindStringSubmatch(cmd)
		if doubleValueCmdEndRe.MatchString(next[2]) {
			// /list join <list_name> <key> [value]
			dv := doubleValueCmdEndRe.FindStringSubmatch(next[2])
			service.List().AddListData(ctx, next[1], dv[1], dv[2])
			catch = true
		} else {
			// /list join <list_name> <key>
			service.List().AddListData(ctx, next[1], next[2])
			catch = true
		}
	}
	return
}

func tryListSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		// /list set <list_name> <json>
		service.List().SetListData(ctx, next[1], next[2])
		catch = true
	}
	return
}

func tryListLeave(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		// /list leave <list_name> <key>
		service.List().RemoveListData(ctx, next[1], next[2])
		catch = true
	}
	return
}
