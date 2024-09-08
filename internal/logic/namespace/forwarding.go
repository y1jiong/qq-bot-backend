package namespace

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
)

const (
	forwardingMapKey = "forwarding"

	// toMapKey forwarding -> to
	toMapKey = "to"
	// authorizationKey forwarding -> to -> authorization
	authorizationKey = "authorization"
	// urlKey forwarding -> to -> url
	urlKey = "url"

	// matchMapKey forwarding -> match
	matchMapKey = "match"
	// userMapKey forwarding -> match -> user
	userMapKey = "user"
	// groupMapKey forwarding -> match -> group
	groupMapKey = "group"

	all = "all"
)

func (s *sNamespace) GetForwardingToAliasList(ctx context.Context) (aliasList map[string]any) {
	// 过程
	namespaceE := getNamespace(ctx, globalNamespace)
	if namespaceE == nil {
		return
	}
	// 解析 setting json
	settingJson, err := sonic.GetFromString(namespaceE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取
	if !settingJson.Get(forwardingMapKey).Get(toMapKey).Valid() {
		return
	}
	aliasList, err = settingJson.Get(forwardingMapKey).Get(toMapKey).Map()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if aliasList == nil {
		aliasList = make(map[string]any)
	}
	return
}

func (s *sNamespace) GetForwardingTo(ctx context.Context, alias string) (url, authorization string) {
	data := &struct {
		Authorization string `orm:"authorization"`
		URL           string `orm:"url"`
	}{}
	// 数据库查询 json
	err := dao.Namespace.Ctx(ctx).
		Fields(fmt.Sprintf(`%v#>>'{%v,%v,%v,%v}' as "authorization"`,
			dao.Namespace.Columns().SettingJson, forwardingMapKey, toMapKey, alias, authorizationKey),
		).
		Fields(fmt.Sprintf(`%v#>>'{%v,%v,%v,%v}' as "url"`,
			dao.Namespace.Columns().SettingJson, forwardingMapKey, toMapKey, alias, urlKey),
		).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Scan(&data)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	url = data.URL
	authorization = data.Authorization
	return
}

func (s *sNamespace) IsForwardingMatchUserId(ctx context.Context, userId string) bool {
	var namespaceE *entity.Namespace
	// 数据库查询 all
	err := dao.Namespace.Ctx(ctx).
		Fields(dao.Namespace.Columns().Namespace).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Where(fmt.Sprintf(`jsonb_path_exists(%v, '$."%v"."%v"."%v"."%v"')`,
			dao.Namespace.Columns().SettingJson, forwardingMapKey, matchMapKey, userMapKey, all),
		).
		Scan(&namespaceE)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if namespaceE != nil {
		return true
	}
	// 数据库查询 json
	err = dao.Namespace.Ctx(ctx).
		Fields(dao.Namespace.Columns().Namespace).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Where(fmt.Sprintf(`jsonb_path_exists(%v, '$."%v"."%v"."%v"."%v"')`,
			dao.Namespace.Columns().SettingJson, forwardingMapKey, matchMapKey, userMapKey, userId),
		).
		Scan(&namespaceE)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if namespaceE == nil {
		return false
	}
	return true
}

func (s *sNamespace) IsForwardingMatchGroupId(ctx context.Context, groupId string) bool {
	var namespaceE *entity.Namespace
	// 数据库查询 all
	err := dao.Namespace.Ctx(ctx).
		Fields(dao.Namespace.Columns().Namespace).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Where(fmt.Sprintf(`jsonb_path_exists(%v, '$."%v"."%v"."%v"."%v"')`,
			dao.Namespace.Columns().SettingJson, forwardingMapKey, matchMapKey, groupMapKey, all),
		).
		Scan(&namespaceE)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if namespaceE != nil {
		return true
	}
	// 数据库查询 json
	err = dao.Namespace.Ctx(ctx).
		Fields(dao.Namespace.Columns().Namespace).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Where(fmt.Sprintf(`jsonb_path_exists(%v, '$."%v"."%v"."%v"."%v"')`,
			dao.Namespace.Columns().SettingJson, forwardingMapKey, matchMapKey, groupMapKey, groupId),
		).
		Scan(&namespaceE)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if namespaceE == nil {
		return false
	}
	return true
}
