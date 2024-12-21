package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility/codec"
	"regexp"
)

func (s *sGroup) AddApprovalPolicyReturnRes(ctx context.Context,
	groupId int64, policyName string, args ...string,
) (retMsg string) {
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
		case consts.WhitelistCmd:
			// 处理白名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
			if _, ok := lists[args[0]]; !ok {
				retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + args[0] + ")"
				return
			}
			// 继续处理
			whitelists := settingJson.Get(approvalWhitelistsMapKey).MustMap(make(map[string]any))
			whitelists[args[0]] = nil
			settingJson.Set(approvalWhitelistsMapKey, whitelists)
		case consts.BlacklistCmd:
			// 处理黑名单
			// 是否存在 list
			lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
			if _, ok := lists[args[0]]; !ok {
				retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + args[0] + ")"
				return
			}
			// 继续处理
			blacklists := settingJson.Get(approvalBlacklistsMapKey).MustMap(make(map[string]any))
			blacklists[args[0]] = nil
			settingJson.Set(approvalBlacklistsMapKey, blacklists)
		case consts.RegexpCmd:
			if codec.IsIncludeCQCode(args[0]) {
				// 包含 CQ Code 时发送表情 gun
				service.Bot().SendMsgIfNotApiReq(ctx, "[CQ:face,id=288]", true)
				return
			}
			// 解码被 CQ Code 转义的字符
			args[0] = codec.DecodeCQCode(args[0])
			// 处理正则表达式
			_, err = regexp.Compile(args[0])
			if err != nil {
				retMsg = "输入的正则表达式无法通过编译"
				return
			}
			settingJson.Set(approvalRegexpKey, args[0])
		case consts.NotificationCmd:
			if v, ok := settingJson.CheckGet(approvalNotificationGroupIdKey); ok {
				retMsg = "早已设置 group(" + gconv.String(groupId) + ") 群入群审核通知群为 group(" +
					gconv.String(v.MustInt64()) + ")"
				return
			}
			// 验证是否存在该群
			_, err = service.Bot().GetGroupInfo(ctx, gconv.Int64(args[0]))
			if err != nil {
				retMsg = "group(" + args[0] + ") 未找到"
				return
			}
			// 继续处理
			settingJson.Set(approvalNotificationGroupIdKey, gconv.Int64(args[0]))
		}
	} else {
		switch policyName {
		case consts.NotifyOnlyCmd:
			if _, ok := settingJson.CheckGet(approvalNotifyOnlyEnabledKey); ok {
				retMsg = "早已启用仅通知"
				return
			}
			settingJson.Set(approvalNotifyOnlyEnabledKey, true)
		case consts.AutoPassCmd:
			if _, ok := settingJson.CheckGet(approvalAutoPassDisabledKey); !ok {
				retMsg = "并未禁用自动通过"
				return
			}
			settingJson.Del(approvalAutoPassDisabledKey)
		case consts.AutoRejectCmd:
			if _, ok := settingJson.CheckGet(approvalAutoRejectDisabledKey); !ok {
				retMsg = "并未禁用自动拒绝"
				return
			}
			settingJson.Del(approvalAutoRejectDisabledKey)
		default:
			// 添加 policyName
			policyMap := settingJson.Get(approvalPolicyMapKey).MustMap(make(map[string]any))
			policyMap[policyName] = nil
			settingJson.Set(approvalPolicyMapKey, policyMap)
		}
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
		retMsg = "已添加 group(" + gconv.String(groupId) + ") 入群审核 " + policyName + "(" + args[0] + ")"
	} else {
		retMsg = "已启用 group(" + gconv.String(groupId) + ") 入群审核 " + policyName
	}
	return
}

func (s *sGroup) RemoveApprovalPolicyReturnRes(ctx context.Context,
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
		case consts.WhitelistCmd:
			// 处理白名单
			whitelists := settingJson.Get(approvalWhitelistsMapKey).MustMap(make(map[string]any))
			if _, ok := whitelists[args[0]]; !ok {
				retMsg = "在 " + consts.WhitelistCmd + " 中未找到 list(" + args[0] + ")"
				return
			}
			delete(whitelists, args[0])
			settingJson.Set(approvalWhitelistsMapKey, whitelists)
		case consts.BlacklistCmd:
			// 处理黑名单
			blacklists := settingJson.Get(approvalBlacklistsMapKey).MustMap(make(map[string]any))
			if _, ok := blacklists[args[0]]; !ok {
				retMsg = "在 " + consts.BlacklistCmd + " 中未找到 list(" + args[0] + ")"
				return
			}
			delete(blacklists, args[0])
			settingJson.Set(approvalBlacklistsMapKey, blacklists)
		}
	} else {
		switch policyName {
		case consts.NotifyOnlyCmd:
			if _, ok := settingJson.CheckGet(approvalNotifyOnlyEnabledKey); !ok {
				retMsg = "并未启用仅通知"
				return
			}
			settingJson.Del(approvalNotifyOnlyEnabledKey)
		case consts.AutoPassCmd:
			if _, ok := settingJson.CheckGet(approvalAutoPassDisabledKey); ok {
				retMsg = "早已禁用自动通过"
				return
			}
			settingJson.Set(approvalAutoPassDisabledKey, true)
		case consts.AutoRejectCmd:
			if _, ok := settingJson.CheckGet(approvalAutoRejectDisabledKey); ok {
				retMsg = "早已禁用自动拒绝"
				return
			}
			settingJson.Set(approvalAutoRejectDisabledKey, true)
		case consts.NotificationCmd:
			if _, ok := settingJson.CheckGet(approvalNotificationGroupIdKey); !ok {
				retMsg = "并未设置通知群"
				return
			}
			settingJson.Del(approvalNotificationGroupIdKey)
		default:
			// 删除 policyName
			policyMap := settingJson.Get(approvalPolicyMapKey).MustMap(make(map[string]any))
			if _, ok := policyMap[policyName]; !ok {
				retMsg = "在 " + approvalPolicyMapKey + " 中未找到 " + policyName
				return
			}
			delete(policyMap, policyName)
			settingJson.Set(approvalPolicyMapKey, policyMap)
		}
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
		retMsg = "已移除 group(" + gconv.String(groupId) + ") 入群审核 " + policyName + "(" + args[0] + ")"
	} else {
		retMsg = "已禁用 group(" + gconv.String(groupId) + ") 入群审核 " + policyName
	}
	return
}
