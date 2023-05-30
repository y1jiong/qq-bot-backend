package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
	"regexp"
)

func (s *sGroup) SetAutoSetListWithRes(ctx context.Context, groupId int64, listName string) {
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
	// 是否存在 list
	lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
	if _, ok := lists[listName]; !ok {
		service.Bot().SendPlainMsg(ctx, "在 namespace("+gEntity.Namespace+") 中未找到 list("+listName+")")
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
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
	service.Bot().SendPlainMsg(ctx, "已设置 group("+gconv.String(groupId)+") 群名片自动设置 list("+listName+")")
	return
}

func (s *sGroup) RemoveAutoSetListWithRes(ctx context.Context, groupId int64) {
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
	if _, ok := settingJson.CheckGet(cardAutoSetListKey); !ok {
		service.Bot().SendPlainMsg(ctx, "并未设置群名片自动设置 list")
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
	service.Bot().SendPlainMsg(ctx, "已移除 group("+gconv.String(groupId)+") 群名片自动设置 list")
}

func (s *sGroup) CheckCardWithRegexpWithRes(ctx context.Context, groupId int64, listName, exp string) {
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
	// 是否存在 list
	lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
	if _, ok := lists[listName]; !ok {
		service.Bot().SendPlainMsg(ctx, "在 namespace("+gEntity.Namespace+") 中未找到 list("+listName+")")
		return
	}
	// callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
			return
		}
		// compile regexp
		exp = service.Codec().DecodeCqCode(exp)
		reg, err := regexp.Compile(exp)
		if err != nil {
			service.Bot().SendPlainMsg(ctx, "正则表达式编译失败")
			return
		}
		// 获取群成员列表
		membersJson := service.Bot().GetData(rsyncCtx)
		if membersJson == nil {
			// 空列表
			service.Bot().SendPlainMsg(ctx, "获取到空的群成员列表")
			return
		}
		// 局部变量
		membersArr := membersJson.MustArray()
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
		service.Bot().SendPlainMsg(ctx,
			"已将 group("+gconv.String(groupId)+") 中不符合群名片规则的 member 导出到 list("+listName+") "+
				gconv.String(len(membersMap))+" 条\n共 "+gconv.String(totalLen)+" 条")
	}
	// 异步获取群成员列表
	service.Bot().GetGroupMemberList(ctx, groupId, callback, true)
}

func (s *sGroup) CheckCardByListWithRes(ctx context.Context, groupId int64, toList, fromList string) {
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
	// 是否存在 list
	lists := service.Namespace().GetNamespaceList(ctx, gEntity.Namespace)
	if _, ok := lists[toList]; !ok {
		service.Bot().SendPlainMsg(ctx, "在 namespace("+gEntity.Namespace+") 中未找到 list("+toList+")")
		return
	}
	if _, ok := lists[fromList]; !ok {
		service.Bot().SendPlainMsg(ctx, "在 namespace("+gEntity.Namespace+") 中未找到 list("+fromList+")")
		return
	}
	// callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
			return
		}
		// 获取 fromList
		fromListData := service.List().GetListData(ctx, fromList)
		// 获取群成员列表
		membersJson := service.Bot().GetData(rsyncCtx)
		if membersJson == nil {
			// 空列表
			service.Bot().SendPlainMsg(ctx, "获取到空的群成员列表")
			return
		}
		// 局部变量
		membersArr := membersJson.MustArray()
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
		service.Bot().SendPlainMsg(ctx,
			"已将 group("+gconv.String(groupId)+") 中不符合群名片规则的 member 导出到 list("+toList+") "+
				gconv.String(len(membersMap))+" 条\n共 "+gconv.String(totalLen)+" 条")
	}
	// 异步获取群成员列表
	service.Bot().GetGroupMemberList(ctx, groupId, callback, true)
}

func (s *sGroup) LockCardWithRes(ctx context.Context, groupId int64) {
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
	if _, ok := settingJson.CheckGet(cardLockKey); ok {
		service.Bot().SendPlainMsg(ctx, "group("+gconv.String(groupId)+") 群名片已锁定")
		return
	}
	settingJson.Set(cardLockKey, true)
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
	service.Bot().SendPlainMsg(ctx, "已设置 group("+gconv.String(groupId)+") 群名片锁定")
}

func (s *sGroup) UnlockCardWithRes(ctx context.Context, groupId int64) {
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
	if _, ok := settingJson.CheckGet(cardLockKey); !ok {
		service.Bot().SendPlainMsg(ctx, "group("+gconv.String(groupId)+") 群名片未锁定")
		return
	}
	settingJson.Del(cardLockKey)
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
	service.Bot().SendPlainMsg(ctx, "已设置 group("+gconv.String(groupId)+") 群名片解锁")
}
