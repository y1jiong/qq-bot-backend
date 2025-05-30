package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
)

func (s *sGroup) SetAntiRecallReturnRes(ctx context.Context, groupId int64, enable bool) (retMsg string) {
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
	if enable {
		if _, ok := settingJson.CheckGet(antiRecallEnabledKey); ok {
			retMsg = "早已启用 group(" + gconv.String(groupId) + ") 反撤回"
			return
		}
		settingJson.Set(antiRecallEnabledKey, true)
	} else {
		if _, ok := settingJson.CheckGet(antiRecallEnabledKey); ok {
			settingJson.Del(antiRecallEnabledKey)
		} else {
			retMsg = "并未启用 group(" + gconv.String(groupId) + ") 反撤回"
			return
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
	if enable {
		retMsg = "已启用 group(" + gconv.String(groupId) + ") 反撤回"
	} else {
		retMsg = "已禁用 group(" + gconv.String(groupId) + ") 反撤回"
	}
	return
}

func (s *sGroup) SetMessageNotificationReturnRes(ctx context.Context,
	groupId, notificationGroupId int64,
) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 || notificationGroupId == 0 {
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
	if v, ok := settingJson.CheckGet(messageNotificationGroupIdKey); ok {
		retMsg = "早已设置 group(" + gconv.String(groupId) + ") 群消息通知群为 group(" +
			gconv.String(v.MustInt64()) + ")"
		return
	}
	// 验证是否存在该群
	if _, err = service.Bot().GetGroupInfo(ctx, notificationGroupId); err != nil {
		retMsg = "group(" + gconv.String(notificationGroupId) + ") 未找到"
		return
	}
	settingJson.Set(messageNotificationGroupIdKey, notificationGroupId)
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
	retMsg = "已设置 group(" + gconv.String(groupId) + ") 群消息通知群为 group(" + gconv.String(notificationGroupId) + ")"
	return
}

func (s *sGroup) RemoveMessageNotificationReturnRes(ctx context.Context, groupId int64) (retMsg string) {
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
	if _, ok := settingJson.CheckGet(messageNotificationGroupIdKey); !ok {
		retMsg = "并未设置 group(" + gconv.String(groupId) + ") 群消息通知群"
		return
	}
	settingJson.Del(messageNotificationGroupIdKey)
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
	retMsg = "已移除 group(" + gconv.String(groupId) + ") 群消息通知群"
	return
}

func (s *sGroup) SetOnlyAntiRecallMemberReturnRes(ctx context.Context, groupId int64, enable bool) (retMsg string) {
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
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if enable {
		if v, _ := settingJson.Get(antiRecallOnlyMemberKey).Bool(); v {
			retMsg = "早已设置 group(" + gconv.String(groupId) + ") 仅反撤回群成员"
			return
		}
		_, _ = settingJson.Set(antiRecallOnlyMemberKey, ast.NewBool(enable))
	} else {
		if settingJson.Get(antiRecallOnlyMemberKey).Exists() {
			_, _ = settingJson.Unset(antiRecallOnlyMemberKey)
		} else {
			retMsg = "并未设置 group(" + gconv.String(groupId) + ") 仅反撤回群成员"
			return
		}
	}
	// 保存数据
	settingStr, err := settingJson.Raw()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, settingStr).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if enable {
		retMsg = "已设置 group(" + gconv.String(groupId) + ") 仅反撤回群成员"
	} else {
		retMsg = "已取消 group(" + gconv.String(groupId) + ") 仅反撤回群成员"
	}
	return
}
