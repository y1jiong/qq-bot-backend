package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sModule) TryGroupKeywordReply(ctx context.Context) (catch bool) {
	// 获取基础信息
	msg := service.Bot().GetMessage(ctx)
	groupId := service.Bot().GetGroupId(ctx)
	// 匹配关键词
	var lists map[string]any
	if service.Group().IsBinding(ctx, groupId) {
		lists = service.Group().GetKeywordReplyLists(ctx, groupId)
	} else {
		lists = service.Namespace().GetPublicNamespaceLists(ctx)
	}
	contains, hit, value := s.isOnKeywordLists(ctx, msg, lists)
	if !contains || value == "" {
		return
	}
	// 限速
	kind := "replyG"
	gid := gconv.String(groupId)
	if limited, _ := s.AutoLimit(ctx, kind, gid, 7, time.Minute); limited {
		g.Log().Info(ctx, kind, gid, "is limited")
		return
	}
	// 匹配成功，回复
	replyMsg := value
	noReplyPrefix := false
	switch {
	case webhookPrefixRe.MatchString(value):
		replyMsg, noReplyPrefix = s.keywordReplyWebhook(ctx, service.Bot().GetUserId(ctx), groupId, msg, hit, value)
	case commandPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyCommand(ctx, msg, hit, value)
	}
	// 内容为空，不回复
	if replyMsg == "" {
		return
	}
	if !noReplyPrefix {
		replyMsg = "[CQ:reply,id=" + gconv.String(service.Bot().GetMsgId(ctx)) + "]" + replyMsg
	}
	service.Bot().SendMsg(ctx, replyMsg)
	catch = true
	return
}
