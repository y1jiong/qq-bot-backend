package user

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
)

func (s *sUser) QueryUserReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 回执
		service.Bot().SendPlainMsg(ctx, "查无此人")
		return
	}
	msg := dao.User.Columns().UserId + ": " + gconv.String(uEntity.UserId) + "\n" +
		dao.User.Columns().SettingJson + ": " + uEntity.SettingJson + "\n" +
		dao.User.Columns().UpdatedAt + ": " + uEntity.UpdatedAt.String()
	// 回执
	service.Bot().SendPlainMsg(ctx, msg)
}

func (s *sUser) SystemTrustUserReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		uEntity, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(trustKey); ok {
		// 重复信任
		service.Bot().SendPlainMsg(ctx, "重复信任")
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
	service.Bot().SendPlainMsg(ctx, "系统已信任 user("+gconv.String(userId)+")")
}

func (s *sUser) SystemDistrustUserReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(trustKey); !ok {
		// 并未信任
		service.Bot().SendPlainMsg(ctx, "并未信任")
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
	service.Bot().SendPlainMsg(ctx, "系统已拒绝信任 user("+gconv.String(userId)+")")
}

func (s *sUser) GrantOpTokenReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		uEntity, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(tokenKey); ok {
		// 重复授予
		service.Bot().SendPlainMsg(ctx, "重复授予操作 token 的权限")
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
	service.Bot().SendPlainMsg(ctx, "系统已授予 user("+gconv.String(userId)+") 操作 token 的权限")
}

func (s *sUser) RevokeOpTokenReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(tokenKey); !ok {
		// 并未授予
		service.Bot().SendPlainMsg(ctx, "并未授予操作 token 的权限")
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
	service.Bot().SendPlainMsg(ctx, "系统已撤销 user("+gconv.String(userId)+") 操作 token 的权限")
}

func (s *sUser) GrantOpNamespaceReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		uEntity, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(namespaceKey); ok {
		// 重复授予
		service.Bot().SendPlainMsg(ctx, "重复授予操作 namespace 的权限")
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
	service.Bot().SendPlainMsg(ctx, "系统已授予 user("+gconv.String(userId)+") 操作 namespace 的权限")
}

func (s *sUser) RevokeOpNamespaceReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(namespaceKey); !ok {
		// 并未授予
		service.Bot().SendPlainMsg(ctx, "并未授予操作 namespace 的权限")
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
	service.Bot().SendPlainMsg(ctx, "系统已撤销 user("+gconv.String(userId)+") 操作 namespace 的权限")
}

func (s *sUser) GrantGetRawMsgReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 如果没有获取到 user 则默认创建
		var err error
		uEntity, err = createUser(ctx, userId)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(rawKey); ok {
		// 重复授予
		service.Bot().SendPlainMsg(ctx, "重复授予获取 raw 的权限")
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
	service.Bot().SendPlainMsg(ctx, "系统已授予 user("+gconv.String(userId)+") 获取 raw 的权限")
}

func (s *sUser) RevokeGetRawMsgReturnRes(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(uEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(rawKey); !ok {
		// 并未授予
		service.Bot().SendPlainMsg(ctx, "并未授予获取 raw 的权限")
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
	service.Bot().SendPlainMsg(ctx, "系统已撤销 user("+gconv.String(userId)+") 获取 raw 的权限")
}
