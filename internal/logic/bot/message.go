package bot

import (
	"context"
	"qq-bot-backend/internal/service"
)

func processMessage(ctx context.Context) {
	msgType, subType := service.Bot().GetMsgType(ctx), service.Bot().GetSubType(ctx)
	switch msgType {
	case "group":
		// 群消息
		switch subType {
		case "normal":
			// 群聊
		case "anonymous":
			// 匿名
		case "notice":
			// 系统提示(可能是生日提醒之类的消息)
		}
	case "private":
		// 私聊消息
		switch subType {
		case "group":
			// 群临时会话
		case "friend":
			// 好友私聊
		}
	}
}
