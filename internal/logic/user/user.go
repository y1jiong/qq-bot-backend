package user

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
)

type sUser struct{}

func New() *sUser {
	return &sUser{}
}

func init() {
	service.RegisterUser(New())
}

var (
	trustKey     = "trust"
	namespaceKey = "namespace"
)

func getUser(ctx context.Context, userId int64) (uEntity *entity.User) {
	// 数据库查询
	err := dao.User.Ctx(ctx).Where(dao.User.Columns().UserId, userId).Scan(&uEntity)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return
}

func (s *sUser) IsSystemTrustUser(ctx context.Context, userId int64) (yes bool) {
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
	if r, ok := settingJson.CheckGet(trustKey); ok {
		yes = r.MustBool()
	}
	return
}

func (s *sUser) SystemTrustUser(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 如果没有获取到 user 则默认创建
		uEntity = &entity.User{
			UserId:      userId,
			SettingJson: "{}",
		}
		// 数据库插入
		_, err := dao.User.Ctx(ctx).
			Data(uEntity).
			OmitEmpty().
			Insert()
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
	service.Bot().SendMsg(ctx, "系统已信任 user("+gconv.String(userId)+")")
}

func (s *sUser) SystemDistrustUser(ctx context.Context, userId int64) {
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
	settingJson.Set(trustKey, false)
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
	service.Bot().SendMsg(ctx, "系统已拒绝信任 user("+gconv.String(userId)+")")
}

func (s *sUser) CouldOperateNamespace(ctx context.Context, userId int64) (yes bool) {
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
	if r, ok := settingJson.CheckGet(namespaceKey); ok {
		yes = r.MustBool()
	}
	return
}

func (s *sUser) GrantOperateNamespace(ctx context.Context, userId int64) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	uEntity := getUser(ctx, userId)
	if uEntity == nil {
		// 如果没有获取到 user 则默认创建
		uEntity = &entity.User{
			UserId:      userId,
			SettingJson: "{}",
		}
		// 数据库插入
		_, err := dao.User.Ctx(ctx).
			Data(uEntity).
			OmitEmpty().
			Insert()
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
	service.Bot().SendMsg(ctx, "系统已授予 user("+gconv.String(userId)+") 操作 namespace 的权限")
}

func (s *sUser) RevokeOperateNamespace(ctx context.Context, userId int64) {
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
	settingJson.Set(namespaceKey, false)
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
	service.Bot().SendMsg(ctx, "系统已撤销 user("+gconv.String(userId)+") 操作 namespace 的权限")
}
