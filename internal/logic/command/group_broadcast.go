package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupBroadcast(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "accept":
			// /group broadcast accept
			retMsg = service.Group().AcceptBroadcastReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = caughtOkay
		case "reject":
			// /group broadcast reject
			retMsg = service.Group().RejectBroadcastReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = caughtOkay
		}
	}
	return
}
