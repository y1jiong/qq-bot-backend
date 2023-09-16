package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	messageNotificationGroupIdKey = "messageNotificationGroupId"
	antiRecallKey                 = "antiRecall"
)

func (s *sGroup) IsEnabledAntiRecall(ctx context.Context, groupId int64) (enabled bool) {
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
	enabled = settingJson.Get(antiRecallKey).MustBool()
	return
}

func (s *sGroup) GetMessageNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64) {
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
	notificationGroupId = settingJson.Get(messageNotificationGroupIdKey).MustInt64()
	return
}
