package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	keywordPolicyMapKey     = "keywordPolicy"
	keywordWhitelistsMapKey = "keywordWhitelists"
	keywordBlacklistsMapKey = "keywordBlacklists"
	keywordReplyListsMapKey = "keywordReplyLists"
)

func (s *sGroup) GetKeywordPolicy(ctx context.Context, groupId int64) (policy map[string]any) {
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
	policy, _ = settingJson.Get(keywordPolicyMapKey).Map()
	if policy == nil {
		policy = make(map[string]any)
	}
	return
}

func (s *sGroup) GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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
	whitelists, _ = settingJson.Get(keywordWhitelistsMapKey).Map()
	if whitelists == nil {
		whitelists = make(map[string]any)
	}
	return
}

func (s *sGroup) GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists, _ = settingJson.Get(keywordBlacklistsMapKey).Map()
	if blacklists == nil {
		blacklists = make(map[string]any)
	}
	return
}

func (s *sGroup) GetKeywordReplyLists(ctx context.Context, groupId int64) (replyLists map[string]any) {
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
	replyLists, _ = settingJson.Get(keywordReplyListsMapKey).Map()
	if replyLists == nil {
		replyLists = make(map[string]any)
	}
	return
}
