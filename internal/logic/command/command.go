package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/y1jiong/go-shellquote"
	"go.opentelemetry.io/otel/attribute"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"strings"
)

type sCommand struct{}

func New() *sCommand {
	return &sCommand{}
}

func init() {
	service.RegisterCommand(New())
}

const (
	cmdPrefix = "/"
)

func (s *sCommand) TryCommand(ctx context.Context, message string) (caught bool, retMsg string) {
	if !strings.HasPrefix(message, cmdPrefix) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.TryCommand")
	defer span.End()

	// 暂停状态时的权限校验
	userId := service.Bot().GetUserId(ctx)
	if !service.Process().IsBotProcessEnabled() &&
		!service.User().IsSystemTrustedUser(ctx, userId) {
		return
	}
	// 命令 log
	defer func() {
		if !caught {
			return
		}
		groupId := service.Bot().GetGroupId(ctx)
		span.SetAttributes(
			attribute.Int64("try_command.user_id", userId),
			attribute.Int64("try_command.group_id", groupId),
			attribute.String("try_command.command", message),
		)
		g.Log().Info(ctx,
			service.Bot().GetCardOrNickname(ctx)+"("+gconv.String(userId)+
				") in group("+gconv.String(groupId)+") send cmd "+message)
		if retMsg == "" {
			_ = service.Bot().Okay(ctx)
		}
	}()

	args, err := shellquote.Split(strings.Replace(message, cmdPrefix, "", 1))
	if err != nil {
		return
	}
	switch {
	case len(args) > 1:
		switch args[0] {
		case "list":
			// /list <>
			caught, retMsg = tryList(ctx, args[1:])
		case "group":
			// /group <>
			caught, retMsg = tryGroup(ctx, args[1:])
		case "namespace":
			// /namespace <>
			caught, retMsg = tryNamespace(ctx, args[1:])
		case "user":
			// /user <>
			caught, retMsg = tryUser(ctx, args[1:])
		case "raw":
			// /raw <>
			caught, retMsg = tryRaw(ctx, args[1:])
		case "like":
			// /like <>
			caught, retMsg = tryLike(ctx, args[1:])
		case "broadcast":
			// /broadcast <>
			caught, retMsg = tryBroadcast(ctx, args[1:])
		case "token":
			// /token <>
			caught, retMsg = tryToken(ctx, args[1:])
		case "sys":
			// /sys <>
			caught, retMsg = trySys(ctx, args[1:])
		case "model":
			// /model <>
			caught, retMsg = tryModelSet(ctx, args[1:])
		}
	case len(args) == 1:
		// 权限校验
		if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
			return
		}

		switch args[0] {
		case "status":
			// /status
			caught, retMsg = queryProcessStatus(ctx)
		case "version":
			// /version
			caught, retMsg = tryVersion(ctx)
		case "continue":
			// /continue
			caught, retMsg = continueProcess(ctx)
		case "pause":
			// /pause
			caught, retMsg = pauseProcess(ctx)
		}
	}
	return
}

func tryRaw(ctx context.Context, args []string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().CanGetRawMessage(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.raw")
	defer span.End()

	caught, retMsg = true, shellquote.Join(args...)
	return
}

func tryVersion(ctx context.Context) (caught bool, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.version")
	defer span.End()

	caught, retMsg = true, consts.Description

	appName, appVersion, protocolVersion, err := service.Bot().GetVersionInfo(ctx)
	if err != nil {
		return
	}
	// appName/appVersion (protocolVersion)
	retMsg += "\n" + appName + "/" + appVersion + " (" + protocolVersion + ")"
	return
}

func tryLike(ctx context.Context, args []string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.like")
	defer span.End()

	if len(args) < 2 {
		return
	}

	caught = true

	// /like <user_id> <times>
	if err := service.Bot().ProfileLike(ctx, gconv.Int64(args[0]), gconv.Int(args[1])); err != nil {
		retMsg = err.Error()
		return
	}
	return
}
