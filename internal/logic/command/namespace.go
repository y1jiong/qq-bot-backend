package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func tryNamespace(ctx context.Context, args []string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryNamespace")
	defer span.End()

	switch {
	case len(args) > 1:
		switch args[0] {
		case "add":
			// /namespace add <namespace>
			retMsg = service.Namespace().AddNewNamespaceReturnRes(ctx, args[1])
			caught = true
		case "rm":
			// /namespace rm <namespace>
			retMsg = service.Namespace().RemoveNamespaceReturnRes(ctx, args[1])
			caught = true
		case "chown":
			// /namespace chown <owner_id> <namespace>
			if len(args) < 3 {
				break
			}
			retMsg = service.Namespace().ChangeNamespaceOwnerReturnRes(ctx, args[1], args[2])
			caught = true
		default:
			// /namespace <namespace> <>
			caught, retMsg = tryNamespaceNext(ctx, args[0], args[1:])
		}
	case len(args) == 1:
		switch args[0] {
		case "query":
			// /namespace query
			retMsg = service.Namespace().QueryOwnNamespaceReturnRes(ctx)
			caught = true
		default:
			// /namespace <namespace>
			retMsg = service.Namespace().QueryNamespaceReturnRes(ctx, args[0])
			caught = true
		}
	}
	return
}

func tryNamespaceNext(ctx context.Context, namespace string, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "load":
			// /namespace <namespace> load <>
			caught, retMsg = tryNamespaceLoad(ctx, namespace, args[1:])
		case "unload":
			// /namespace <namespace> unload <>
			caught, retMsg = tryNamespaceUnload(ctx, namespace, args[1:])
		case "set":
			// /namespace <namespace> set <>
			caught, retMsg = tryNamespaceSet(ctx, namespace, args[1:])
		case "reset":
			// /namespace <namespace> reset <>
			caught, retMsg = tryNamespaceReset(ctx, namespace, args[1:])
		}
	}
	return
}

func tryNamespaceLoad(ctx context.Context, namespace string, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "list":
			// /namespace <namespace> load list <list_name>
			retMsg = service.Namespace().LoadNamespaceListReturnRes(ctx, namespace, args[1])
			caught = true
		}
	}
	return
}

func tryNamespaceUnload(ctx context.Context, namespace string, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case "list":
			// /namespace <namespace> unload list <list_name>
			retMsg = service.Namespace().UnloadNamespaceListReturnRes(ctx, namespace, args[1])
			caught = true
		}
	}
	return
}

func tryNamespaceSet(ctx context.Context, namespace string, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch args[0] {
		case consts.PublicCmd:
			// /namespace <namespace> set public <true|false>
			retMsg = service.Namespace().SetNamespacePropertyPublicReturnRes(ctx, namespace, args[1] == "true")
			caught = true
		}
	}
	return
}

func tryNamespaceReset(ctx context.Context, namespace string, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) == 1:
		switch args[0] {
		case "admin":
			// /namespace <namespace> reset admin
			retMsg = service.Namespace().ResetNamespaceAdminReturnRes(ctx, namespace)
			caught = true
		}
	}
	return
}
