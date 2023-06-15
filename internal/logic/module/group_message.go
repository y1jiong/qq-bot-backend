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
	// 异步反撤回
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
			return
		}
		// 获取发送者信息
		nickname, userId := service.Bot().GetSenderFromData(rsyncCtx)
		// 获取撤回消息
		message := service.Bot().GetMessageFromData(rsyncCtx)
		// 反撤回
		msg := gconv.String(userId) + "(" + nickname + ")" + "撤回了一条消息：\n" + message
		service.Bot().SendMsg(ctx, msg)
	}
	service.Bot().RequestMessage(ctx, msgId, callback)
	// 防止过度触发反撤回
	s.AutoMute(ctx, "recall", groupId, service.Bot().GetUserId(ctx),
		2, 5, 5, gconv.Duration("1m"))
	catch = true
	return
}
