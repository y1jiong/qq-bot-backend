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

func (s *sGroup) AddKeywordPolicyReturnRes(ctx context.Context,
	groupId int64, policyName string, args ...string) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(groupE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if len(args) > 0 {
		// 处理 args
		switch policyName {
		case consts.ReplyCmd:
			// 添加回复列表
			// 是否存在 list
			lists := service.Namespace().GetNamespaceListsIncludingGlobal(ctx, groupE.Namespace)
			if _, ok := lists[args[0]]; !ok {
				retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + args[0] + ")"
				return
			}
			// 继续处理
			replyLists := settingJson.Get(keywordReplyListsMapKey).MustMap(make(map[string]any))
			replyLists[args[0]] = nil
			settingJson.Set(keywordReplyListsMapKey, replyLists)
		case consts.BlacklistCmd:
			// 添加一个黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
			if _, ok := lists[args[0]]; !ok {
				retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + args[0] + ")"
				return
			}
			// 继续处理
			blacklists := settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(keywordBlacklistsMapKey, blacklists)
		case consts.WhitelistCmd:
			// 添加一个白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
			if _, ok := lists[args[0]]; !ok {
				retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + args[0] + ")"
				return
			}
			// 继续处理
			whitelists := settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(keywordWhitelistsMapKey, whitelists)
		}
	} else {
		// 添加 policyName
		policyMap := settingJson.Get(keywordPolicyMapKey).MustMap(make(map[string]any))
		policyMap[policyName] = nil
		settingJson.Set(keywordPolicyMapKey, policyMap)
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
		retMsg = "已添加 group(" + gconv.String(groupId) + ") 关键词检查 " + policyName + "(" + args[0] + ")"
	} else {
		retMsg = "已启用 group(" + gconv.String(groupId) + ") 关键词检查 " + policyName
	}
	return
}

func (s *sGroup) RemoveKeywordPolicyReturnRes(ctx context.Context,
	groupId int64, policyName string, args ...string) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(groupE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if len(args) > 0 {
		// 处理 args
		switch policyName {
		case consts.ReplyCmd:
			// 移除回复列表
			replyLists := settingJson.Get(keywordReplyListsMapKey).MustMap(make(map[string]any))
			if _, ok := replyLists[args[0]]; !ok {
				retMsg = "在 " + consts.ReplyCmd + " 中未找到 list(" + args[0] + ")"
				return
			}
			delete(replyLists, args[0])
			settingJson.Set(keywordReplyListsMapKey, replyLists)
		case consts.BlacklistCmd:
			// 移除某个黑名单
			blacklists := settingJson.Get(keywordBlacklistsMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; !ok {
				retMsg = "在 " + consts.BlacklistCmd + " 中未找到 list(" + args[0] + ")"
				return
			}
			delete(blacklists, args[0])
			settingJson.Set(keywordBlacklistsMapKey, blacklists)
		case consts.WhitelistCmd:
			// 移除某个白名单
			whitelists := settingJson.Get(keywordWhitelistsMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; !ok {
				retMsg = "在 " + consts.WhitelistCmd + " 中未找到 list(" + args[0] + ")"
				return
			}
			delete(whitelists, args[0])
			settingJson.Set(keywordWhitelistsMapKey, whitelists)
		}
	} else {
		// 删除 policyName
		policyMap := settingJson.Get(keywordPolicyMapKey).MustMap(make(map[string]any))
		if _, ok := policyMap[policyName]; !ok {
			retMsg = "在 " + approvalPolicyMapKey + " 中未找到 " + policyName
			return
		}
		delete(policyMap, policyName)
		settingJson.Set(keywordPolicyMapKey, policyMap)
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
		retMsg = "已移除 group(" + gconv.String(groupId) + ") 关键词检查 " + policyName + "(" + args[0] + ")"
	} else {
		retMsg = "已禁用 group(" + gconv.String(groupId) + ") 关键词检查 " + policyName
	}
	return
}
