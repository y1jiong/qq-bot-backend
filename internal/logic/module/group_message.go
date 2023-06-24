package module

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sModule) TryUndoMessageRecall(ctx context.Context) (catch bool) {
	if service.Bot().IsGroupOwnerOrAdmin(ctx) {
		// owner or admin 不需要反撤回
		return
	}
	// 获取当前 group message anti-recall 策略
	groupId := service.Bot().GetGroupId(ctx)
	if !service.Group().IsEnabledAntiRecall(ctx, groupId) {
		return
	}
	if service.Bot().GetOperatorId(ctx) != service.Bot().GetUserId(ctx) {
		// 不是自己操作的撤回，不需要反撤回
		return
	}
	// 获取撤回消息的 id
	msgId := service.Bot().GetMsgId(ctx)
	// 反撤回
	messageMap, err := service.Bot().RequestMessage(ctx, msgId)
	if err != nil {
		service.Bot().SendPlainMsg(ctx, "获取历史消息失败")
		return
	}
	// 获取发送者信息
	senderMap := gconv.Map(messageMap["sender"])
	nickname := gconv.String(senderMap["nickname"])
	userId := gconv.Int64(senderMap["user_id"])
	// 获取撤回消息
	message := gconv.String(messageMap["message"])
	// 防止过度触发反撤回
	s.AutoMute(ctx, "recall", groupId, service.Bot().GetUserId(ctx),
		2, 5, 5, gconv.Duration("1m"))
	// 反撤回
	notificationGroupId := service.Group().GetMessageNotificationGroupId(ctx, groupId)
	var msg string
	if notificationGroupId < 1 {
		notificationGroupId = groupId
		msg = gconv.String(userId) + "(" + nickname + ") 撤回了一条消息：\n"
	} else {
		msg = gconv.String(userId) + "(" + nickname + ") 在 group(" + gconv.String(groupId) + ") 撤回了一条消息：\n"
	}
	msg += message
	service.Bot().SendMessage(ctx,
		"group", 0, notificationGroupId, msg, false)
	catch = true
	return
}
