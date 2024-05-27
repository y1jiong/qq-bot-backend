package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	messageNotificationGroupIdKey = "messageNotificationGroupId"
	antiRecallEnabledKey          = "antiRecallEnabled"
	antiRecallOnlyMemberKey       = "antiRecallOnlyMember"
)

func (s *sGroup) IsAntiRecallEnabled(ctx context.Context, groupId int64) bool {
	// 参数合法性校验
	if groupId == 0 {
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
	b, _ := settingJson.Get(antiRecallEnabledKey).Bool()
	return b
}

func (s *sGroup) GetMessageNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64) {
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
	notificationGroupId, _ = settingJson.Get(messageNotificationGroupIdKey).StrictInt64()
	return
}

func (s *sGroup) IsSetOnlyAntiRecallMember(ctx context.Context, groupId int64) bool {
	// 参数合法性校验
	if groupId == 0 {
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
	b, _ := settingJson.Get(antiRecallOnlyMemberKey).Bool()
	return b
}
