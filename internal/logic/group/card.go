package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	cardAutoSetListKey = "cardAutoSetList"
	cardLockKey        = "cardLock"
)

func (s *sGroup) GetCardAutoSetList(ctx context.Context, groupId int64) (listName string) {
	// 参数合法性校验
	if groupId < 1 {
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
	listName, _ = settingJson.Get(cardAutoSetListKey).StrictString()
	return
}

func (s *sGroup) IsCardLocked(ctx context.Context, groupId int64) bool {
	// 参数合法性校验
	if groupId < 1 {
		return false
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return false
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	locked, _ := settingJson.Get(cardLockKey).Bool()
	return locked
}
