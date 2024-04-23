package group

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sGroup) ExportGroupMemberListReturnRes(ctx context.Context,
	groupId int64, listName string) (retMsg string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 是否存在 list
	lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
	if _, ok := lists[listName]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + listName + ")"
		return
	}
	// 获取群成员列表
	membersArr, err := service.Bot().GetGroupMemberList(ctx, groupId)
	if err != nil {
		retMsg = "获取群成员列表失败"
		return
	}
	// 局部变量
	membersMap := make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if vv, ok := v.(map[string]any); ok {
			// 写入数据
			membersMap[gconv.String(vv["user_id"])] = struct {
				Card     string `json:"card"`
				JoinTime string `json:"join_time"`
			}{
				Card:     gconv.String(vv["card"]),
				JoinTime: gtime.New(gconv.Int(vv["join_time"])).String(),
			}
		}
	}
	// 保存数据
	totalLen, err := service.List().AppendListData(ctx, listName, membersMap)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 group(" + gconv.String(groupId) + ") 的 member 导出到 list(" + listName + ") " +
		gconv.String(len(membersMap)) + " 条\n共 " + gconv.String(totalLen) + " 条"
	return
}
