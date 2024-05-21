package namespace

import (
	"context"
	"fmt"
	sj "github.com/bitly/go-simplejson"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
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
		OmitEmptyData().
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
	// 权限校验 判断 owner
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
	// 权限校验 判断 owner or admin
	if permission, public := isNamespaceOwnerOrAdmin(ctx, service.Bot().GetUserId(ctx), namespaceE),
		s.IsNamespacePropertyPublic(ctx, namespace) || s.IsSharedNamespace(namespace); !permission && !public {
		return
	} else if permission {
		retMsg = dao.Namespace.Columns().OwnerId + ": " + gconv.String(namespaceE.OwnerId) + "\n"
	}
	// 回执
	retMsg = dao.Namespace.Columns().Namespace + ": " + namespaceE.Namespace + "\n" +
		retMsg +
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
	query := dao.Namespace.Ctx(ctx).
		Fields(
			dao.Namespace.Columns().Namespace,
			dao.Namespace.Columns().CreatedAt,
		)
	if !service.User().CouldOpNamespace(ctx, userId) {
		query = query.Where(dao.Namespace.Columns().OwnerId, userId)
		query = query.WhereOr(fmt.Sprintf("%v->'%v'->>'%v'='%v'",
			dao.Namespace.Columns().SettingJson, propertiesMapKey, propertyPublic, true),
		)
	}
	err := query.Scan(&nEntities)
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
		msg.WriteString(fmt.Sprintf("%s: %s\n%s: %s",
			dao.Namespace.Columns().Namespace, v.Namespace,
			dao.Namespace.Columns().CreatedAt, v.CreatedAt.String()))
		if i != nEntitiesLen-1 {
			msg.WriteString("\n")
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
	// 权限校验 判断 owner 或者 namespace op
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) &&
		!service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(namespaceE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin map
	admins := settingJson.Get(adminsMapKey).MustMap(make(map[string]any))
	// 添加 userId 的 admin 权限
	admins[gconv.String(userId)] = nil
	// 保存数据
	settingJson.Set(adminsMapKey, admins)
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
	// 权限校验 判断 owner 或者 namespace op
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) &&
		!service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(namespaceE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin map
	admins := settingJson.Get(adminsMapKey).MustMap(make(map[string]any))
	if _, ok := admins[gconv.String(userId)]; !ok {
		retMsg = "在 namespace(" + namespaceE.Namespace + ") 的 " + adminsMapKey + " 中未找到 user(" + gconv.String(userId) + ")"
		return
	}
	// 删除 userId 的 admin 权限
	delete(admins, gconv.String(userId))
	// 保存数据
	if len(admins) == 0 {
		settingJson.Del(adminsMapKey)
	} else {
		settingJson.Set(adminsMapKey, admins)
	}
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
	// 权限校验 判断 owner
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
	settingJson.Del(adminsMapKey)
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

func (s *sNamespace) ChangeNamespaceOwnerReturnRes(ctx context.Context,
	namespace, ownerId string) (retMsg string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 权限校验 判断 owner 或者 namespace op
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) &&
		!service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据库更新
	_, err := dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().OwnerId, ownerId).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 namespace(" + namespace + ") 的 owner 修改为 user(" + ownerId + ")"
	return
}

func (s *sNamespace) SetNamespacePropertyPublicReturnRes(ctx context.Context,
	namespace string, value bool) (retMsg string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 权限校验 判断 owner 或者 namespace op
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), namespaceE) &&
		!service.User().CouldOpNamespace(ctx, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if value {
		if !settingJson.Get(propertiesMapKey).Valid() {
			_, _ = settingJson.Set(propertiesMapKey, ast.NewNull())
		}
		if v, _ := settingJson.Get(propertiesMapKey).Get(propertyPublic).Bool(); v {
			retMsg = "早已设置 namespace(" + namespace + ") 为 public"
			return
		}
		_, _ = settingJson.Get(propertiesMapKey).Set(propertyPublic, ast.NewBool(value))
	} else {
		if settingJson.Get(propertiesMapKey).Get(propertyPublic).Exists() {
			_, _ = settingJson.Get(propertiesMapKey).Unset(propertyPublic)
		} else {
			retMsg = "并未设置 namespace(" + namespace + ") 为 public"
			return
		}
		// 删除 properties map 如果为空
		if l, _ := settingJson.Get(propertiesMapKey).Len(); l == 0 {
			_, _ = settingJson.Unset(propertiesMapKey)
		}
	}
	// 保存数据
	settingBytes, err := settingJson.MarshalJSON()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
	}
	// 回执
	if value {
		retMsg = "已设置 namespace(" + namespace + ") 为 public"
	} else {
		retMsg = "已取消 namespace(" + namespace + ") 的 public 属性"
	}
	return
}
