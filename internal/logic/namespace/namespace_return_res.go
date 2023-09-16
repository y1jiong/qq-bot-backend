package namespace

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"strings"
)

func (s *sNamespace) AddNewNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 数据库查存在
	one := getNamespace(ctx, namespace)
	if one != nil {
		retMsg = "namespace(" + namespace + ") 已被占用"
		return
	}
	// 初始化 namespace 对象
	namespaceE := entity.Namespace{
		Namespace:   namespace,
		OwnerId:     service.Bot().GetUserId(ctx),
		SettingJson: "{}",
	}
	// 数据库插入
	_, err := dao.Namespace.Ctx(ctx).
		Data(namespaceE).
		OmitEmpty().
		Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		// 返回错误
		retMsg = "新增 namespace 失败"
		return
	}
	// 回执
	retMsg = "已新增 namespace(" + namespace + ")"
	return
}

func (s *sNamespace) RemoveNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		retMsg = "未找到 namespace(" + namespace + ")"
		return
	}
	// 判断 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) {
		retMsg = "未找到 namespace(" + namespace + ")"
		return
	}
	// 数据库软删除
	_, err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已删除 namespace(" + namespace + ")"
	return
}

func (s *sNamespace) QueryNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 判断 owner or admin
	if !isNamespaceOwnerOrAdmin(ctx, service.Bot().GetUserId(ctx), namespaceE) {
		return
	}
	// 回执
	retMsg = dao.Namespace.Columns().Namespace + ": " + namespaceE.Namespace + "\n" +
		dao.Namespace.Columns().OwnerId + ": " + gconv.String(namespaceE.OwnerId) + "\n" +
		dao.Namespace.Columns().SettingJson + ": " + namespaceE.SettingJson + "\n" +
		dao.Namespace.Columns().UpdatedAt + ": " + namespaceE.UpdatedAt.String()
	return
}

func (s *sNamespace) QueryOwnNamespaceReturnRes(ctx context.Context) (retMsg string) {
	userId := service.Bot().GetUserId(ctx)
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 创建数组指针
	var nEntities []*entity.Namespace
	// 数据库查询
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().OwnerId, userId).
		Scan(&nEntities)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 判断空
	if len(nEntities) == 0 {
		return
	}
	// 回执
	nEntitiesLen := len(nEntities)
	var msg strings.Builder
	for i, v := range nEntities {
		msg.WriteString(dao.Namespace.Columns().Namespace)
		msg.WriteString(": ")
		msg.WriteString(v.Namespace)
		msg.WriteString("\n")
		msg.WriteString(dao.Namespace.Columns().CreatedAt)
		msg.WriteString(": ")
		msg.WriteString(v.CreatedAt.String())
		if i != nEntitiesLen-1 {
			msg.WriteString("\n---\n")
		}
	}
	retMsg = msg.String()
	return
}

func (s *sNamespace) AddNamespaceAdminReturnRes(ctx context.Context,
	namespace string, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId < 1 || !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 判断 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(namespaceE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin map
	admins := settingJson.Get(adminMapKey).MustMap(make(map[string]any))
	// 添加 userId 的 admin 权限
	admins[gconv.String(userId)] = nil
	// 保存数据
	settingJson.Set(adminMapKey, admins)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已添加 user(" + gconv.String(userId) + ") 的 namespace(" + namespace + ") admin 权限"
	return
}

func (s *sNamespace) RemoveNamespaceAdminReturnRes(ctx context.Context,
	namespace string, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId < 1 || !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 判断是否是 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(namespaceE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin map
	admins := settingJson.Get(adminMapKey).MustMap(make(map[string]any))
	if _, ok := admins[gconv.String(userId)]; !ok {
		retMsg = "在 namespace(" + namespaceE.Namespace + ") 的 " + adminMapKey + " 中未找到 user(" + gconv.String(userId) + ")"
		return
	}
	// 删除 userId 的 admin 权限
	delete(admins, gconv.String(userId))
	// 保存数据
	settingJson.Set(adminMapKey, admins)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已取消 user(" + gconv.String(userId) + ") 的 namespace(" + namespace + ") admin 权限"
	return
}

func (s *sNamespace) ResetNamespaceAdminReturnRes(ctx context.Context, namespace string) (retMsg string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 判断是否是 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(namespaceE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 重置 admin
	settingJson.Del(adminMapKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已重置 namespace(" + namespace + ") 的 admin"
	return
}
