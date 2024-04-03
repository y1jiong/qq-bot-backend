package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sModule) TryKeywordReply(ctx context.Context) (catch bool) {
	// 获取基础信息
	msg := service.Bot().GetMessage(ctx)
	userId := service.Bot().GetUserId(ctx)
	// 匹配关键词
	contains, hit, value := s.isOnKeywordLists(ctx, msg, service.Namespace().GetPublicNamespaceLists(ctx))
	if !contains || value == "" {
		return
	}
	// 限速
	kind := "replyU"
	uid := gconv.String(userId)
	if limited, _ := s.AutoLimit(ctx, kind, uid, 5, time.Minute); limited {
		g.Log().Info(ctx, kind, uid, "is limited")
		return
	}
	// 匹配成功，回复
	replyMsg := value
	switch {
	case webhookPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyWebhook(ctx, userId, 0, msg, hit, value)
	case commandPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyCommand(ctx, msg, hit, value)
	}
	// 内容为空，不回复
	if replyMsg == "" {
		return
	}
	pre := "[CQ:reply,id=" + gconv.String(service.Bot().GetMsgId(ctx)) + "]" + replyMsg
	service.Bot().SendMsg(ctx, pre)
	catch = true
	return
}
