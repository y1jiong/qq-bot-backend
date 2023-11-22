package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
)

type sCommand struct{}

func New() *sCommand {
	return &sCommand{}
}

func init() {
	service.RegisterCommand(New())
}

var (
	commandPrefixRe     = regexp.MustCompile(`^/(.+)$`)
	nextBranchRe        = regexp.MustCompile(`^(\S+)\s+(.+)$`)
	endBranchRe         = regexp.MustCompile(`^\S+$`)
	doubleValueCmdEndRe = regexp.MustCompile(`^(\S+)\s+(\S+)$`)
)

func (s *sCommand) TryCommand(ctx context.Context) (catch bool, retMsg string) {
	msg := service.Bot().GetMessage(ctx)
	if !commandPrefixRe.MatchString(msg) {
		return
	}
	// 暂停状态时的权限校验
	userId := service.Bot().GetUserId(ctx)
	if !service.Process().IsEnabledBotProcess() &&
		!service.User().IsSystemTrustUser(ctx, userId) {
		return
	}
	// 命令 log
	g.Log().Info(ctx,
		"user("+gconv.String(userId)+") in group("+gconv.String(service.Bot().GetGroupId(ctx))+") send cmd "+msg)
	// 继续处理
	cmd := commandPrefixRe.FindStringSubmatch(msg)[1]
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "list":
			// /list <>
			catch, retMsg = tryList(ctx, next[2])
		case "group":
			// /group <>
			catch, retMsg = tryGroup(ctx, next[2])
		case "namespace":
			// /namespace <>
			catch, retMsg = tryNamespace(ctx, next[2])
		case "user":
			// /user <>
			catch, retMsg = tryUser(ctx, next[2])
		case "raw":
			// 权限校验
			if !service.User().CouldGetRawMsg(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// /raw <>
			catch, retMsg = true, next[2]
		case "model":
			// /model <>
			catch, retMsg = tryModelSet(ctx, next[2])
		case "token":
			// /token <>
			catch, retMsg = tryToken(ctx, next[2])
		case "sys":
			// /sys <>
			catch, retMsg = trySys(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		// 权限校验
		if !service.User().IsSystemTrustUser(ctx, service.Bot().GetUserId(ctx)) {
			return
		}
		switch cmd {
		case "status":
			// /status
			catch, retMsg = queryProcessStatus(ctx)
		case "version":
			// /version
			retMsg = consts.Description
			catch = true
		case "continue":
			// /continue
			catch, retMsg = continueProcess(ctx)
		case "pause":
			// /pause
			catch, retMsg = pauseProcess(ctx)
		}
	}
	return
}
