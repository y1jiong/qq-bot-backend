package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"strings"
	"time"
)

func (s *sEvent) TryGroupKeywordReply(ctx context.Context) (catch bool) {
	// 获取基础信息
	msg := service.Bot().GetMessage(ctx)
	groupId := service.Bot().GetGroupId(ctx)
	// 匹配 @bot
	if atPrefixRe.MatchString(msg) {
		sub := atPrefixRe.FindStringSubmatch(msg)
		if sub[1] == gconv.String(service.Bot().GetSelfId(ctx)) {
			msg = strings.Replace(msg, sub[0], "", 1)
		}
	}
	// 匹配关键词
	var lists map[string]any
	if service.Group().IsBinding(ctx, groupId) {
		lists = service.Group().GetKeywordReplyLists(ctx, groupId)
	} else {
		lists = service.Namespace().GetGlobalNamespaceLists(ctx)
	}
	contains, hit, value := service.Util().IsOnKeywordLists(ctx, msg, lists)
	if !contains || value == "" {
		return
	}
	// 匹配成功，回复
	replyMsg := value
	noReplyPrefix := false
	switch {
	case webhookPrefixRe.MatchString(value):
		replyMsg, noReplyPrefix = s.keywordReplyWebhook(ctx,
			service.Bot().GetUserId(ctx), groupId, service.Bot().GetCardOrNickname(ctx),
			msg, hit, value)
	case rewritePrefixRe.MatchString(value):
		catch = s.keywordReplyRewrite(ctx, s.TryGroupKeywordReply, msg, hit, value)
		replyMsg = ""
	case commandPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyCommand(ctx, msg, hit, value)
	}
	// 内容为空，不回复
	if replyMsg == "" {
		return
	}
	// 限速
	const kind = "replyG"
	gid := gconv.String(groupId)
	if limited, _ := service.Util().AutoLimit(ctx, kind, gid, 7, time.Minute); limited {
		g.Log().Notice(ctx, kind, gid, "is limited")
		return
	}
	if !noReplyPrefix {
		replyMsg = "[CQ:reply,id=" + gconv.String(service.Bot().GetMsgId(ctx)) + "]" + replyMsg
	}
	service.Bot().SendMsg(ctx, replyMsg)
	catch = true
	return
}
