package group

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	approvalPolicyMapKey           = "approvalPolicy"
	approvalRegexpKey              = "approvalRegexp"
	approvalLevelKey               = "approvalLevel"
	approvalReasonKey              = "approvalReason"
	approvalWhitelistsMapKey       = "approvalWhitelists"
	approvalBlacklistsMapKey       = "approvalBlacklists"
	approvalNotifyOnlyEnabledKey   = "approvalNotifyOnlyEnabled"
	approvalAutoPassDisabledKey    = "approvalAutoPassDisabled"
	approvalAutoRejectDisabledKey  = "approvalAutoRejectDisabled"
	approvalNotificationGroupIdKey = "approvalNotificationGroupId"
)

func (s *sGroup) GetApprovalPolicy(ctx context.Context, groupId int64) (policy map[string]any) {
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
	policy, _ = settingJson.Get(approvalPolicyMapKey).Map()
	if policy == nil {
		policy = make(map[string]any)
	}
	return
}

func (s *sGroup) GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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
	whitelists, _ = settingJson.Get(approvalWhitelistsMapKey).Map()
	if whitelists == nil {
		whitelists = make(map[string]any)
	}
	return
}

func (s *sGroup) GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists, _ = settingJson.Get(approvalBlacklistsMapKey).Map()
	if blacklists == nil {
		blacklists = make(map[string]any)
	}
	return
}

func (s *sGroup) GetApprovalRegexp(ctx context.Context, groupId int64) (exp string) {
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
	exp, _ = settingJson.Get(approvalRegexpKey).StrictString()
	return
}

func (s *sGroup) GetApprovalNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64) {
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
	notificationGroupId, _ = settingJson.Get(approvalNotificationGroupIdKey).StrictInt64()
	return
}

func (s *sGroup) GetApprovalLevel(ctx context.Context, groupId int64) (level int64) {
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
	level, _ = settingJson.Get(approvalLevelKey).StrictInt64()
	return
}

func (s *sGroup) GetApprovalReason(ctx context.Context, groupId int64) (reason string) {
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
	reason, _ = settingJson.Get(approvalReasonKey).StrictString()
	return
}

func (s *sGroup) IsApprovalNotifyOnlyEnabled(ctx context.Context, groupId int64) bool {
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
	b, _ := settingJson.Get(approvalNotifyOnlyEnabledKey).Bool()
	return b
}

func (s *sGroup) IsApprovalAutoPassEnabled(ctx context.Context, groupId int64) bool {
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
	b, _ := settingJson.Get(approvalAutoPassDisabledKey).Bool()
	return !b
}

func (s *sGroup) IsApprovalAutoRejectEnabled(ctx context.Context, groupId int64) bool {
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
	b, _ := settingJson.Get(approvalAutoRejectDisabledKey).Bool()
	return !b
}
