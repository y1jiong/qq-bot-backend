package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryLogLeave(ctx context.Context) (catch bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryLogLeave")
	defer span.End()

	// 获取当前 group log leave list
	groupId := service.Bot().GetGroupId(ctx)
	listName := service.Group().GetLogLeaveList(ctx, groupId)
	// 预处理
	if listName == "" {
		// 没有设置 leave list
		return
	}

	catch = true

	// 处理
	action := service.Bot().GetSubType(ctx)
	userId := service.Bot().GetUserId(ctx)
	operatorId := service.Bot().GetOperatorId(ctx)
	// 初始化数据
	one := struct {
		SubType    string `json:"sub_type"`
		Time       string `json:"time"`
		OperatorId string `json:"operator_id"`
	}{
		SubType:    action,
		Time:       gtime.New(service.Bot().GetTimestamp(ctx)).String(),
		OperatorId: gconv.String(operatorId),
	}
	listMap := make(map[string]any)
	listMap[gconv.String(userId)] = one
	// 保存数据
	_, err := service.List().AppendListData(ctx, listName, listMap)
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
	return
}

func (s *sEvent) TryLogApproval(ctx context.Context) (catch bool) {
	// 获取当前 group log approval list
	groupId := service.Bot().GetGroupId(ctx)
	listName := service.Group().GetLogApprovalList(ctx, groupId)
	// 预处理
	if listName == "" {
		// 没有设置 approval list
		return
	}
	// 处理
	catch = true
	userId := service.Bot().GetUserId(ctx)
	// 初始化数据
	one := struct {
		Comment string `json:"comment"`
		Time    string `json:"time"`
	}{
		Comment: service.Bot().GetComment(ctx),
		Time:    gtime.New(service.Bot().GetTimestamp(ctx)).String(),
	}
	listMap := make(map[string]any)
	listMap[gconv.String(userId)] = one
	// 保存数据
	_, err := service.List().AppendListData(ctx, listName, listMap)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}
