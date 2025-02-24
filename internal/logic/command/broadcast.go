package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/y1jiong/go-shellquote"
	"qq-bot-backend/internal/service"
	"strings"
)

func tryBroadcast(ctx context.Context, args []string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryBroadcast")
	defer span.End()

	switch {
	case len(args) > 1:
		switch args[0] {
		case "group":
			// /broadcast group <>
			if len(args) < 3 {
				break
			}

			// /broadcast group <group_id> <...content>
			dstGroupId := gconv.Int64(args[1])
			userId := service.Bot().GetUserId(ctx)
			if dstNamespace := service.Group().GetNamespace(ctx, dstGroupId); dstNamespace == "" ||
				!service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, dstNamespace, userId) {
				break
			}

			suffix := "\n\nbroadcast from " + service.Bot().GetCardOrNickname(ctx) + "(" + gconv.String(userId) + ")"
			_, _ = service.Bot().SendMessage(ctx,
				service.Bot().GetMsgType(ctx),
				0,
				dstGroupId,
				strings.TrimSpace(shellquote.Join(args[2:]...))+suffix,
				false,
			)
			caught = true
		}
	}
	return
}
