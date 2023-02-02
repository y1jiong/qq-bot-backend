package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
)

const (
	keywordProcessMapKey    = "keywordProcess"
	keywordWhitelistsMapKey = "keywordWhitelists"
	keywordBlacklistsMapKey = "keywordBlacklists"
)

func (s *sGroup) AddKeywordProcess(ctx context.Context, groupId int64, processName string, args ...string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if len(args) > 0 {
		// 处理 args
		switch processName {
		case consts.BlacklistCmd:
			// 处理黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			blacklists := settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(keywordBlacklistsMapKey, blacklists)
		case consts.WhitelistCmd:
			// 处理白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
			if _, ok := lists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			// 继续处理
			whitelists := settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(keywordWhitelistsMapKey, whitelists)
		}
	} else {
		// 添加 processName
		processMap := settingJson.Get(keywordProcessMapKey).MustMap(make(map[string]any))
		processMap[processName] = nil
		settingJson.Set(keywordProcessMapKey, processMap)
	}
	// 保存数据
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if len(args) > 0 {
		service.Bot().SendPlainMsg(ctx,
			"已添加 group("+gconv.String(groupId)+") 关键词检查流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已启用 group("+gconv.String(groupId)+") 关键词检查流程 "+processName)
	}
}

func (s *sGroup) RemoveKeywordProcess(ctx context.Context, groupId int64, processName string, args ...string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if len(args) > 0 {
		// 处理 args
		switch processName {
		case consts.BlacklistCmd:
			// 处理黑名单
			blacklists := settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(blacklists, args[0])
			settingJson.Set(keywordBlacklistsMapKey, blacklists)
		case consts.WhitelistCmd:
			// 处理白名单
			whitelists := settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; !ok {
				service.Bot().SendPlainMsg(ctx, args[0]+" 不存在")
				return
			}
			delete(whitelists, args[0])
			settingJson.Set(keywordWhitelistsMapKey, whitelists)
		}
	} else {
		// 删除 processName
		processMap := settingJson.Get(keywordProcessMapKey).MustMap(make(map[string]any))
		if _, ok := processMap[processName]; !ok {
			service.Bot().SendPlainMsg(ctx, processName+" 不存在")
			return
		}
		delete(processMap, processName)
		settingJson.Set(keywordProcessMapKey, processMap)
	}
	// 保存数据
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if len(args) > 0 {
		service.Bot().SendPlainMsg(ctx,
			"已移除 group("+gconv.String(groupId)+") 关键词检查流程 "+processName+"("+args[0]+")")
	} else {
		service.Bot().SendPlainMsg(ctx, "已禁用 group("+gconv.String(groupId)+") 关键词检查流程 "+processName)
	}
}

func (s *sGroup) GetKeywordProcess(ctx context.Context, groupId int64) (process map[string]any) {
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
	process = settingJson.Get(keywordProcessMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any) {
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
	whitelists = settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
	return
}

func (s *sGroup) GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any) {
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
	blacklists = settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
	return
}
