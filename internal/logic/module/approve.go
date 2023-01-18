package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/service"
)

func (s *sModule) TryApproveAddGroup(ctx context.Context) (catch bool) {
	catch = true

	comment := service.Bot().GetComment(ctx)
	// 正版验证
	genuine, name, uuid, err := service.ThirdParty().QueryMinecraftGenuineUser(ctx, comment)
	if err != nil {
		g.Log().Notice(ctx, err)
	}
	// 发送审批请求
	service.Bot().ApproveAddGroup(ctx,
		service.Bot().GetFlag(ctx),
		service.Bot().GetSubType(ctx),
		genuine,
		"invalid name")
	// 打印通过的日志
	if genuine {
		g.Log().Infof(ctx, "approve user(%v) join group(%v) with %v(%v) in %v",
			service.Bot().GetUserId(ctx),
			service.Bot().GetGroupId(ctx),
			name, uuid,
			comment)
	}
	return
}
