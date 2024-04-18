package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
	"regexp"
)

func (s *sGroup) AddApprovalProcessReturnRes(ctx context.Context,
	groupId int64, processName string, args ...string) (retMsg string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) {
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
		switch processName {
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
			if service.Codec().IsIncludeCqCode(args[0]) {
				// 包含 CQ Code 时发送表情 gun
				service.Bot().SendMsg(ctx, "[CQ:face,id=288]")
				return
			}
			// 解码被 CQ Code 转义的字符
			args[0] = service.Codec().DecodeCqCode(args[0])
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
		switch processName {
		case consts.AutoPassCmd:
			if _, ok := settingJson.CheckGet(approvalDisabledAutoPassKey); !ok {
				retMsg = "并未禁用自动通过"
				return
			}
			settingJson.Del(approvalDisabledAutoPassKey)
		case consts.AutoRejectCmd:
			if _, ok := settingJson.CheckGet(approvalDisabledAutoRejectKey); !ok {
				retMsg = "并未禁用自动拒绝"
				return
			}
			settingJson.Del(approvalDisabledAutoRejectKey)
		default:
			// 添加 processName
			processMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
			processMap[processName] = nil
			settingJson.Set(approvalProcessMapKey, processMap)
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
		retMsg = "已添加 group(" + gconv.String(groupId) + ") 入群审核 " + processName + "(" + args[0] + ")"
	} else {
		retMsg = "已启用 group(" + gconv.String(groupId) + ") 入群审核 " + processName
	}
	return
}

func (s *sGroup) RemoveApprovalProcessReturnRes(ctx context.Context,
	groupId int64, processName string, args ...string) (retMsg string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) {
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
		switch processName {
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
		switch processName {
		case consts.AutoPassCmd:
			if _, ok := settingJson.CheckGet(approvalDisabledAutoPassKey); ok {
				retMsg = "早已禁用自动通过"
				return
			}
			settingJson.Set(approvalDisabledAutoPassKey, true)
		case consts.AutoRejectCmd:
			if _, ok := settingJson.CheckGet(approvalDisabledAutoRejectKey); ok {
				retMsg = "早已禁用自动拒绝"
				return
			}
			settingJson.Set(approvalDisabledAutoRejectKey, true)
		case consts.NotificationCmd:
			if _, ok := settingJson.CheckGet(approvalNotificationGroupIdKey); !ok {
				retMsg = "并未设置通知群"
				return
			}
			settingJson.Del(approvalNotificationGroupIdKey)
		default:
			// 删除 processName
			processMap := settingJson.Get(approvalProcessMapKey).MustMap(make(map[string]any))
			if _, ok := processMap[processName]; !ok {
				retMsg = "在 " + approvalProcessMapKey + " 中未找到 " + processName
				return
			}
			delete(processMap, processName)
			settingJson.Set(approvalProcessMapKey, processMap)
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
		retMsg = "已移除 group(" + gconv.String(groupId) + ") 入群审核 " + processName + "(" + args[0] + ")"
	} else {
		retMsg = "已禁用 group(" + gconv.String(groupId) + ") 入群审核 " + processName
	}
	return
}
