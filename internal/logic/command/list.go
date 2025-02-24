package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/y1jiong/go-shellquote"
	"qq-bot-backend/internal/service"
)

func tryList(ctx context.Context, args []string) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.tryList")
	defer span.End()

	switch {
	case len(args) > 1:
		switch args[0] {
		case "join":
			// /list join <>
			caught, retMsg = tryListJoin(ctx, args[1:])
		case "leave":
			// /list leave <>
			if len(args) < 3 {
				break
			}
			// /list leave <list_name> <key>
			retMsg = service.List().RemoveListDataReturnRes(ctx, args[1], args[2])
			caught = true
		case "copy-key":
			// /list copy-key <>
			if len(args) < 4 {
				break
			}
			// /list copy-key <list_name> <src_key> <dst_key>
			retMsg = service.List().CopyListKeyReturnRes(ctx, args[1], args[2], args[3])
			caught = true
		case "glance":
			// /list glance <list_name>
			retMsg = service.List().GlanceListDataReturnRes(ctx, args[1])
			caught = true
		case "query":
			// /list query <>
			caught, retMsg = tryListQuery(ctx, args[1:])
		case "len":
			// /list len <list_name>
			retMsg = service.List().QueryListLenReturnRes(ctx, args[1])
			caught = true
		case "export":
			// /list export <list_name>
			retMsg = service.List().ExportListReturnRes(ctx, args[1])
			caught = true
		case "append":
			// /list append <>
			if len(args) < 3 {
				break
			}
			// /list append <list_name> <...json>
			retMsg = service.List().AppendListDataReturnRes(ctx, args[1], shellquote.Join(args[2:]...))
			caught = true
		case "set":
			// /list set <>
			if len(args) < 3 {
				break
			}
			// /list set <list_name> <...json>
			retMsg = service.List().SetListDataReturnRes(ctx, args[1], shellquote.Join(args[2:]...))
			caught = true
		case "reset":
			// /list reset <list_name>
			retMsg = service.List().ResetListDataReturnRes(ctx, args[1])
			caught = true
		case "add":
			if len(args) < 3 {
				break
			}
			// /list add <list_name> <namespace>
			retMsg = service.List().AddListReturnRes(ctx, args[1], args[2])
			caught = true
		case "rm":
			// /list rm <list_name>
			retMsg = service.List().RemoveListReturnRes(ctx, args[1])
			caught = true
		case "recover":
			// /list recover <list_name>
			retMsg = service.List().RecoverListReturnRes(ctx, args[1])
			caught = true
		case "op":
			// /list op <>
			caught, retMsg = tryListOp(ctx, args[1:])
		}
	}
	return
}

func tryListJoin(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		switch {
		case len(args) > 2:
			// /list join <list_name> <key> [...value]
			retMsg = service.List().AddListDataReturnRes(ctx, args[0], args[1], shellquote.Join(args[2:]...))
			caught = true
		case len(args) == 2:
			// /list join <list_name> <key>
			retMsg = service.List().AddListDataReturnRes(ctx, args[0], args[1])
			caught = true
		}
	}
	return
}

func tryListQuery(ctx context.Context, args []string) (caught bool, retMsg string) {
	switch {
	case len(args) > 1:
		// /list query <list_name> [key]
		retMsg = service.List().QueryListReturnRes(ctx, args[1], args[2])
		caught = true
	case len(args) == 1:
		// /list query <list_name>
		retMsg = service.List().QueryListReturnRes(ctx, args[1])
		caught = true
	}
	return
}

func tryListOp(ctx context.Context, args []string) (caught bool, retMsg string) {
	if len(args) < 4 {
		return
	}
	caught = true
	switch args[1] {
	case "U":
		// /list op <A> U <B> <C>
		// `A` Union `B` equals `C`
		retMsg = service.List().UnionListReturnRes(ctx, args[0], args[2], args[3])
	case "I":
		// /list op <A> I <B> <C>
		// `A` Intersect `B` equals `C`
		retMsg = service.List().IntersectListReturnRes(ctx, args[0], args[2], args[3])
	case "D":
		// /list op <A> D <B> <C>
		// `A` Difference `B` equals `C`
		retMsg = service.List().DifferenceListReturnRes(ctx, args[0], args[2], args[3])
	default:
		caught = false
		return
	}
	return
}
