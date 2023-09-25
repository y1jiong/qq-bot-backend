package user

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
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

const (
	trustKey     = "trust"
	tokenKey     = "token"
	namespaceKey = "namespace"
	rawKey       = "raw"
)

func getUser(ctx context.Context, userId int64) (userE *entity.User) {
	// 数据库查询
	err := dao.User.Ctx(ctx).
		Where(dao.User.Columns().UserId, userId).
		Scan(&userE)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return
}

func createUser(ctx context.Context, userId int64) (userE *entity.User, err error) {
	userE = &entity.User{
		UserId:      userId,
		SettingJson: "{}",
	}
	// 数据库插入
	_, err = dao.User.Ctx(ctx).
		Data(userE).
		OmitEmpty().
		Insert()
	return
}

func (s *sUser) IsSystemTrustUser(ctx context.Context, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	yes, _ = settingJson.Get(trustKey).Bool()
	return
}

func (s *sUser) CouldOpToken(ctx context.Context, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	yes, _ = settingJson.Get(tokenKey).Bool()
	return
}

func (s *sUser) CouldOpNamespace(ctx context.Context, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	yes, _ = settingJson.Get(namespaceKey).Bool()
	return
}

func (s *sUser) CouldGetRawMsg(ctx context.Context, userId int64) (yes bool) {
	// 参数合法性校验
	if userId < 1 {
		return
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	yes, _ = settingJson.Get(rawKey).Bool()
	return
}
