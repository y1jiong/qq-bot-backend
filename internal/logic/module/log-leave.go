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
	action := service.Bot().GetSubType(ctx)
	userId := service.Bot().GetUserId(ctx)
	operatorId := service.Bot().GetOperatorId(ctx)
	// 初始化数据
	one := struct {
		SubType    string `json:"subType"`
		Time       string `json:"time"`
		OperatorId string `json:"operatorId"`
	}{
		SubType:    action,
		Time:       gtime.New(service.Bot().GetTimestamp(ctx)).String(),
		OperatorId: gconv.String(operatorId),
	}
	listMap := make(map[string]any)
	listMap[gconv.String(userId)] = one
	// 保存数据
	err := service.List().AppendListData(ctx, listName, listMap)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 打印离群日志
	g.Log().Infof(ctx, "%v user(%v) from group(%v) by user(%v)",
		action,
		userId,
		groupId,
		operatorId)
	catch = true
	return
}
