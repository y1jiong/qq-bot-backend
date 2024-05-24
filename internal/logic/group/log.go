package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	logLeaveListKey    = "logLeaveList"
	logApprovalListKey = "logApprovalList"
)

func (s *sGroup) GetLogLeaveList(ctx context.Context, groupId int64) (listName string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listName, _ = settingJson.Get(logLeaveListKey).StrictString()
	return
}

func (s *sGroup) GetLogApprovalList(ctx context.Context, groupId int64) (listName string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listName, _ = settingJson.Get(logApprovalListKey).StrictString()
	return
}
