package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryModelSet(ctx context.Context, cmd string) (catch bool, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	// /model set <model>
	if !nextBranchRe.MatchString(cmd) {
		return
	}
	next := nextBranchRe.FindStringSubmatch(cmd)
	if next[1] != "set" {
		return
	}
	service.Bot().SetModel(ctx, next[2])
	catch = true
	return
}
