package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroup(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "bind":
			// /group bind <>
			if singleValueCmdEndRe.MatchString(next[2]) {
				// /group bind <namespace>
				service.Group().BindNamespace(ctx, service.Bot().GetGroupId(ctx), next[2])
				catch = true
			}
		}
	case singleValueCmdEndRe.MatchString(cmd):
		v := singleValueCmdEndRe.FindString(cmd)
		switch v {
		case "bind":
			// /group bind
			service.Group().BindQuery(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		case "unbind":
			// /group unbind
			service.Group().Unbind(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}
