package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryNamespace(ctx context.Context, cmd string) (catch bool, retMsg string) {
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
			retMsg = service.Namespace().AddNewNamespaceReturnRes(ctx, next[2])
			catch = true
		case "rm":
			// 权限校验
			if !service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// /namespace rm <namespace>
			// 继续处理
			retMsg = service.Namespace().RemoveNamespaceReturnRes(ctx, next[2])
			catch = true
		case "chown":
			// /namespace chown <owner_id> <namespace>
			// 继续处理
			if !dualValueCmdEndRe.MatchString(next[2]) {
				return
			}
			dv := dualValueCmdEndRe.FindStringSubmatch(next[2])
			retMsg = service.Namespace().ChangeNamespaceOwnerReturnRes(ctx, dv[2], dv[1])
			catch = true
		default:
			// /namespace <namespace> <>
			catch, retMsg = tryNamespaceNext(ctx, next[1], next[2])
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// /namespace query
			retMsg = service.Namespace().QueryOwnNamespaceReturnRes(ctx)
			catch = true
		default:
			// /namespace <namespace>
			retMsg = service.Namespace().QueryNamespaceReturnRes(ctx, cmd)
			catch = true
		}
	}
	return
}

func tryNamespaceNext(ctx context.Context, namespace, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "set":
			// /namespace <namespace> set <>
			catch, retMsg = tryNamespaceSet(ctx, namespace, next[2])
		case "reset":
			// /namespace <namespace> reset <>
			catch, retMsg = tryNamespaceReset(ctx, namespace, next[2])
		}
	}
	return
}

func tryNamespaceSet(ctx context.Context, namespace, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.PublicCmd:
			// /namespace <namespace> set public <true|false>
			retMsg = service.Namespace().SetNamespacePropertyPublicReturnRes(ctx, namespace, next[2] == "true")
			catch = true
		}
	}
	return
}

func tryNamespaceReset(ctx context.Context, namespace, cmd string) (catch bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "admin":
			// /namespace <namespace> reset admin
			retMsg = service.Namespace().ResetNamespaceAdminReturnRes(ctx, namespace)
			catch = true
		}
	}
	return
}
