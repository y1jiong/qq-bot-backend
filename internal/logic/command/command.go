package command

import (
	"context"
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
	nextBranchRe        = regexp.MustCompile(`^(\S+) (.+)$`)
	singleValueCmdEndRe = regexp.MustCompile(`^\S+$`)
	doubleValueCmdEndRe = regexp.MustCompile(`^(\S+) (\S+)$`)
)

func (s *sCommand) TryCommand(ctx context.Context) (catch bool) {
	msg := service.Bot().GetMessage(ctx)
	if !commandPrefixRe.MatchString(msg) {
		return
	}
	cmd := commandPrefixRe.FindStringSubmatch(msg)[1]
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "user":
			// /user <>
			catch = tryUser(ctx, next[2])
		case "group":
			// /group <>
			catch = tryGroup(ctx, next[2])
		case "namespace":
			// /namespace <>
			catch = tryNamespace(ctx, next[2])
		case "model":
			// /model <>
			catch = tryModelSet(ctx, next[2])
		case "token":
			// /token <>
			catch = tryToken(ctx, next[2])
		case "sys":
			// /sys <>
			catch = trySys(ctx, next[2])
		}
	case singleValueCmdEndRe.MatchString(cmd):
		// 权限校验
		if !service.User().IsSystemTrustUser(ctx, service.Bot().GetUserId(ctx)) {
			return
		}
		switch singleValueCmdEndRe.FindString(cmd) {
		case "continue":
			// /continue
			catch = continueProcess(ctx)
		case "pause":
			// /pause
			catch = pauseProcess(ctx)
		case "state":
			// /state
			catch = queryProcessState(ctx)
		}
	}
	return
}
