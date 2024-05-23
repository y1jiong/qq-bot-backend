package event

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryUndoMessageRecall(ctx context.Context) (catch bool) {
	groupId := service.Bot().GetGroupId(ctx)
	if service.Group().IsSetOnlyAntiRecallMember(ctx, groupId) && service.Bot().IsGroupOwnerOrAdmin(ctx) {
		// owner or admin 在 only-anti-recall-member 情况下不需要反撤回
		return
	}
	// 获取当前 group message anti-recall 策略
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
	nickname := gconv.String(senderMap["card"])
	if nickname == "" {
		nickname = gconv.String(senderMap["nickname"])
	}
	userId := gconv.Int64(senderMap["user_id"])
	// 获取撤回消息
	message := gconv.String(messageMap["message"])
	// 防止过度触发反撤回
	service.Util().AutoMute(ctx, "recall", groupId, userId,
		2, 5, 5, gconv.Duration("1m"))
	// 反撤回
	notificationGroupId := service.Group().GetMessageNotificationGroupId(ctx, groupId)
	var msg string
	if notificationGroupId < 1 {
		notificationGroupId = groupId
		msg = "user[" + nickname + "](" + gconv.String(userId) + ") 撤回了：\n"
	} else {
		msg = "user[" + nickname + "](" + gconv.String(userId) +
			") 在 group(" + gconv.String(groupId) + ") 撤回了：\n"
	}
	msg += message
	_ = service.Bot().SendMessage(ctx,
		"group", 0, notificationGroupId, msg, false)
	catch = true
	return
}
