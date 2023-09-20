package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	approvalProcessMapKey          = "approvalProcess"
	approvalRegexpKey              = "approvalRegexp"
	approvalWhitelistsMapKey       = "approvalWhitelists"
	approvalBlacklistsMapKey       = "approvalBlacklists"
	approvalDisabledAutoPassKey    = "approvalDisabledAutoPass"
	approvalDisabledAutoRejectKey  = "approvalDisabledAutoReject"
	approvalNotificationGroupIdKey = "approvalNotificationGroupId"
)

func (s *sGroup) GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any) {
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
	process, _ = settingJson.Get(approvalProcessMapKey).Map()
	if process == nil {
		process = make(map[string]any)
	}
	return
}

func (s *sGroup) GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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
	whitelists, _ = settingJson.Get(approvalWhitelistsMapKey).Map()
	if whitelists == nil {
		whitelists = make(map[string]any)
	}
	return
}

func (s *sGroup) GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists, _ = settingJson.Get(approvalBlacklistsMapKey).Map()
	if blacklists == nil {
		blacklists = make(map[string]any)
	}
	return
}

func (s *sGroup) GetApprovalRegexp(ctx context.Context, groupId int64) (exp string) {
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
	exp, _ = settingJson.Get(approvalRegexpKey).StrictString()
	return
}

func (s *sGroup) GetApprovalNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64) {
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
	notificationGroupId, _ = settingJson.Get(approvalNotificationGroupIdKey).StrictInt64()
	return
}

func (s *sGroup) IsEnabledApprovalAutoPass(ctx context.Context, groupId int64) (enabled bool) {
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
	enabled, _ = settingJson.Get(approvalDisabledAutoPassKey).Bool()
	enabled = !enabled
	return
}

func (s *sGroup) IsEnabledApprovalAutoReject(ctx context.Context, groupId int64) (enabled bool) {
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
	enabled, _ = settingJson.Get(approvalDisabledAutoRejectKey).Bool()
	enabled = !enabled
	return
}
