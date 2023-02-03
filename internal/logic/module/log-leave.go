package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sModule) TryLogLeave(ctx context.Context) (catch bool) {
	// 获取当前 group log leave list
	groupId := service.Bot().GetGroupId(ctx)
	listName := service.Group().GetLogLeaveList(ctx, groupId)
	// 预处理
	if listName == "" {
		// 没有设置离群记录 list
		return
	}
	// 初始化数据
	one := struct {
		SubType string `json:"subType"`
		Time    string `json:"time"`
	}{
		SubType: service.Bot().GetSubType(ctx),
		Time:    gtime.New(service.Bot().GetTimestamp(ctx)).String(),
	}
	listMap := make(map[string]any)
	listMap[gconv.String(service.Bot().GetUserId(ctx))] = one
	// 保存数据
	err := service.List().AppendListData(ctx, listName, listMap)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	catch = true
	return
}
