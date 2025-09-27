package event

import (
	"context"
	"fmt"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
)

func (s *sEvent) TryKeywordRecall(ctx context.Context) (caught bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryKeywordRecall")
	defer span.End()

	if service.Bot().IsGroupOwnerOrAdmin(ctx) {
		// owner or admin 不需要撤回
		return
	}
	// 获取当前 group keyword 策略
	groupId := service.Bot().GetGroupId(ctx)
	policy := service.Group().GetKeywordPolicy(ctx, groupId)
	// 预处理
	if len(policy) == 0 {
		// 没有关键词检查策略，跳过撤回功能
		return
	}
	// 获取聊天信息
	msg := service.Util().ToPlainText(ctx, service.Bot().GetMessage(ctx))
	shouldRecall := false
	// 命中规则
	hit := ""
	// 处理
	if _, ok := policy[consts.BlacklistCmd]; ok {
		shouldRecall, hit, _ = service.Util().FindBestKeywordMatch(ctx, msg, service.Group().GetKeywordBlacklists(ctx, groupId))
	}
	if _, ok := policy[consts.WhitelistCmd]; ok && shouldRecall {
		found, _, _ := service.Util().FindBestKeywordMatch(ctx, msg, service.Group().GetKeywordWhitelists(ctx, groupId))
		shouldRecall = !found
	}
	// 结果处理
	if !shouldRecall {
		// 不需要撤回
		return
	}
	userId := service.Bot().GetUserId(ctx)
	// 撤回日志
	logMsg := fmt.Sprintf("recall group(%v) %v(%v) hit(%v) detail %v",
		groupId,
		service.Bot().GetCardOrNickname(ctx),
		userId,
		hit,
		msg,
	)
	// 撤回
	if err := service.Bot().RecallMessage(ctx, service.Bot().GetMsgId(ctx)); err != nil {
		g.Log().Warning(ctx, logMsg, err)
	} else {
		g.Log().Info(ctx, logMsg)
	}
	// 通知
	if notificationGroupId :=
		service.Group().GetMessageNotificationGroupId(ctx, groupId); notificationGroupId != 0 {
		if _, err := service.Bot().SendMessage(ctx, 0, notificationGroupId, logMsg, true); err != nil {
			g.Log().Warning(ctx, err)
		}
	}
	// 禁言
	service.Util().AutoMute(ctx, "keyword", groupId, userId,
		1, 5, 0, 16*time.Hour)

	caught = true
	return
}
