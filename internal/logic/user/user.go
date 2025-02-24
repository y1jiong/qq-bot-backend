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
	recallKey    = "recall"
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
		OmitEmptyData().
		Insert()
	return
}

func (s *sUser) IsSystemTrustedUser(ctx context.Context, userId int64) bool {
	// 参数合法性校验
	if userId == 0 {
		return false
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return false
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	b, _ := settingJson.Get(trustKey).Bool()
	return b
}

func (s *sUser) CanOpToken(ctx context.Context, userId int64) bool {
	// 参数合法性校验
	if userId == 0 {
		return false
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return false
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	b, _ := settingJson.Get(tokenKey).Bool()
	t, _ := settingJson.Get(trustKey).Bool()
	return b || t
}

func (s *sUser) CanOpNamespace(ctx context.Context, userId int64) bool {
	// 参数合法性校验
	if userId == 0 {
		return false
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return false
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	b, _ := settingJson.Get(namespaceKey).Bool()
	t, _ := settingJson.Get(trustKey).Bool()
	return b || t
}

func (s *sUser) CanGetRawMessage(ctx context.Context, userId int64) bool {
	// 参数合法性校验
	if userId == 0 {
		return false
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return false
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	b, _ := settingJson.Get(rawKey).Bool()
	t, _ := settingJson.Get(trustKey).Bool()
	return b || t
}

func (s *sUser) CanRecallMessage(ctx context.Context, userId int64) bool {
	// 参数合法性校验
	if userId == 0 {
		return false
	}
	// 获取 user
	userE := getUser(ctx, userId)
	if userE == nil {
		return false
	}
	// 数据处理
	settingJson, err := sonic.GetFromString(userE.SettingJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	b, _ := settingJson.Get(recallKey).Bool()
	t, _ := settingJson.Get(trustKey).Bool()
	return b || t
}
