package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupLog(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "set":
			// /group log set <>
			catch = tryGroupLogSet(ctx, next[2])
		case "rm":
			// /group log rm <>
			catch = tryGroupLogRemove(ctx, next[2])
		}
	}
	return
}

func tryGroupLogSet(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "leave":
			// /group log set leave <list_name>
			service.Group().SetLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	}
	return
}

func tryGroupLogRemove(ctx context.Context, cmd string) (catch bool) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "leave":
			// /group log rm leave
			service.Group().RemoveLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}
