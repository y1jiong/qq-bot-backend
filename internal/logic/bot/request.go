package bot

import (
	"context"
	"qq-bot-backend/internal/service"
)

func processRequest(ctx context.Context) {
	switch service.Bot().GetRequestType(ctx) {
	case "group":
		switch service.Bot().GetSubType(ctx) {
		case "add":
			// 申请入群
			service.Module().TryApproveAddGroup(ctx)
		case "invite":
			// 群邀请
		}
	case "friend":
		// 好友请求
	}
}
