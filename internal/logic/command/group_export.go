package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroupExport(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "member":
			// /group export member <list_name>
			service.Group().ExportGroupMemberListReturnRes(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	}
	return
}
