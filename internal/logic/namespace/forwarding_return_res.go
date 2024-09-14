package namespace

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
	"strconv"
)

func (s *sNamespace) AddForwardingToReturnRes(ctx context.Context, alias, url, key string) (retMsg string) {
	// 参数合法性校验
	if alias == "" || url == "" {
		return
	}
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
	// forwardingMapKey
	if !settingJson.Get(forwardingMapKey).Valid() {
		_, _ = settingJson.Set(forwardingMapKey, ast.NewNull())
	}
	// toMapKey
	if !settingJson.Get(forwardingMapKey).Get(toMapKey).Valid() {
		_, _ = settingJson.Get(forwardingMapKey).Set(toMapKey, ast.NewNull())
	}
	// alias
	if settingJson.Get(forwardingMapKey).Get(toMapKey).Get(alias).Valid() {
		retMsg = "早已设置 forwarding " + alias
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(toMapKey).Set(alias, ast.NewNull())
	// urlKey
	_, _ = settingJson.Get(forwardingMapKey).Get(toMapKey).Get(alias).Set(
		urlKey, ast.NewString(url))
	// keyKey
	if key != "" {
		_, _ = settingJson.Get(forwardingMapKey).Get(toMapKey).Get(alias).Set(
			keyKey, ast.NewString(key))
	}
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已设置 forwarding " + alias
	return
}

func (s *sNamespace) RemoveForwardingToReturnRes(ctx context.Context, alias string) (retMsg string) {
	// 参数合法性校验
	if alias == "" {
		return
	}
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
	// alias
	if !settingJson.Get(forwardingMapKey).Get(toMapKey).Get(alias).Valid() {
		retMsg = "并未设置 forwarding " + alias
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(toMapKey).Unset(alias)
	// toMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(toMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Unset(toMapKey)
	}
	// forwardingMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Len(); l == 0 {
		_, _ = settingJson.Unset(forwardingMapKey)
	}
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已删除 forwarding " + alias
	return
}

func (s *sNamespace) AddForwardingMatchUserIdReturnRes(ctx context.Context, userId string) (retMsg string) {
	// 参数合法性校验
	if userId == "" {
		return
	}
	if _, err := strconv.Atoi(userId); err != nil {
		if userId != all {
			return
		}
	}
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
	// forwardingMapKey
	if !settingJson.Get(forwardingMapKey).Valid() {
		_, _ = settingJson.Set(forwardingMapKey, ast.NewNull())
	}
	// matchMapKey
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Valid() {
		_, _ = settingJson.Get(forwardingMapKey).Set(matchMapKey, ast.NewNull())
	}
	// userMapKey
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Valid() {
		_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Set(userMapKey, ast.NewNull())
	}
	// userId
	if settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Get(userId).Valid() {
		retMsg = "早已设置 forwarding match user(" + userId + ")"
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Set(userId, ast.NewNull())
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已设置 forwarding match user(" + userId + ")"
	return
}

func (s *sNamespace) AddForwardingMatchGroupIdReturnRes(ctx context.Context, groupId string) (retMsg string) {
	// 参数合法性校验
	if groupId == "" {
		return
	}
	if _, err := strconv.Atoi(groupId); err != nil {
		if groupId != all {
			return
		}
	}
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
	// forwardingMapKey
	if !settingJson.Get(forwardingMapKey).Valid() {
		_, _ = settingJson.Set(forwardingMapKey, ast.NewNull())
	}
	// matchMapKey
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Valid() {
		_, _ = settingJson.Get(forwardingMapKey).Set(matchMapKey, ast.NewNull())
	}
	// groupMapKey
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Valid() {
		_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Set(groupMapKey, ast.NewNull())
	}
	// groupId
	if settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Get(groupId).Valid() {
		retMsg = "早已设置 forwarding match group(" + groupId + ")"
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Set(groupId, ast.NewNull())
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已设置 forwarding match group(" + groupId + ")"
	return
}

func (s *sNamespace) RemoveForwardingMatchUserIdReturnRes(ctx context.Context, userId string) (retMsg string) {
	// 参数合法性校验
	if userId == "" {
		return
	}
	if _, err := strconv.Atoi(userId); err != nil {
		if userId != all {
			return
		}
	}
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
	// userId
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Get(userId).Valid() {
		retMsg = "并未设置 forwarding match user(" + userId + ")"
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Unset(userId)
	// userMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Unset(userMapKey)
	}
	// matchMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(matchMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Unset(matchMapKey)
	}
	// forwardingMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Len(); l == 0 {
		_, _ = settingJson.Unset(forwardingMapKey)
	}
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已删除 forwarding match user(" + userId + ")"
	return
}

func (s *sNamespace) RemoveForwardingMatchGroupIdReturnRes(ctx context.Context, groupId string) (retMsg string) {
	// 参数合法性校验
	if groupId == "" {
		return
	}
	if _, err := strconv.Atoi(groupId); err != nil {
		if groupId != all {
			return
		}
	}
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
	// groupId
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Get(groupId).Valid() {
		retMsg = "并未设置 forwarding match group(" + groupId + ")"
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Unset(groupId)
	// groupMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Unset(groupMapKey)
	}
	// matchMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(matchMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Unset(matchMapKey)
	}
	// forwardingMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Len(); l == 0 {
		_, _ = settingJson.Unset(forwardingMapKey)
	}
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已删除 forwarding match group(" + groupId + ")"
	return
}

func (s *sNamespace) ResetForwardingMatchUserIdReturnRes(ctx context.Context) (retMsg string) {
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
	// userMapKey
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(userMapKey).Valid() {
		retMsg = "并未设置 forwarding match user"
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Unset(userMapKey)
	// matchMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(matchMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Unset(matchMapKey)
	}
	// forwardingMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Len(); l == 0 {
		_, _ = settingJson.Unset(forwardingMapKey)
	}
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已重置 forwarding match user"
	return
}

func (s *sNamespace) ResetForwardingMatchGroupIdReturnRes(ctx context.Context) (retMsg string) {
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
	// groupMapKey
	if !settingJson.Get(forwardingMapKey).Get(matchMapKey).Get(groupMapKey).Valid() {
		retMsg = "并未设置 forwarding match group"
		return
	}
	_, _ = settingJson.Get(forwardingMapKey).Get(matchMapKey).Unset(groupMapKey)
	// matchMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Get(matchMapKey).Len(); l == 0 {
		_, _ = settingJson.Get(forwardingMapKey).Unset(matchMapKey)
	}
	// forwardingMapKey
	if l, _ := settingJson.Get(forwardingMapKey).Len(); l == 0 {
		_, _ = settingJson.Unset(forwardingMapKey)
	}
	settingStr, _ := settingJson.Raw()
	// 更新
	_, err = dao.Namespace.Ctx(ctx).
		Data(dao.Namespace.Columns().SettingJson, settingStr).
		Where(dao.Namespace.Columns().Namespace, globalNamespace).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	retMsg = "已重置 forwarding match group"
	return
}
