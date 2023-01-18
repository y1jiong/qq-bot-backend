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
			// /namespace add <>
			// 权限校验
			if !service.User().CouldOperateNamespace(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// 继续处理
			service.Namespace().AddNewNamespace(ctx, next[2])
			catch = true
		case "rm":
			// /namespace rm <>
			// 权限校验
			if !service.User().CouldOperateNamespace(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// 继续处理
			service.Namespace().RemoveNamespace(ctx, next[2])
			catch = true
		default:
			// /namespace <namespace> <>
			catch = tryNamespaceSet(ctx, next[1], next[2])
		}
	case singleValueCmdEndRe.MatchString(cmd):
		v := singleValueCmdEndRe.FindString(cmd)
		switch v {
		case "query":
			// /namespace query
			service.Namespace().QueryOwnNamespace(ctx, service.Bot().GetUserId(ctx))
			catch = true
		default:
			// /namespace <namespace>
			service.Namespace().QueryNamespace(ctx, v)
			catch = true
		}
	}
	return
}

func tryNamespaceSet(ctx context.Context, namespace, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "reset":
			// /namespace <namespace> reset <>
			catch = tryNamespaceReset(ctx, namespace, next[2])
		}
	case singleValueCmdEndRe.MatchString(cmd):
		end := singleValueCmdEndRe.FindString(cmd)
		switch end {
		case "reset":
			// /namespace <namespace> reset
			service.Namespace().ResetNamespace(ctx, namespace, "all")
			catch = true
		}
	}
	return
}

func tryNamespaceReset(ctx context.Context, namespace, cmd string) (catch bool) {
	switch cmd {
	case "admin":
		// /namespace <namespace> reset admin
		service.Namespace().ResetNamespace(ctx, namespace, cmd)
		catch = true
	case "whitelist":
		// /namespace <namespace> reset whitelist
		service.Namespace().ResetNamespace(ctx, namespace, cmd)
		catch = true
	}
	return
}
