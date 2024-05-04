package command

import (
	"context"
	"qq-bot-backend/internal/service"
	"strings"
)

func tryList(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "join":
			// /list join <>
			catch, retMsg = tryListJoin(ctx, next[2])
		case "leave":
			// /list leave <>
			if !nextBranchRe.MatchString(next[2]) {
				break
			}
			next := nextBranchRe.FindStringSubmatch(next[2])
			// /list leave <list_name> <key>
			retMsg = service.List().RemoveListDataReturnRes(ctx, next[1], next[2])
			catch = true
		case "copy-key":
			// /list copy-key <>
			if !nextBranchRe.MatchString(next[2]) {
				break
			}
			next := nextBranchRe.FindStringSubmatch(next[2])
			if !dualValueCmdEndRe.MatchString(next[2]) {
				break
			}
			dv := dualValueCmdEndRe.FindStringSubmatch(next[2])
			// /list copy-key <list_name> <src_key> <dst_key>
			retMsg = service.List().CopyListKeyReturnRes(ctx, next[1], dv[1], dv[2])
			catch = true
		case "glance":
			// /list glance <list_name>
			retMsg = service.List().GlanceListDataReturnRes(ctx, next[2])
			catch = true
		case "query":
			// /list query <>
			catch, retMsg = tryListQuery(ctx, next[2])
		case "len":
			// /list len <list_name>
			retMsg = service.List().QueryListLenReturnRes(ctx, next[2])
			catch = true
		case "export":
			// /list export <list_name>
			retMsg = service.List().ExportListReturnRes(ctx, next[2])
			catch = true
		case "append":
			// /list append <>
			if !nextBranchRe.MatchString(next[2]) {
				break
			}
			next := nextBranchRe.FindStringSubmatch(next[2])
			// /list append <list_name> <json>
			retMsg = service.List().AppendListDataReturnRes(ctx, next[1], next[2])
			catch = true
		case "set":
			// /list set <>
			if !nextBranchRe.MatchString(next[2]) {
				break
			}
			next := nextBranchRe.FindStringSubmatch(next[2])
			// /list set <list_name> <json>
			retMsg = service.List().SetListDataReturnRes(ctx, next[1], next[2])
			catch = true
		case "reset":
			// /list reset <list_name>
			retMsg = service.List().ResetListDataReturnRes(ctx, next[2])
			catch = true
		case "add":
			if !dualValueCmdEndRe.MatchString(next[2]) {
				break
			}
			// /list add <list_name> <namespace>
			dv := dualValueCmdEndRe.FindStringSubmatch(next[2])
			retMsg = service.List().AddListReturnRes(ctx, dv[1], dv[2])
			catch = true
		case "rm":
			// /list rm <list_name>
			retMsg = service.List().RemoveListReturnRes(ctx, next[2])
			catch = true
		case "recover":
			// /list recover <list_name>
			retMsg = service.List().RecoverListReturnRes(ctx, next[2])
			catch = true
		case "op":
			// /list op <>
			catch, retMsg = tryListOp(ctx, next[2])
		}
	}
	return
}

func tryListJoin(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch {
		case nextBranchRe.MatchString(next[2]):
			// /list join <list_name> <key> [value]
			ne := nextBranchRe.FindStringSubmatch(next[2])
			retMsg = service.List().AddListDataReturnRes(ctx, next[1], ne[1], ne[2])
			catch = true
		case endBranchRe.MatchString(next[2]):
			// /list join <list_name> <key>
			retMsg = service.List().AddListDataReturnRes(ctx, next[1], next[2])
			catch = true
		}
	}
	return
}

func tryListQuery(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		// /list query <list_name> [key]
		next := nextBranchRe.FindStringSubmatch(cmd)
		retMsg = service.List().QueryListReturnRes(ctx, next[1], next[2])
		catch = true
	case endBranchRe.MatchString(cmd):
		// /list query <list_name>
		retMsg = service.List().QueryListReturnRes(ctx, cmd)
		catch = true
	}
	return
}

func tryListOp(ctx context.Context, cmd string) (catch bool, retMsg string) {
	args := strings.Fields(cmd)
	if len(args) != 4 {
		return
	}
	catch = true
	switch args[1] {
	case "U":
		// /list op <A> U <B> <C>
		// `A` Union `B` equals `C`
		retMsg = service.List().UnionListReturnRes(ctx, args[0], args[2], args[3])
	case "I":
		// /list op <A> I <B> <C>
		// `A` Intersect `B` equals `C`
		retMsg = service.List().IntersectListReturnRes(ctx, args[0], args[2], args[3])
	case "D":
		// /list op <A> D <B> <C>
		// `A` Difference `B` equals `C`
		retMsg = service.List().DifferenceListReturnRes(ctx, args[0], args[2], args[3])
	default:
		catch = false
		return
	}
	return
}
