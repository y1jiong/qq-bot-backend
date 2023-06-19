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

func (s *sToken) IsCorrectToken(ctx context.Context, token string) (correct bool, name string) {
	// 过滤非法 token
	if !legalTokenRe.MatchString(token) {
		return
	}
	// 数据库查询
	var tEntity *entity.Token
	err := dao.Token.Ctx(ctx).
		Fields(dao.Token.Columns().Name).
		Where(dao.Token.Columns().Token, token).
		Scan(&tEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if tEntity == nil {
		return
	}
	// 数据处理
	correct = true
	name = tEntity.Name
	return
}

func (s *sToken) UpdateLoginTime(ctx context.Context, token string) {
	// 数据库更新
	_, err := dao.Token.Ctx(ctx).
		Data(g.Map{
			dao.Token.Columns().LastLoginAt: gtime.Now(),
		}).
		Where(dao.Token.Columns().Token, token).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
	}
}
