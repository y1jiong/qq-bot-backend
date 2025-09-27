package command

import (
	"context"
	"qq-bot-backend/internal/service"

	"github.com/gogf/gf/v2/net/gtrace"
)

func tryModelSet(ctx context.Context, cmd string) (caught catch, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.modelSet")
	defer span.End()

	// /model set <model>
	if !nextBranchRe.MatchString(cmd) {
		return
	}
	next := nextBranchRe.FindStringSubmatch(cmd)
	if next[1] != "set" {
		return
	}

	caught = caughtOkay

	if err := service.Bot().SetModel(ctx, next[2]); err != nil {
		retMsg = err.Error()
		return
	}

	retMsg = "已更改机型为 '" + next[2] + "'"
	return
}
