package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryToken(ctx context.Context, cmd string) (catch bool) {
	// 权限校验
	if !service.User().IsSystemTrustUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 继续处理
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /token add <>
			catch = tryTokenAdd(ctx, next[2])
		case "rm":
			// /token rm <>
			service.Token().RemoveToken(ctx, next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		if endBranchRe.FindString(cmd) == "query" {
			// /token query
			service.Token().QueryToken(ctx)
			catch = true
		}
	}
	return
}

func tryTokenAdd(ctx context.Context, cmd string) (catch bool) {
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token add <name> <token>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	service.Token().AddNewToken(ctx, dv[1], dv[2], service.Bot().GetUserId(ctx))
	catch = true
	return
}
