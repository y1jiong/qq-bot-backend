package token

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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
	var tokenE *entity.Token
	err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Scan(&tokenE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	var exist bool
	// 判断是否存在
	if tokenE != nil {
		// 判断所有人是否一致
		if tokenE.OwnerId != owner {
			retMsg = "token(" + name + ") 已被占用"
			return
		}
		exist = true
	}
	tokenE = &entity.Token{
		Name:    name,
		Token:   token,
		OwnerId: owner,
	}
	if !exist {
		// 数据库插入
		_, err = dao.Token.Ctx(ctx).
			Data(tokenE).
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
			Data(tokenE).
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
		Where(dao.Token.Columns().Name, name).
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
		msg.WriteString(dao.Token.Columns().UpdatedAt)
		msg.WriteString(": ")
		msg.WriteString(v.UpdatedAt.String())
		msg.WriteString("\n")
		msg.WriteString(dao.Token.Columns().LastLoginAt)
		msg.WriteString(": ")
		msg.WriteString(v.LastLoginAt.String())
		msg.WriteString("\n")
		msg.WriteString(dao.Token.Columns().BindingBotId)
		msg.WriteString(": ")
		msg.WriteString(gconv.String(v.BindingBotId))
		if i != tEntitiesLen-1 {
			msg.WriteString("\n---\n")
		}
	}
	retMsg = msg.String()
	return
}

func (s *sToken) ChangeTokenOwnerReturnRes(ctx context.Context, name, ownerId string) (retMsg string) {
	// 过滤非法 name
	if !legalTokenNameRe.MatchString(name) {
		return
	}
	// 数据库查存在
	var tokenE *entity.Token
	err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Scan(&tokenE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if tokenE == nil {
		retMsg = "未找到 token(" + name + ")"
		return
	}
	// 数据库更新
	_, err = dao.Token.Ctx(ctx).
		Data(dao.Token.Columns().OwnerId, gconv.Int64(ownerId)).
		Where(dao.Token.Columns().Name, name).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 token(" + name + ") 的所有者修改为 " + ownerId
	return
}

func (s *sToken) BindTokenBotId(ctx context.Context, name, botId string) (retMsg string) {
	// 过滤非法 name
	if !legalTokenNameRe.MatchString(name) {
		return
	}
	// 数据库查存在
	var tokenE *entity.Token
	err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Scan(&tokenE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if tokenE == nil {
		retMsg = "未找到 token(" + name + ")"
		return
	}
	// 数据库更新
	_, err = dao.Token.Ctx(ctx).
		Data(dao.Token.Columns().BindingBotId, gconv.Int64(botId)).
		Where(dao.Token.Columns().Name, name).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 token(" + name + ") 绑定的 bot_id 修改为 " + botId
	return
}

func (s *sToken) UnbindTokenBotId(ctx context.Context, name string) (retMsg string) {
	// 过滤非法 name
	if !legalTokenNameRe.MatchString(name) {
		return
	}
	// 数据库查存在
	var tokenE *entity.Token
	err := dao.Token.Ctx(ctx).
		Where(dao.Token.Columns().Name, name).
		Scan(&tokenE)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if tokenE == nil {
		retMsg = "未找到 token(" + name + ")"
		return
	}
	// 数据库更新
	_, err = dao.Token.Ctx(ctx).
		Data(dao.Token.Columns().BindingBotId, nil).
		Where(dao.Token.Columns().Name, name).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已将 token(" + name + ") 解绑 bot_id"
	return
}
