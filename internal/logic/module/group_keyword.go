package module

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func (s *sModule) TryKeywordRecall(ctx context.Context) (catch bool) {
	if service.Bot().IsGroupOwnerOrAdmin(ctx) {
		// owner or admin 不需要撤回
		return
	}
	// 获取当前 group keyword 策略
	groupId := service.Bot().GetGroupId(ctx)
	process := service.Group().GetKeywordProcess(ctx, groupId)
	// 预处理
	if len(process) < 1 {
		// 没有关键词检查策略，跳过撤回功能
		return
	}
	// 获取聊天信息
	msg := service.Bot().GetMessage(ctx)
	shouldRecall := false
	// 命中规则
	hit := ""
	// 处理
	if _, ok := process[consts.BlacklistCmd]; ok {
		shouldRecall, hit, _ = s.isOnKeywordLists(ctx, msg, service.Group().GetKeywordBlacklists(ctx, groupId))
	}
	if _, ok := process[consts.WhitelistCmd]; ok && shouldRecall {
		in, _, _ := s.isOnKeywordLists(ctx, msg, service.Group().GetKeywordWhitelists(ctx, groupId))
		shouldRecall = !in
	}
	// 结果处理
	if !shouldRecall {
		// 不需要撤回
		return
	}
	// 撤回
	service.Bot().RecallMessage(ctx, service.Bot().GetMsgId(ctx))
	userId := service.Bot().GetUserId(ctx)
	// 打印撤回日志
	logMsg := fmt.Sprintf("recall group(%v) user(%v) hit(%v) detail %v",
		groupId,
		userId,
		hit,
		msg)
	g.Log().Info(ctx, logMsg)
	// 通知
	notificationGroupId := service.Group().GetMessageNotificationGroupId(ctx, groupId)
	if notificationGroupId > 0 {
		service.Bot().SendMessage(ctx,
			"group", 0, notificationGroupId, logMsg, true)
	}
	// 禁言
	s.AutoMute(ctx, "keyword", groupId, userId,
		1, 5, 0, gconv.Duration("16h"))
	catch = true
	return
}
