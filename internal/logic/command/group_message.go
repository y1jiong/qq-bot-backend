package command

import (
	"context"
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
			service.Group().SetAntiRecallWithRes(ctx, service.Bot().GetGroupId(ctx), true)
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
			service.Group().SetAntiRecallWithRes(ctx, service.Bot().GetGroupId(ctx), false)
			catch = true
		}
	}
	return
}
