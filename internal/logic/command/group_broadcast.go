package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupBroadcast(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "accept":
			// /group broadcast accept
			retMsg = service.Group().AcceptBroadcastReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "reject":
			// /group broadcast reject
			retMsg = service.Group().RejectBroadcastReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		}
	}
	return
}
