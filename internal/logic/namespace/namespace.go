package namespace

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
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

func getNamespace(ctx context.Context, namespace string) (namespaceE *entity.Namespace) {
	// 数据库查询
	err := dao.Namespace.Ctx(ctx).
		Where(dao.Namespace.Columns().Namespace, namespace).
		Scan(&namespaceE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return
}

func isNamespaceOwner(userId int64, namespaceE *entity.Namespace) (yes bool) {
	// 判断 owner
	if userId != namespaceE.OwnerId {
		return
	}
	yes = true
	return
}

func isNamespaceOwnerOrAdmin(ctx context.Context, userId int64, namespaceE *entity.Namespace) (yes bool) {
	// 判断 owner
	if userId == namespaceE.OwnerId {
		yes = true
		return
	}
	// 解析 setting json
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 判断 admin
	yes = settingJson.Get(adminMapKey).Get(gconv.String(userId)).Valid()
	return
}

func (s *sNamespace) IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 || !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 过程
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	return isNamespaceOwnerOrAdmin(ctx, userId, namespaceE)
}

func (s *sNamespace) AddNamespaceList(ctx context.Context, namespace, listName string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	_, _ = settingJson.Get(listMapKey).Set(listName, ast.NewNull())
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
}

func (s *sNamespace) RemoveNamespaceList(ctx context.Context, namespace, listName string) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if !settingJson.Get(listMapKey).Get(listName).Valid() {
		return
	}
	_, _ = settingJson.Get(listMapKey).Unset(listName)
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
}

func (s *sNamespace) GetNamespaceList(ctx context.Context, namespace string) (lists map[string]any) {
	// 参数合法性校验
	if !legalNamespaceNameRe.MatchString(namespace) {
		return
	}
	// 获取 namespace
	namespaceE := getNamespace(ctx, namespace)
	if namespaceE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	lists, _ = settingJson.Get(listMapKey).Map()
	if lists == nil {
		lists = make(map[string]any)
	}
	return
}
