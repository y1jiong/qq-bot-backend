package command

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel/attribute"
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

type catch uint8

const (
	notCaught catch = iota
	caughtNeedOkay
	caughtNoOkay
)

var (
	nextBranchRe      = regexp.MustCompile(`^(\S+)\s+([\s\S]+)$`)
	endBranchRe       = regexp.MustCompile(`^\S+$`)
	dualValueCmdEndRe = regexp.MustCompile(`^(\S+)\s+(\S+)$`)
)

func (s *sCommand) TryCommand(ctx context.Context, message string) (caught bool, retMsg string) {
	c, retMsg := s.tryCommand(ctx, message)
	return c != notCaught, retMsg
}

func (s *sCommand) tryCommand(ctx context.Context, message string) (caught catch, retMsg string) {
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
		if caught == notCaught {
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
		if caught == caughtNeedOkay && retMsg == "" {
			if err := service.Bot().Okay(ctx); err != nil {
				g.Log().Warning(ctx, err)
			}
		}
	}()

	cmd := strings.Replace(message, cmdPrefix, "", 1)
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)

		switch next[1] {
		case "say":
			// /say <>
			caught, retMsg = trySay(ctx, next[2])
		case "raw":
			// /raw <>
			caught, retMsg = tryRaw(ctx, next[2])
		case "like":
			// /like <>
			caught, retMsg = tryLike(ctx, next[2])
		case "list":
			// /list <>
			caught, retMsg = tryList(ctx, next[2])
		case "group":
			// /group <>
			caught, retMsg = tryGroup(ctx, next[2])
		case "crontab":
			// /crontab <>
			caught, retMsg = tryCrontab(ctx, next[2])
		case "namespace":
			// /namespace <>
			caught, retMsg = tryNamespace(ctx, next[2])
		case "broadcast":
			// /broadcast <>
			caught, retMsg = tryBroadcast(ctx, next[2])
		case "user":
			// /user <>
			caught, retMsg = tryUser(ctx, next[2])
		case "sys":
			// /sys <>
			caught, retMsg = trySys(ctx, next[2])
		case "token":
			// /token <>
			caught, retMsg = tryToken(ctx, next[2])
		case "model":
			// /model <>
			caught, retMsg = tryModelSet(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		// 权限校验
		if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
			break
		}

		switch cmd {
		case "readall":
			// /readall
			caught, retMsg = tryReadAll(ctx)
		case "version":
			// /version
			caught, retMsg = tryVersion(ctx)
		case "status":
			// /status
			caught, retMsg = queryProcessStatus(ctx)
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

func tryReadAll(ctx context.Context) (caught catch, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.readall")
	defer span.End()

	caught = caughtNeedOkay

	if err := service.Bot().MarkAllAsRead(ctx); err != nil {
		retMsg = err.Error()
	}
	return
}

func tryVersion(ctx context.Context) (caught catch, retMsg string) {
	ctx, span := gtrace.NewSpan(ctx, "command.version")
	defer span.End()

	caught, retMsg = caughtNeedOkay, consts.Description

	appName, appVersion, protocolVersion, err := service.Bot().GetVersionInfo(ctx)
	if err != nil {
		return
	}
	// appName/appVersion (protocolVersion)
	retMsg += "\n" + appName + "/" + appVersion + " (" + protocolVersion + ")"
	return
}

func trySay(ctx context.Context, cmd string) (caught catch, retMsg string) {
	// 权限校验
	if !service.User().CanGetRawMessage(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.say")
	defer span.End()

	service.Bot().SendMsgCacheContext(ctx, cmd)

	caught = caughtNoOkay
	return
}

func tryRaw(ctx context.Context, cmd string) (caught catch, retMsg string) {
	// 权限校验
	if !service.User().CanGetRawMessage(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.raw")
	defer span.End()

	caught, retMsg = caughtNeedOkay, cmd
	return
}

func tryLike(ctx context.Context, cmd string) (caught catch, retMsg string) {
	// 权限校验
	if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.like")
	defer span.End()

	if !dualValueCmdEndRe.MatchString(cmd) {
		return
	}
	dv := dualValueCmdEndRe.FindStringSubmatch(cmd)

	caught = caughtNeedOkay

	// /like <user_id> <times>
	if err := service.Bot().ProfileLike(ctx, gconv.Int64(dv[1]), gconv.Int(dv[2])); err != nil {
		retMsg = err.Error()
		return
	}
	return
}
