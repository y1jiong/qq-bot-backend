package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
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
	if groupE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(groupE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listName = settingJson.Get(cardAutoSetListKey).MustString()
	return
}

func (s *sGroup) IsCardLocked(ctx context.Context, groupId int64) (locked bool) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(groupE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	locked = settingJson.Get(cardLockKey).MustBool()
	return
}
