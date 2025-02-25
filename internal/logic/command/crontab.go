package command

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
	"regexp"
)

var (
	crontabRe = regexp.MustCompile(`^([0-5]?[0-9]|\*)\s([01]?[0-9]|2[0-3]|\*)\s([012]?[0-9]|3[01]|\*)\s([01]?[0-9]|1[0-2]|\*)\s([0-6]|\*)\s+([\s\S]+)`)
)

func tryCrontab(ctx context.Context, cmd string) (caught bool, retMsg string) {
	// 权限校验
	if !service.User().CanOpCrontab(ctx, service.Bot().GetUserId(ctx)) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.crontab")
	defer span.End()

	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "add":
			// /crontab add <name> <expr> <message>
			caught, retMsg = tryCrontabAdd(ctx, next[2])
		case "rm":
			// /crontab rm <name>
			retMsg = service.Crontab().Remove(ctx, next[2])
			caught = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "glance":
			// /crontab glance
			retMsg = service.Crontab().Glance(ctx)
			caught = true
		}
	}
	return
}

func tryCrontabAdd(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		// /crontab add <name> <>
		next := nextBranchRe.FindStringSubmatch(cmd)
		name := next[1]
		if !crontabRe.MatchString(next[2]) {
			return
		}
		next = crontabRe.FindStringSubmatch(next[2])
		// /crontab add <name> <expr> <message>
		retMsg = service.Crontab().Add(ctx, name, next[1], next[2])
		caught = true
	}
	return
}
