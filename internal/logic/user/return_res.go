package user

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
)

func (s *sUser) QueryUserReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 回执
		retMsg = "查无此人"
		return
	}
	// 回执
	retMsg = dao.User.Columns().UserId + ": " + gconv.String(userE.UserId) + "\n" +
		dao.User.Columns().SettingJson + ": " + userE.SettingJson + "\n" +
		dao.User.Columns().UpdatedAt + ": " + userE.UpdatedAt.String()
	return
}

func (s *sUser) SystemTrustUserReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		userE, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(trustKey); ok {
		// 重复信任
		retMsg = "重复信任"
		return
	}
	settingJson.Set(trustKey, true)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已信任 user(" + gconv.String(userId) + ")"
	return
}

func (s *sUser) SystemDistrustUserReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(trustKey); !ok {
		// 并未信任
		retMsg = "并未信任"
		return
	}
	settingJson.Del(trustKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已不再信任 user(" + gconv.String(userId) + ")"
	return
}

func (s *sUser) GrantOpTokenReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		userE, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(tokenKey); ok {
		retMsg = "重复授予操作 token 的权限"
		return
	}
	settingJson.Set(tokenKey, true)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已授予 user(" + gconv.String(userId) + ") 操作 token 的权限"
	return
}

func (s *sUser) RevokeOpTokenReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(tokenKey); !ok {
		retMsg = "并未授予操作 token 的权限"
		return
	}
	settingJson.Del(tokenKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已撤销 user(" + gconv.String(userId) + ") 操作 token 的权限"
	return
}

func (s *sUser) GrantOpNamespaceReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		userE, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(namespaceKey); ok {
		retMsg = "重复授予操作 namespace 的权限"
		return
	}
	settingJson.Set(namespaceKey, true)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已授予 user(" + gconv.String(userId) + ") 操作 namespace 的权限"
	return
}

func (s *sUser) RevokeOpNamespaceReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(namespaceKey); !ok {
		retMsg = "并未授予操作 namespace 的权限"
		return
	}
	settingJson.Del(namespaceKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已撤销 user(" + gconv.String(userId) + ") 操作 namespace 的权限"
	return
}

func (s *sUser) GrantGetRawMsgReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		userE, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(rawKey); ok {
		retMsg = "重复授予获取 raw 的权限"
		return
	}
	settingJson.Set(rawKey, true)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已授予 user(" + gconv.String(userId) + ") 获取 raw 的权限"
	return
}

func (s *sUser) RevokeGetRawMsgReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(rawKey); !ok {
		retMsg = "并未授予获取 raw 的权限"
		return
	}
	settingJson.Del(rawKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已撤销 user(" + gconv.String(userId) + ") 获取 raw 的权限"
	return
}

func (s *sUser) GrantRecallReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		userE, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(recallKey); ok {
		retMsg = "重复授予 recall 的权限"
		return
	}
	settingJson.Set(recallKey, true)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已授予 user(" + gconv.String(userId) + ") recall 的权限"
	return
}

func (s *sUser) RevokeRecallReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(recallKey); !ok {
		retMsg = "并未授予 recall 的权限"
		return
	}
	settingJson.Del(recallKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已撤销 user(" + gconv.String(userId) + ") recall 的权限"
	return
}

func (s *sUser) GrantOpCrontabReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		userE, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(crontabKey); ok {
		retMsg = "重复授予操作 crontab 的权限"
		return
	}
	settingJson.Set(crontabKey, true)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已授予 user(" + gconv.String(userId) + ") 操作 crontab 的权限"
	return
}

func (s *sUser) RevokeOpCrontabReturnRes(ctx context.Context, userId int64) (retMsg string) {
	// 参数合法性校验
	if userId == 0 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(userE.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(crontabKey); !ok {
		retMsg = "并未授予操作 crontab 的权限"
		return
	}
	settingJson.Del(crontabKey)
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Data(dao.User.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已撤销 user(" + gconv.String(userId) + ") 操作 crontab 的权限"
	return
}
