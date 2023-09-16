package token

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"strings"
)

func (s *sToken) AddNewTokenReturnRes(ctx context.Context, name, token string) (retMsg string) {
	owner := service.Bot().GetUserId(ctx)
	// 过滤非法 token 或 name
	if !legalTokenRe.MatchString(token) || !legalTokenNameRe.MatchString(name) {
		return
	}
	// 数据库查存在
	var tokenEntity *entity.Token
	err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Scan(&tokenEntity)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	var exist bool
	// 判断是否存在
	if tokenEntity != nil {
		// 判断所有人是否一致
		if tokenEntity.OwnerId != owner {
			retMsg = "token(" + name + ") 已被占用"
			return
		}
		exist = true
	}
	tokenEntity = &entity.Token{
		Name:    name,
		Token:   token,
		OwnerId: owner,
	}
	if !exist {
		// 数据库插入
		_, err = dao.Token.Ctx(ctx).
			Data(tokenEntity).
			OmitEmpty().
			Insert()
		if err != nil {
			g.Log().Error(ctx, err)
			// 返回错误
			retMsg = "新增 token 失败"
			return
		}
		// 回执
		retMsg = "已新增 token(" + name + ")"
	} else {
		// 数据库更新
		_, err = dao.Token.Ctx(ctx).
			Data(tokenEntity).
			OmitEmpty().
			Where(dao.Token.Columns().Name, name).
			Update()
		if err != nil {
			g.Log().Error(ctx, err)
			// 返回错误
			retMsg = "更新 token 失败"
			return
		}
		// 回执
		retMsg = "已更新 token(" + name + ")"
	}
	return
}

func (s *sToken) RemoveTokenReturnRes(ctx context.Context, name string) (retMsg string) {
	// 过滤非法 name
	if !legalTokenNameRe.MatchString(name) {
		return
	}
	// 数据库查存在
	one, err := dao.Token.Ctx(ctx).
		Where(g.Map{
			dao.Token.Columns().Name:    name,
			dao.Token.Columns().OwnerId: service.Bot().GetUserId(ctx),
		}).
		One()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if one.IsEmpty() {
		retMsg = "未找到 token(" + name + ")"
		return
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
	retMsg = "已删除 token(" + name + ")"
	return
}

func (s *sToken) QueryTokenReturnRes(ctx context.Context) (retMsg string) {
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
	tEntitiesLen := len(tEntities)
	for i, v := range tEntities {
		msg.WriteString(dao.Token.Columns().Name)
		msg.WriteString(": ")
		msg.WriteString(v.Name)
		msg.WriteString("\n")
		msg.WriteString(dao.Token.Columns().CreatedAt)
		msg.WriteString(": ")
		msg.WriteString(v.CreatedAt.String())
		msg.WriteString("\n")
		msg.WriteString(dao.Token.Columns().LastLoginAt)
		msg.WriteString(": ")
		msg.WriteString(v.LastLoginAt.String())
		if i != tEntitiesLen-1 {
			msg.WriteString("\n---\n")
		}
	}
	retMsg = msg.String()
	return
}
