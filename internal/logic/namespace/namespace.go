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
	if _, ok := lists[listName]; !ok {
		return
	}
	delete(lists, listName)
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
