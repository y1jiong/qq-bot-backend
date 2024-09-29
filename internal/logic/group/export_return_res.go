package group

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"strings"
)

func (s *sGroup) ExportGroupMemberListReturnRes(ctx context.Context,
	groupId int64, listName string) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 是否存在 list
	lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
	if _, ok := lists[listName]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + listName + ")"
		return
	}
	// 获取群成员列表
	membersArr, err := service.Bot().GetGroupMemberList(ctx, groupId, true)
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
			infoM := make(map[string]any)
			for k, vvv := range vv {
				// 过滤字段
				switch k {
				case "sex", "user_id", "card_changeable":
					continue
				}
				// 过滤空值
				if gvar.New(vvv).IsEmpty() {
					continue
				}
				if str, okay := vvv.(string); okay && str == "0" {
					continue
				}
				// 时间戳转换
				if strings.HasSuffix(k, "_time") || strings.HasSuffix(k, "_timestamp") {
					vvv = gtime.New(gconv.Int64(vvv)).String()

					k = strings.TrimSuffix(k, "_time")
					k = strings.TrimSuffix(k, "_timestamp")
				}

				infoM[k] = vvv
			}
			// 写入数据
			membersMap[gconv.String(vv["user_id"])] = infoM
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
