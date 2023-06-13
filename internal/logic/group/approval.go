package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	approvalProcessMapKey       = "approvalProcess"
	approvalRegexpKey           = "approvalRegexp"
	approvalWhitelistsMapKey    = "approvalWhitelists"
	approvalBlacklistsMapKey    = "approvalBlacklists"
	approvalDisabledAutoPassKey = "approvalDisabledAutoPass"
)

func (s *sGroup) GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	process = settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	whitelists = settingJson.Get(approvalWhitelistsMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	blacklists = settingJson.Get(approvalBlacklistsMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetApprovalRegexp(ctx context.Context, groupId int64) (exp string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	exp = settingJson.Get(approvalRegexpKey).MustString()
	return
}

func (s *sGroup) IsEnabledApprovalAutoPass(ctx context.Context, groupId int64) (enabled bool) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	enabled = !settingJson.Get(approvalDisabledAutoPassKey).MustBool()
	return
}
