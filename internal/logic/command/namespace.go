package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"

	"github.com/gogf/gf/v2/net/gtrace"
)

func tryNamespace(ctx context.Context, cmd string) (caught catch, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.namespace")
	defer span.End()

	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /namespace add <namespace>
			retMsg = service.Namespace().AddNewNamespaceReturnRes(ctx, next[2])
			caught = caughtOkay
		case "rm":
			// /namespace rm <namespace>
			retMsg = service.Namespace().RemoveNamespaceReturnRes(ctx, next[2])
			caught = caughtOkay
		case "chown":
			// /namespace chown <owner_id> <namespace>
			if !dualValueCmdEndRe.MatchString(next[2]) {
				break
			}
			dv := dualValueCmdEndRe.FindStringSubmatch(next[2])
			retMsg = service.Namespace().ChangeNamespaceOwnerReturnRes(ctx, dv[1], dv[2])
			caught = caughtOkay
		default:
			// /namespace <namespace> <>
			caught, retMsg = tryNamespaceNext(ctx, next[1], next[2])
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "query":
			// /namespace query
			retMsg = service.Namespace().QueryOwnNamespaceReturnRes(ctx)
			caught = caughtOkay
		default:
			// /namespace <namespace>
			retMsg = service.Namespace().QueryNamespaceReturnRes(ctx, cmd)
			caught = caughtOkay
		}
	}
	return
}

func tryNamespaceNext(ctx context.Context, namespace, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "load":
			// /namespace <namespace> load <>
			caught, retMsg = tryNamespaceLoad(ctx, namespace, next[2])
		case "unload":
			// /namespace <namespace> unload <>
			caught, retMsg = tryNamespaceUnload(ctx, namespace, next[2])
		case "set":
			// /namespace <namespace> set <>
			caught, retMsg = tryNamespaceSet(ctx, namespace, next[2])
		case "reset":
			// /namespace <namespace> reset <>
			caught, retMsg = tryNamespaceReset(ctx, namespace, next[2])
		}
	}
	return
}

func tryNamespaceLoad(ctx context.Context, namespace, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "list":
			if !endBranchRe.MatchString(next[2]) {
				break
			}
			// /namespace <namespace> load list <list_name>
			retMsg = service.Namespace().LoadNamespaceListReturnRes(ctx, namespace, next[2])
			caught = caughtOkay
		}
	}
	return
}

func tryNamespaceUnload(ctx context.Context, namespace, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "list":
			if !endBranchRe.MatchString(next[2]) {
				break
			}
			// /namespace <namespace> unload list <list_name>
			retMsg = service.Namespace().UnloadNamespaceListReturnRes(ctx, namespace, next[2])
			caught = caughtOkay
		}
	}
	return
}

func tryNamespaceSet(ctx context.Context, namespace, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case consts.PublicCmd:
			// /namespace <namespace> set public <true|false>
			retMsg = service.Namespace().SetNamespacePropertyPublicReturnRes(ctx, namespace, next[2] == "true")
			caught = caughtOkay
		}
	}
	return
}

func tryNamespaceReset(ctx context.Context, namespace, cmd string) (caught catch, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "admin":
			// /namespace <namespace> reset admin
			retMsg = service.Namespace().ResetNamespaceAdminReturnRes(ctx, namespace)
			caught = caughtOkay
		}
	}
	return
}
