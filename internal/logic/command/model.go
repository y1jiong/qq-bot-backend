package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func tryModelSet(ctx context.Context, cmd string) (caught bool, retMsg string) {
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

	caught = true

	if err := service.Bot().SetModel(ctx, next[2]); err != nil {
		retMsg = err.Error()
		return
	}

	retMsg = "已更改机型为 '" + next[2] + "'"
	return
}
