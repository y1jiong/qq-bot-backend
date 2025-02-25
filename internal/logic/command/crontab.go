package command

import (
	"context"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
	"regexp"
)

var (
	crontabRe = regexp.MustCompile(`^(\S+ \S+ \S+ \S+ \S+)\s+([\s\S]+)`)
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
		case "query":
			var creatorId int64
			if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
				creatorId = userId
			}
			// /crontab query <name>
			retMsg = service.Crontab().QueryReturnRes(ctx, next[2], creatorId)
			caught = true
		case "add":
			// /crontab add <>
			caught, retMsg = tryCrontabAdd(ctx, next[2])
		case "rm":
			// /crontab rm <name>
			retMsg = service.Crontab().RemoveReturnRes(ctx, next[2], service.Bot().GetUserId(ctx))
			caught = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "glance":
			var creatorId int64
			if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
				creatorId = userId
			}
			// /crontab glance
			retMsg = service.Crontab().GlanceReturnRes(ctx, creatorId)
			caught = true
		case "reload":
			if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
				break
			}
			// /crontab reload
			service.Crontab().Run(ctx)
			retMsg = "crontab reloaded"
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
		node := service.Bot().CloneReqNode(ctx)
		if node == nil {
			break
		}
		_, _ = node.Set("raw_message", ast.NewString(next[2]))
		_, _ = node.Set("message", ast.NewString(next[2]))
		reqJSON, _ := node.MarshalJSON()
		// /crontab add <name> <expr> <message>
		retMsg = service.Crontab().AddReturnRes(ctx,
			name,
			next[1],
			service.Bot().GetUserId(ctx),
			service.Bot().GetSelfId(ctx),
			reqJSON,
		)
		caught = true
	}
	return
}
