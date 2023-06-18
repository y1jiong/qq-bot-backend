package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryNamespace(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// 权限校验
			if !service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// /namespace add <namespace>
			// 继续处理
			service.Namespace().AddNewNamespaceReturnRes(ctx, next[2])
			catch = true
		case "rm":
			// 权限校验
			if !service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// /namespace rm <namespace>
			// 继续处理
			service.Namespace().RemoveNamespaceReturnRes(ctx, next[2])
			catch = true
		default:
			// /namespace <namespace> <>
			catch = tryNamespaceReset(ctx, next[1], next[2])
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// 权限校验
			if !service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// /namespace query
			service.Namespace().QueryOwnNamespaceReturnRes(ctx)
			catch = true
		default:
			// /namespace <namespace>
			service.Namespace().QueryNamespaceReturnRes(ctx, cmd)
			catch = true
		}
	}
	return
}

func tryNamespaceReset(ctx context.Context, namespace, cmd string) (catch bool) {
	if endBranchRe.MatchString(cmd) {
		switch cmd {
		case "reset":
			// /namespace <namespace> reset
			service.Namespace().ResetNamespaceAdminReturnRes(ctx, namespace)
			catch = true
		}
	}
	return
}
