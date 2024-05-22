package namespace

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
)

const (
	listsMapKey = "lists"
)

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
	if !settingJson.Get(listsMapKey).Valid() {
		_, _ = settingJson.Set(listsMapKey, ast.NewNull())
	}
	_, _ = settingJson.Get(listsMapKey).Set(listName, ast.NewNull())
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
	if !settingJson.Get(listsMapKey).Get(listName).Exists() {
		return
	}
	_, _ = settingJson.Get(listsMapKey).Unset(listName)
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

func (s *sNamespace) GetNamespaceLists(ctx context.Context, namespace string) (lists map[string]any) {
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
	lists, _ = settingJson.Get(listsMapKey).Map()
	if lists == nil {
		lists = make(map[string]any)
	}
	return
}

func (s *sNamespace) GetNamespaceListsIncludingShared(ctx context.Context, namespace string) (lists map[string]any) {
	// 先加载 namespace list
	lists = s.GetNamespaceLists(ctx, namespace)
	// 加载公共 list
	namespaceE := getNamespace(ctx, sharedNamespace)
	if namespaceE == nil {
		return
	}
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	sharedLists, _ := settingJson.Get(listsMapKey).Map()
	if sharedLists == nil {
		return
	}
	for k, v := range sharedLists {
		lists[k] = v
	}
	return
}

func (s *sNamespace) GetSharedNamespaceLists(ctx context.Context) (lists map[string]any) {
	return s.GetNamespaceLists(ctx, sharedNamespace)
}
