package process

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
			go service.Bot().CacheMessageAstNode(ctx)
			go service.Event().TryForward(ctx)
			go service.Event().TryKeywordRecall(ctx)
			go service.Event().TryGroupKeywordReply(ctx)
		case "anonymous":
			// 匿名
		case "notice":
			// 系统提示(可能是生日提醒之类的消息)
		}
	case "private":
		// 私聊消息
		switch subType {
		case "friend":
			// 好友私聊
			go service.Event().TryForward(ctx)
			go service.Event().TryKeywordReply(ctx)
		case "group":
			// 群临时会话
			go service.Event().TryForward(ctx)
			go service.Event().TryKeywordReply(ctx)
		}
	}
}
