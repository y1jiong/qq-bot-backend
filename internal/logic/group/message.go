package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	messageNotificationGroupIdKey = "messageNotificationGroupId"
	antiRecallKey                 = "antiRecall"
	antiRecallOnlyMemberKey       = "antiRecallOnlyMember"
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
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	enabled, _ = settingJson.Get(antiRecallKey).Bool()
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
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	notificationGroupId, _ = settingJson.Get(messageNotificationGroupIdKey).StrictInt64()
	return
}

func (s *sGroup) IsSetOnlyAntiRecallMember(ctx context.Context, groupId int64) (set bool) {
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
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	set, _ = settingJson.Get(antiRecallOnlyMemberKey).Bool()
	return
}
