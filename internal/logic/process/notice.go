package process

import (
	"context"
	"qq-bot-backend/internal/service"
)

func processNotice(ctx context.Context) {
	switch service.Bot().GetNoticeType(ctx) {
	case "group_recall":
		// 群消息撤回
		go service.Event().TryCascadingRecall(ctx)
		go service.Event().TryUndoMessageRecall(ctx)
	case "group_increase":
		// 群成员增加
		go service.Event().TryAutoSetCard(ctx)
	case "group_decrease":
		// 群成员减少
		go service.Event().TryLogLeave(ctx)
	case "group_card":
		// 群名片变更
		go service.Event().TryLockCard(ctx)
	case "group_msg_emoji_like":
		// 群消息表情回应（仅自己的消息）
		go service.Event().TryEmojiRecall(ctx)
	case "group_upload":
		// 群文件上传
	case "group_admin":
		// 群管理员变更
	case "group_ban":
		// 群成员禁言
	case "friend_add":
		// 好友添加
	case "friend_recall":
		// 好友消息撤回
	case "offline_file":
		// 离线文件上传
	case "client_status":
		// 客户端状态变更
	case "essence":
		// 精华消息
	case "notify":
		// 系统通知
	}
}
