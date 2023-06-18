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
		case "len":
			// /list len <list_name>
			service.List().QueryListLenReturnRes(ctx, next[2])
			catch = true
		case "query":
			// /list query <>
			catch = tryListQuery(ctx, next[2])
		case "leave":
			// /list leave <>
			catch = tryListLeave(ctx, next[2])
		case "export":
			// /list export <list_name>
			service.List().ExportListReturnRes(ctx, next[2])
			catch = true
		case "append":
			// /list append <>
			catch = tryListAppend(ctx, next[2])
		case "set":
			// /list set <>
			catch = tryListSet(ctx, next[2])
		case "reset":
			// /list reset <list_name>
			service.List().ResetListDataReturnRes(ctx, next[2])
			catch = true
		case "add":
			if !doubleValueCmdEndRe.MatchString(next[2]) {
				break
			}
			// /list add <list_name> <namespace>
			dv := doubleValueCmdEndRe.FindStringSubmatch(next[2])
			service.List().AddListReturnRes(ctx, dv[1], dv[2])
			catch = true
		case "rm":
			// /list rm <list_name>
			service.List().RemoveListReturnRes(ctx, next[2])
			catch = true
		}
	}
	return
}

func tryListJoin(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch {
		case nextBranchRe.MatchString(next[2]):
			// /list join <list_name> <key> [value]
			ne := nextBranchRe.FindStringSubmatch(next[2])
			service.List().AddListDataReturnRes(ctx, next[1], ne[1], ne[2])
			catch = true
		case endBranchRe.MatchString(next[2]):
			// /list join <list_name> <key>
			service.List().AddListDataReturnRes(ctx, next[1], next[2])
			catch = true
		}
	}
	return
}

func tryListQuery(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		// /list query <list_name> [key]
		next := nextBranchRe.FindStringSubmatch(cmd)
		service.List().QueryListReturnRes(ctx, next[1], next[2])
		catch = true
	case endBranchRe.MatchString(cmd):
		// /list query <list_name>
		service.List().QueryListReturnRes(ctx, cmd)
		catch = true
	}
	return
}

func tryListAppend(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		// /list append <list_name> <json>
		next := nextBranchRe.FindStringSubmatch(cmd)
		service.List().AppendListDataReturnRes(ctx, next[1], next[2])
		catch = true
	}
	return
}

func tryListSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		if endBranchRe.MatchString(next[2]) {
			// /list set <list_name> <json>
			service.List().SetListDataReturnRes(ctx, next[1], next[2])
			catch = true
		}
	}
	return
}

func tryListLeave(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		// /list leave <list_name> <key>
		service.List().RemoveListDataReturnRes(ctx, next[1], next[2])
		catch = true
	}
	return
}
