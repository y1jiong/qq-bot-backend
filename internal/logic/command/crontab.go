package command

import (
	"context"
	"qq-bot-backend/internal/service"
	"regexp"

	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	crontabRe = regexp.MustCompile(`^(\S+ \S+ \S+ \S+ \S+)\s+([\s\S]+)`)
)

func tryCrontab(ctx context.Context, cmd string) (caught catch, retMsg string) {
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
		case "oneshot":
			var creatorId int64
			if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
				creatorId = userId
			}
			// /crontab oneshot <name>
			retMsg = service.Crontab().OneshotReturnRes(ctx, next[2], creatorId)
			caught = caughtNeedOkay
		case "query":
			var creatorId int64
			if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
				creatorId = userId
			}
			// /crontab query <name>
			retMsg = service.Crontab().QueryReturnRes(ctx, next[2], creatorId)
			caught = caughtNeedOkay
		case "add":
			// /crontab add <>
			caught, retMsg = tryCrontabAdd(ctx, next[2])
		case "rm":
			var creatorId int64
			if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
				creatorId = userId
			}
			// /crontab rm <name>
			retMsg = service.Crontab().RemoveReturnRes(ctx, next[2], creatorId)
			caught = caughtNeedOkay
		case "ch-expr":
			// /crontab ch-expr <>
			caught, retMsg = tryCrontabChExpr(ctx, next[2])
		case "ch-bind":
			// /crontab ch-bind <>
			caught, retMsg = tryCrontabChBind(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "show":
			var creatorId int64
			if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
				creatorId = userId
			}
			// /crontab show
			retMsg = service.Crontab().ShowReturnRes(ctx, creatorId)
			caught = caughtNeedOkay
		case "reload":
			if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
				break
			}
			// /crontab reload
			service.Crontab().Run(ctx)
			retMsg = "crontab reloaded"
			caught = caughtNeedOkay
		}
	}
	return
}

func tryCrontabAdd(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		// /crontab add <name> <>
		next := nextBranchRe.FindStringSubmatch(cmd)
		name := next[1]
		if !crontabRe.MatchString(next[2]) {
			break
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
		caught = caughtNeedOkay
	}
	return
}

func tryCrontabChExpr(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case crontabRe.MatchString(cmd):
		// /crontab ch-expr <expr> <>
		next := crontabRe.FindStringSubmatch(cmd)
		if !endBranchRe.MatchString(next[2]) {
			break
		}
		expr := next[1]
		name := next[2]
		var creatorId int64
		if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
			creatorId = userId
		}
		// /crontab ch-expr <expr> <name>
		retMsg = service.Crontab().ChangeExpressionReturnRes(ctx, expr, name, creatorId)
		caught = caughtNeedOkay
	}
	return
}

func tryCrontabChBind(ctx context.Context, cmd string) (caught catch, retMsg string) {
	switch {
	case dualValueCmdEndRe.MatchString(cmd):
		dv := dualValueCmdEndRe.FindStringSubmatch(cmd)
		var creatorId int64
		if userId := service.Bot().GetUserId(ctx); !service.User().IsSystemTrustedUser(ctx, userId) {
			creatorId = userId
		}
		// /crontab ch-bind <bot_id> <name>
		retMsg = service.Crontab().ChangeBotIdReturnRes(ctx, gconv.Int64(dv[1]), dv[2], creatorId)
		caught = caughtNeedOkay
	}
	return
}
