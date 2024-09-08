package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/internal/util/codec"
	"regexp"
)

func (s *sGroup) SetAutoSetListReturnRes(ctx context.Context,
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
	settingJson, err := sj.NewJson([]byte(groupE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	settingJson.Set(cardAutoSetListKey, listName)
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
	retMsg = "已设置 group(" + gconv.String(groupId) + ") 群名片自动设置 list(" + listName + ")"
	return
}

func (s *sGroup) RemoveAutoSetListReturnRes(ctx context.Context, groupId int64) (retMsg string) {
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
	if _, ok := settingJson.CheckGet(cardAutoSetListKey); !ok {
		retMsg = "并未设置群名片自动设置 list"
		return
	}
	settingJson.Del(cardAutoSetListKey)
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
	retMsg = "已移除 group(" + gconv.String(groupId) + ") 群名片自动设置 list"
	return
}

func (s *sGroup) CheckCardWithRegexpReturnRes(ctx context.Context,
	groupId int64, listName, exp string) (retMsg string) {
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
	// 获取群成员列表
	membersArr, err := service.Bot().GetGroupMemberList(ctx, groupId, true)
	if err != nil {
		retMsg = "获取群成员列表失败"
		return
	}
	// compile regexp
	exp = codec.DecodeCqCode(exp)
	reg, err := regexp.Compile(exp)
	if err != nil {
		retMsg = "正则表达式编译失败"
		return
	}
	// 局部变量
	membersMap := make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if vv, ok := v.(map[string]any); ok {
			if vv["role"] != "member" {
				continue
			}
			userCard := gconv.String(vv["card"])
			// 正则匹配
			if !reg.MatchString(userCard) {
				membersMap[gconv.String(vv["user_id"])] = userCard
			}
		}
	}
	// 保存数据
	totalLen, err := service.List().AppendListData(ctx, listName, membersMap)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 group(" + gconv.String(groupId) + ") 中不符合群名片规则的 member 导出到 list(" + listName + ") " +
		gconv.String(len(membersMap)) + " 条\n共 " + gconv.String(totalLen) + " 条"
	return
}

func (s *sGroup) CheckCardByListReturnRes(ctx context.Context,
	groupId int64, toList, fromList string) (retMsg string) {
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
	if _, ok := lists[toList]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + toList + ")"
		return
	}
	if _, ok := lists[fromList]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + fromList + ")"
		return
	}
	// 获取群成员列表
	membersArr, err := service.Bot().GetGroupMemberList(ctx, groupId, true)
	// 获取 fromList
	fromListData := service.List().GetListData(ctx, fromList)
	// 局部变量
	membersMap := make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if vv, ok := v.(map[string]any); ok {
			if vv["role"] != "member" {
				continue
			}
			userId := gconv.String(vv["user_id"])
			// 匹配
			if _, okk := fromListData[userId].(string); okk {
				userCard := gconv.String(vv["card"])
				if userCard != fromListData[userId] {
					membersMap[userId] = userCard
				}
			}
		}
	}
	// 保存数据
	totalLen, err := service.List().AppendListData(ctx, toList, membersMap)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 group(" + gconv.String(groupId) + ") 中不符合群名片规则的 member 导出到 list(" + toList + ") " +
		gconv.String(len(membersMap)) + " 条\n共 " + gconv.String(totalLen) + " 条"
	return
}

func (s *sGroup) LockCardReturnRes(ctx context.Context, groupId int64) (retMsg string) {
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
	if _, ok := settingJson.CheckGet(cardLockEnabledKey); ok {
		retMsg = "group(" + gconv.String(groupId) + ") 群名片已锁定"
		return
	}
	settingJson.Set(cardLockEnabledKey, true)
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
	retMsg = "已设置 group(" + gconv.String(groupId) + ") 群名片锁定"
	return
}

func (s *sGroup) UnlockCardReturnRes(ctx context.Context, groupId int64) (retMsg string) {
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
	if _, ok := settingJson.CheckGet(cardLockEnabledKey); !ok {
		retMsg = "group(" + gconv.String(groupId) + ") 群名片未锁定"
		return
	}
	settingJson.Del(cardLockEnabledKey)
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
	retMsg = "已设置 group(" + gconv.String(groupId) + ") 群名片解锁"
	return
}
