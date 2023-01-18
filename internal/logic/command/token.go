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
			catch = tryTokenRemove(ctx, next[2])
		}
	case singleValueCmdEndRe.MatchString(cmd):
		if singleValueCmdEndRe.FindString(cmd) == "query" {
			// /token query
			catch = tryTokenQuery(ctx)
		}
	}
	return
}

func tryTokenAdd(ctx context.Context, cmd string) (catch bool) {
	catch = true
	if !doubleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token add <name> <token>
	dv := doubleValueCmdEndRe.FindStringSubmatch(cmd)
	// 执行
	service.Token().AddNewToken(ctx, dv[1], dv[2], service.Bot().GetUserId(ctx))
	return
}

func tryTokenRemove(ctx context.Context, cmd string) (catch bool) {
	catch = true
	if !singleValueCmdEndRe.MatchString(cmd) {
		return
	}
	// /token rm <name>
	name := singleValueCmdEndRe.FindString(cmd)
	// 执行
	service.Token().RemoveToken(ctx, name)
	return
}

func tryTokenQuery(ctx context.Context) (catch bool) {
	catch = true
	// /token query
	// 执行
	service.Token().QueryToken(ctx)
	return
}
