package group

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"sync"
	"time"
)

func (s *sGroup) BindNamespaceReturnRes(ctx context.Context, groupId int64, namespace string) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) ||
		!service.Namespace().IsNamespaceOwnerOrAdmin(ctx, namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	var err error
	if gEntity == nil {
		// 初始化 group 对象
		gEntity = &entity.Group{
			GroupId:     groupId,
			Namespace:   namespace,
			SettingJson: "{}",
		}
		// 数据库插入
		_, err = dao.Group.Ctx(ctx).
			Data(gEntity).
			OmitEmpty().
			Insert()
	} else {
		if gEntity.Namespace != "" {
			service.Bot().SendPlainMsg(ctx,
				"当前 group("+gconv.String(groupId)+") 已经绑定了 namespace("+gEntity.Namespace+")")
			return
		}
		// 重置 setting
		gEntity = &entity.Group{
			Namespace:   namespace,
			SettingJson: "{}",
		}
		// 数据库更新
		_, err = dao.Group.Ctx(ctx).
			Where(dao.Group.Columns().GroupId, groupId).
			Data(gEntity).
			OmitEmpty().
			Update()
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已绑定当前 group("+gconv.String(groupId)+") 到 namespace("+namespace+")")
}

func (s *sGroup) UnbindReturnRes(ctx context.Context, groupId int64) {
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
	if gEntity == nil || gEntity.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据库更新
	_, err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().Namespace, "").
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已解除 group("+gconv.String(groupId)+") 的 namespace 绑定")
}

func (s *sGroup) QueryGroupReturnRes(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
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
	// 回执
	msg := dao.Group.Columns().Namespace + ": " + gEntity.Namespace + "\n" +
		dao.Group.Columns().SettingJson + ": " + gEntity.SettingJson + "\n" +
		dao.Group.Columns().UpdatedAt + ": " + gEntity.UpdatedAt.String()
	service.Bot().SendPlainMsg(ctx, msg)
}

func (s *sGroup) KickFromListReturnRes(ctx context.Context, groupId int64, listName string) {
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
	// 获取 list
	listMap := service.List().GetListData(ctx, listName)
	listMapLen := len(listMap)
	// 异步 callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
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
		kickMap := make(map[string]any)
		// 解析数组
		for _, v := range membersArr {
			// map 断言
			if kk, ok := v.(map[string]any); ok {
				if kk["role"] != "member" {
					continue
				}
				// 放入待踢出 map
				userId := gconv.String(kk["user_id"])
				if _, ok := listMap[userId]; ok {
					kickMap[userId] = nil
				}
			}
		}
		// 不需要执行踢出
		if len(kickMap) == 0 {
			service.Bot().SendPlainMsg(ctx, "list("+listName+") 中的 group("+gconv.String(groupId)+
				") member 无需踢出")
			return
		}
		// 踢人过程
		service.Bot().SendPlainMsg(ctx, "正在踢出 list("+listName+") 中的 group("+gconv.String(groupId)+
			") member\n共 "+gconv.String(listMapLen)+" 条，有 "+gconv.String(len(kickMap))+" 条需要踢出")
		for k := range kickMap {
			// 踢人
			service.Bot().Kick(ctx, groupId, gconv.Int64(k))
			// 随机延时
			time.Sleep(time.Duration(grand.N(1000, 10000)) * time.Millisecond)
		}
		// 异步检查有没有踢出失败的
		callback2 := func(ctx context.Context, rsyncCtx context.Context) {
			if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
				return
			}
			// 获取群成员列表
			membersJson = service.Bot().GetData(rsyncCtx)
			if membersJson == nil {
				// 空列表
				service.Bot().SendPlainMsg(ctx, "获取到空的群成员列表")
				return
			}
			// 局部变量
			membersArr = membersJson.MustArray()
			kickMap = make(map[string]any)
			// 解析数组
			for _, v := range membersArr {
				// map 断言
				if kk, ok := v.(map[string]any); ok {
					if kk["role"] != "member" {
						continue
					}
					// 放入待踢出 map
					userId := gconv.String(kk["user_id"])
					if _, ok := listMap[userId]; ok {
						kickMap[userId] = nil
					}
				}
			}
			// 回执
			service.Bot().SendPlainMsg(ctx, "已踢出 list("+listName+") 中的 group("+gconv.String(groupId)+
				") member\n共 "+gconv.String(listMapLen)+" 条，有 "+gconv.String(len(kickMap))+" 条未踢出")
		}
		// 异步获取群成员列表
		service.Bot().GetGroupMemberList(ctx, groupId, callback2, true)
	}
	// 异步获取群成员列表
	service.Bot().GetGroupMemberList(ctx, groupId, callback, true)
}

func (s *sGroup) KeepFromListReturnRes(ctx context.Context, groupId int64, listName string) {
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
	// 获取 list
	listMap := service.List().GetListData(ctx, listName)
	listMapLen := len(listMap)
	// 异步 callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
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
		kickMap := make(map[string]any)
		// 解析数组
		for _, v := range membersArr {
			// map 断言
			if kk, ok := v.(map[string]any); ok {
				if kk["role"] != "member" {
					continue
				}
				// 放入待踢出 map
				userId := gconv.String(kk["user_id"])
				if _, ok := listMap[userId]; !ok {
					kickMap[userId] = nil
				}
			}
		}
		// 不需要执行踢出
		if len(kickMap) == 0 {
			service.Bot().SendPlainMsg(ctx, "list("+listName+") 中的 group("+gconv.String(groupId)+
				") member 无需踢出")
			return
		}
		// 踢人过程
		service.Bot().SendPlainMsg(ctx, "正在踢出不在 list("+listName+") 中的 group("+gconv.String(groupId)+
			") member\n共 "+gconv.String(listMapLen)+" 条，有 "+gconv.String(len(kickMap))+" 条需要踢出")
		for k := range kickMap {
			// 踢人
			service.Bot().Kick(ctx, groupId, gconv.Int64(k))
			// 随机延时
			time.Sleep(time.Duration(grand.N(1000, 10000)) * time.Millisecond)
		}
		// 异步检查有没有踢出失败的
		callback2 := func(ctx context.Context, rsyncCtx context.Context) {
			if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
				return
			}
			// 获取群成员列表
			membersJson = service.Bot().GetData(rsyncCtx)
			if membersJson == nil {
				// 空列表
				service.Bot().SendPlainMsg(ctx, "获取到空的群成员列表")
				return
			}
			// 局部变量
			membersArr = membersJson.MustArray()
			kickMap = make(map[string]any)
			// 解析数组
			for _, v := range membersArr {
				// map 断言
				if kk, ok := v.(map[string]any); ok {
					if kk["role"] != "member" {
						continue
					}
					// 放入待踢出 map
					userId := gconv.String(kk["user_id"])
					if _, ok := listMap[userId]; !ok {
						kickMap[userId] = nil
					}
				}
			}
			// 回执
			service.Bot().SendPlainMsg(ctx, "已踢出不在 list("+listName+") 中的 group("+gconv.String(groupId)+
				") member\n共 "+gconv.String(listMapLen)+" 条，有 "+gconv.String(len(kickMap))+" 条未踢出")
		}
		// 异步获取群成员列表
		service.Bot().GetGroupMemberList(ctx, groupId, callback2, true)
	}
	// 异步获取群成员列表
	service.Bot().GetGroupMemberList(ctx, groupId, callback, true)
}

func (s *sGroup) CheckExistReturnRes(ctx context.Context) {
	wg := sync.WaitGroup{}
	pageNum, pageSize := 1, 10
	deleted := 0
	msg := ""
	for {
		var gEntity []*entity.Group
		err := dao.Group.Ctx(ctx).
			Fields(dao.Group.Columns().GroupId).
			Page(pageNum, pageSize).
			Scan(&gEntity)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		for _, v := range gEntity {
			vGroupId := v.GroupId
			// 异步 callback
			callback := func(ctx context.Context, rsyncCtx context.Context) {
				defer wg.Done()
				// 不使用默认 echo 处理流程
				//if service.Bot().DefaultEchoProcess(ctx, rsyncCtx) {
				//	return
				//}
				// 判断群是否已经解散
				if service.Bot().GetEchoStatus(rsyncCtx) != "failed" || service.Bot().GetEchoFailedMsg(rsyncCtx) != "群聊不存在" {
					return
				}
				// 记录信息
				msg += "\ngroup(" + gconv.String(vGroupId) + ")"
				// 删除 group
				_, err = dao.Group.Ctx(ctx).Where(dao.Group.Columns().GroupId, vGroupId).Delete()
				if err != nil {
					g.Log().Error(ctx, err)
					return
				}
				deleted++
			}
			wg.Add(1)
			service.Bot().GetGroupInfo(ctx, vGroupId, callback, true)
		}
		// 结束条件
		if len(gEntity) < pageSize {
			break
		}
		pageNum++
	}
	// 回执
	wg.Wait()
	service.Bot().SendPlainMsg(ctx, "已删除 "+gconv.String(deleted)+" 条不存在的 group"+msg)
}
