package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupLog(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "set":
			// /group log set <>
			caught, retMsg = tryGroupLogSet(ctx, next[2])
		case "rm":
			// /group log rm <>
			caught, retMsg = tryGroupLogRemove(ctx, next[2])
		}
	}
	return
}

func tryGroupLogSet(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "approval":
			// /group log set approval <list_name>
			retMsg = service.Group().SetLogApprovalListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			caught = true
		case "leave":
			// /group log set leave <list_name>
			retMsg = service.Group().SetLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			caught = true
		}
	}
	return
}

func tryGroupLogRemove(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "approval":
			// /group log rm approval
			retMsg = service.Group().RemoveLogApprovalListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "leave":
			// /group log rm leave
			retMsg = service.Group().RemoveLogLeaveListReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		}
	}
	return
}
