package token

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
)

type sToken struct{}

func New() *sToken {
	return &sToken{}
}

func init() {
	service.RegisterToken(New())
}

var (
	legalTokenRe     = regexp.MustCompile(`^\w{16,48}$`)
	legalTokenNameRe = regexp.MustCompile(`^\S{1,16}$`)
)

func (s *sToken) IsCorrectToken(ctx context.Context,
	token string) (correct bool, name string, ownerId, botId int64) {
	// 过滤非法 token
	if !legalTokenRe.MatchString(token) {
		return
	}
	// 数据库查询
	var tokenE *entity.Token
	err := dao.Token.Ctx(ctx).
		Fields(
			dao.Token.Columns().Name,
			dao.Token.Columns().OwnerId,
			dao.Token.Columns().BindingBotId,
		).
		Where(dao.Token.Columns().Token, token).
		Scan(&tokenE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if tokenE == nil {
		return
	}
	// 数据处理
	correct = true
	name = tokenE.Name
	ownerId = tokenE.OwnerId
	botId = tokenE.BindingBotId
	return
}

func (s *sToken) UpdateLoginTime(ctx context.Context, token string) {
	// 数据库更新
	_, err := dao.Token.Ctx(ctx).
		Data(dao.Token.Columns().LastLoginAt, gtime.Now()).
		Where(dao.Token.Columns().Token, token).
		Unscoped().
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
	}
}
