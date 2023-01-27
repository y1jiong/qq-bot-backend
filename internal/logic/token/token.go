package token

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"
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

func (s *sToken) IsCorrectToken(ctx context.Context, token string) (yes bool, name string) {
	// 过滤非法 token
	if !legalTokenRe.MatchString(token) {
		return
	}
	// 数据库查询
	var tEntity *entity.Token
	err := dao.Token.Ctx(ctx).
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
	yes = true
	name = tEntity.Name
	return
}

func (s *sToken) AddNewToken(ctx context.Context, name, token string, owner int64) {
	// 过滤非法 token 或 name
	if !legalTokenRe.MatchString(token) || !legalTokenNameRe.MatchString(name) {
		return
	}
	tokenEntity := entity.Token{
		Name:    name,
		Token:   token,
		OwnerId: owner,
	}
	// 数据库插入
	_, err := dao.Token.Ctx(ctx).
		Data(tokenEntity).
		OmitEmpty().
		Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已新增 token "+name)
}

func (s *sToken) RemoveToken(ctx context.Context, name string) {
	// 过滤非法 name
	if !legalTokenNameRe.MatchString(name) {
		return
	}
	// 数据库计数
	n, err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if n < 1 {
		service.Bot().SendPlainMsg(ctx, name+" 不存在")
	}
	// 数据库软删除
	_, err = dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已删除 token "+name)
}

func (s *sToken) QueryToken(ctx context.Context) {
	// 数据库查询
	var tEntities []*entity.Token
	err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().OwnerId, service.Bot().GetUserId(ctx)).
		Scan(&tEntities)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 判断空
	if len(tEntities) == 0 {
		return
	}
	// 回执
	var msg strings.Builder
	for _, v := range tEntities {
		msg.WriteString(dao.Token.Columns().Name)
		msg.WriteString(": ")
		msg.WriteString(v.Name)
		msg.WriteString("\n")
		msg.WriteString(dao.Token.Columns().CreatedAt)
		msg.WriteString(": ")
		msg.WriteString(v.CreatedAt.String())
		msg.WriteString("\n---\n")
	}
	service.Bot().SendPlainMsg(ctx, msg.String())
}
