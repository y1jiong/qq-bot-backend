package message

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"qq-bot-backend/internal/service"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"qq-bot-backend/api/message/v1"
)

func (c *ControllerV1) Message(ctx context.Context, req *v1.MessageReq) (res *v1.MessageRes, err error) {
	if req.Token == "" {
		// 忽视前置的 Bearer 或 Token 进行鉴权
		authorizations := strings.Fields(g.RequestFromCtx(ctx).Header.Get("Authorization"))
		if len(authorizations) < 2 {
			err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
				http.StatusText(http.StatusForbidden))
			return
		}
		req.Token = authorizations[1]
	}
	// token 验证
	pass, tokenName, ownerId, botId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
			http.StatusText(http.StatusForbidden))
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, service.Namespace().GetGlobalNamespace(), ownerId) {
		if req.GroupId == 0 {
			err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
				"permission denied")
			return
		}
		namespace := service.Group().GetNamespace(ctx, req.GroupId)
		if namespace == "" {
			err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
				"permission denied")
			return
		}
		if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, namespace, ownerId) {
			err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
				"permission denied")
			return
		}
	}
	// 记录访问时间
	service.Token().UpdateLoginTime(ctx, req.Token)
	// 加载 botId 对应的 botCtx
	botCtx := service.Bot().LoadConnectionPool(botId)
	if botCtx == nil {
		err = gerror.NewCode(gcode.New(http.StatusInternalServerError, "", nil),
			"bot not connected")
		return
	}
	// 规范请求参数
	if req.GroupId != 0 && req.UserId != 0 {
		req.UserId = 0
	}
	// for log
	{
		inner := struct {
			UserId  int64  `json:"user_id,omitempty"`
			GroupId int64  `json:"group_id,omitempty"`
			Message string `json:"message"`
		}{
			UserId:  req.UserId,
			GroupId: req.GroupId,
			Message: req.Message,
		}
		var innerStr string
		innerStr, err = sonic.ConfigDefault.MarshalToString(inner)
		if err != nil {
			return
		}
		g.Log().Info(ctx, tokenName+" access successfully with "+innerStr)
	}
	// 限速 一分钟只能发送 7 条消息
	if limit, _ := service.Util().AutoLimit(ctx,
		"send_msg", gconv.String(req.UserId+req.GroupId), 7, time.Minute); limit {
		err = gerror.NewCode(gcode.New(http.StatusTooManyRequests, "", nil),
			http.StatusText(http.StatusTooManyRequests))
		return
	}
	// send message
	_, err = service.Bot().SendMessage(botCtx,
		service.Bot().GuessMsgType(req.GroupId), req.UserId, req.GroupId, req.Message, false)
	if err != nil {
		err = gerror.NewCode(gcode.New(http.StatusInternalServerError, "", nil),
			err.Error())
		return
	}
	return
}
