package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel/attribute"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
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

var (
	nextBranchRe      = regexp.MustCompile(`^(\S+)\s+([\s\S]+)$`)
	endBranchRe       = regexp.MustCompile(`^\S+$`)
	dualValueCmdEndRe = regexp.MustCompile(`^(\S+)\s+(\S+)$`)
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
	}()
	cmd := strings.Replace(message, cmdPrefix, "", 1)
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)

		switch next[1] {
		case "list":
			// /list <>
			caught, retMsg = tryList(ctx, next[2])
		case "group":
			// /group <>
			caught, retMsg = tryGroup(ctx, next[2])
		case "namespace":
			// /namespace <>
			caught, retMsg = tryNamespace(ctx, next[2])
		case "user":
			// /user <>
			caught, retMsg = tryUser(ctx, next[2])
		case "raw":
			// /raw <>
			caught, retMsg = tryRaw(ctx, next[2])
		case "broadcast":
			// /broadcast <>
			caught, retMsg = tryBroadcast(ctx, next[2])
		case "token":
			// /token <>
			caught, retMsg = tryToken(ctx, next[2])
		case "sys":
			// /sys <>
			caught, retMsg = trySys(ctx, next[2])
		case "model":
			// /model <>
			caught, retMsg = tryModelSet(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		// 权限校验
		if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
			return
		}

		switch cmd {
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

func tryRaw(ctx context.Context, cmd string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().CanGetRawMsg(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.raw")
	defer span.End()

	caught, retMsg = true, cmd
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
