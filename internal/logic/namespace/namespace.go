package namespace

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"
)

type sNamespace struct{}

func New() *sNamespace {
	return &sNamespace{}
}

func init() {
	service.RegisterNamespace(New())
}

var (
	legalNamespaceNameRe = regexp.MustCompile(`^\S{1,16}$`)
)

// setting json key
const (
	adminMapKey = "admins"
	listMapKey  = "lists"
)

func getNamespace(ctx context.Context, namespace string) (nEntity *entity.Namespace) {
	// 数据库查询
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Scan(&nEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 没找到
	if nEntity == nil {
		service.Bot().SendPlainMsg(ctx, "没找到 namespace("+namespace+")")
	}
	return
}

func isNamespaceOwner(userId int64, nEntity *entity.Namespace) (yes bool) {
	// 判断 owner
	if userId != nEntity.OwnerId {
		return
	}
	yes = true
	return
}

func isNamespaceOwnerOrAdmin(ctx context.Context, userId int64, nEntity *entity.Namespace) (yes bool) {
	// 判断 owner
	if userId == nEntity.OwnerId {
		yes = true
		return
	}
	// 解析 setting json
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin map
	admins := settingJson.Get(adminMapKey).MustMap(make(map[string]any))
	// 判断 admin
	if _, ok := admins[gconv.String(userId)]; ok {
		yes = true
	}
	return
}

func (s *sNamespace) IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 || !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 过程
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	return isNamespaceOwnerOrAdmin(ctx, userId, nEntity)
}

func (s *sNamespace) AddNewNamespace(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 初始化 namespace 对象
	nEntity := entity.Namespace{
		Namespace:   namespace,
		OwnerId:     service.Bot().GetUserId(ctx),
		SettingJson: "{}",
	}
	// 数据库插入
	_, err := dao.Namespace.Ctx(ctx).
		Data(nEntity).
		OmitEmpty().
		Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		// 返回错误
		service.Bot().SendPlainMsg(ctx, "新增 namespace 失败")
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已新增 namespace("+namespace+")")
}

func (s *sNamespace) RemoveNamespace(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
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
	service.Bot().SendPlainMsg(ctx, "已删除 namespace("+namespace+")")
}

func (s *sNamespace) QueryNamespace(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断 owner or admin
	if !isNamespaceOwnerOrAdmin(ctx, service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 回执
	msg := dao.Namespace.Columns().Namespace + ": " + nEntity.Namespace + "\n" +
		dao.Namespace.Columns().OwnerId + ": " + gconv.String(nEntity.OwnerId) + "\n" +
		dao.Namespace.Columns().SettingJson + ": " + nEntity.SettingJson + "\n" +
		dao.Namespace.Columns().UpdatedAt + ": " + nEntity.UpdatedAt.String()
	service.Bot().SendPlainMsg(ctx, msg)
}

func (s *sNamespace) QueryOwnNamespace(ctx context.Context, userId int64) {
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
	var msg strings.Builder
	for _, v := range nEntities {
		msg.WriteString(dao.Namespace.Columns().Namespace)
		msg.WriteString(": ")
		msg.WriteString(v.Namespace)
		msg.WriteString("\n")
		msg.WriteString(dao.Namespace.Columns().CreatedAt)
		msg.WriteString(": ")
		msg.WriteString(v.CreatedAt.String())
		msg.WriteString("\n---\n")
	}
	service.Bot().SendPlainMsg(ctx, msg.String())
}

func (s *sNamespace) AddNamespaceAdmin(ctx context.Context, namespace string, userId int64) {
	// 参数合法性校验
	if userId < 1 || !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace 对象
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
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
	// 数据库更新 添加 admin 权限
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, string(settingBytes)).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已添加 user("+gconv.String(userId)+") 的 namespace("+namespace+") admin 权限")
}

func (s *sNamespace) RemoveNamespaceAdmin(ctx context.Context, namespace string, userId int64) {
	// 参数合法性校验
	if userId < 1 || !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 数据库查询
	var nEntity *entity.Namespace
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Scan(&nEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 没找到
	if nEntity == nil {
		service.Bot().SendPlainMsg(ctx, "没找到 namespace "+namespace)
		return
	}
	// 判断是否是 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取 admin map
	admins := settingJson.Get(adminMapKey).MustMap(make(map[string]any))
	// 删除 userId 的 admin 权限
	if _, ok := admins[gconv.String(userId)]; ok {
		delete(admins, gconv.String(userId))
	} else {
		service.Bot().SendPlainMsg(ctx, gconv.String(userId)+"不存在")
		return
	}
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
	service.Bot().SendPlainMsg(ctx, "已取消 user("+gconv.String(userId)+") 的 namespace("+namespace+") admin 权限")
}

func (s *sNamespace) ResetNamespaceAdmin(ctx context.Context, namespace string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 数据库查询
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 判断是否是 owner
	if !isNamespaceOwner(service.Bot().GetUserId(ctx), nEntity) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
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
	service.Bot().SendPlainMsg(ctx, "已重置 namespace("+namespace+") 的 admin")
}

func (s *sNamespace) AddNamespaceList(ctx context.Context, namespace, listName string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	listMap := settingJson.Get(listMapKey).MustMap(make(map[string]any))
	listMap[listName] = nil
	// 保存数据
	settingJson.Set(listMapKey, listMap)
	settingBytes, err := settingJson.Encode()
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
}

func (s *sNamespace) RemoveNamespaceList(ctx context.Context, namespace, listName string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	lists := settingJson.Get(listMapKey).MustMap(make(map[string]any))
	if _, ok := lists[listName]; ok {
		delete(lists, listName)
	} else {
		return
	}
	// 保存数据
	settingJson.Set(listMapKey, lists)
	settingBytes, err := settingJson.Encode()
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
}

func (s *sNamespace) GetNamespaceList(ctx context.Context, namespace string) (lists map[string]any) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	nEntity := getNamespace(ctx, namespace)
	if nEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(nEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	lists = settingJson.Get(listMapKey).MustMap(make(map[string]any))
	return
}
