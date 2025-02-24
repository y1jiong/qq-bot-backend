package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupExport(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "member":
			// /group export member <list_name>
			retMsg = service.Group().ExportGroupMemberListReturnRes(ctx, service.Bot().GetGroupId(ctx), args[1])
			caught = true
		}
	}
	return
}
