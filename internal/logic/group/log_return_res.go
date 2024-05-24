package group

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
)

func (s *sGroup) SetLogLeaveListReturnRes(ctx context.Context,
	groupId int64, listName string) (retMsg string) {
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
	// 是否存在 list
	lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
	if _, ok := lists[listName]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + listName + ")"
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	_, _ = settingJson.Set(logLeaveListKey, ast.NewString(listName))
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
	retMsg = "已设置 group(" + gconv.String(groupId) + ") 记录离群 list(" + listName + ")"
	return
}

func (s *sGroup) RemoveLogLeaveListReturnRes(ctx context.Context, groupId int64) (retMsg string) {
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
	if !settingJson.Get(logLeaveListKey).Exists() {
		retMsg = "并未设置记录离群 list"
		return
	}
	_, _ = settingJson.Unset(logLeaveListKey)
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
	retMsg = "已移除 group(" + gconv.String(groupId) + ") 记录离群 list"
	return
}

func (s *sGroup) SetLogApprovalListReturnRes(ctx context.Context,
	groupId int64, listName string) (retMsg string) {
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
	// 是否存在 list
	lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
	if _, ok := lists[listName]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + listName + ")"
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(groupE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	_, _ = settingJson.Set(logApprovalListKey, ast.NewString(listName))
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
	retMsg = "已设置 group(" + gconv.String(groupId) + ") 记录入群审核 list(" + listName + ")"
	return
}

func (s *sGroup) RemoveLogApprovalListReturnRes(ctx context.Context, groupId int64) (retMsg string) {
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
	if !settingJson.Get(logApprovalListKey).Exists() {
		retMsg = "并未设置记录入群审核 list"
		return
	}
	_, _ = settingJson.Unset(logApprovalListKey)
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
	retMsg = "已移除 group(" + gconv.String(groupId) + ") 记录入群审核 list"
	return
}
