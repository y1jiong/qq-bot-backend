package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupLog(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "set":
			// /group log set <>
			catch, retMsg = tryGroupLogSet(ctx, next[2])
		case "rm":
			// /group log rm <>
			catch, retMsg = tryGroupLogRemove(ctx, next[2])
		}
	}
	return
}

func tryGroupLogSet(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "leave":
			// /group log set leave <list_name>
			retMsg = service.Group().SetLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	}
	return
}

func tryGroupLogRemove(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "leave":
			// /group log rm leave
			retMsg = service.Group().RemoveLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}
