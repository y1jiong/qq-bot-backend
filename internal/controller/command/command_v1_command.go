package command

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"qq-bot-backend/internal/service"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"qq-bot-backend/api/command/v1"
)

func (c *ControllerV1) Command(ctx context.Context, req *v1.CommandReq) (res *v1.CommandRes, err error) {
	// 验证请求时间有效性
	{
		msgTime := gtime.New(time.Unix(req.Timestamp, 0))
		if diff := gtime.Now().Sub(msgTime); diff > 5*time.Second {
			err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil),
				"message expired")
			return
		} else if diff < -5*time.Second {
			err = gerror.NewCode(gcode.New(http.StatusTooEarly, "", nil),
				http.StatusText(http.StatusTooEarly))
			return
		}
	}
	// 验证 token
	pass, tokenName, ownerId, botId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
			"permission denied")
		return
	}
	// 防止重放攻击
	if limit, _ := service.Util().AutoLimit(ctx,
		"api.command", req.Signature, 1, 10*time.Second); limit {
		err = gerror.NewCode(gcode.New(http.StatusConflict, "", nil),
			http.StatusText(http.StatusConflict))
		return
	}
	// 验证签名
	{
		// 以 token+command+group_id+timestamp+message_sync+async 为原文，
		// 以 token_name 为 key 的 HmacSha1 值的 base64 值
		s := req.Token + req.Command + gconv.String(req.GroupId) +
			gconv.String(req.Timestamp) + gconv.String(req.MessageSync) +
			gconv.String(req.Async)
		// HmacSha1
		hmacSha1 := hmac.New(sha1.New, []byte(tokenName))
		hmacSha1.Write([]byte(s))
		macBase64 := gbase64.Encode(hmacSha1.Sum(nil))
		if !hmac.Equal(macBase64, []byte(req.Signature)) {
			err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil),
				"signature error")
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
	// 初始化内部请求
	innerReq := struct {
		ApiReq  struct{} `json:"api_req"`
		UserId  int64    `json:"user_id"`
		GroupId int64    `json:"group_id"`
	}{
		ApiReq:  struct{}{},
		UserId:  ownerId,
		GroupId: req.GroupId,
	}
	rawJson, err := sonic.ConfigDefault.MarshalToString(innerReq)
	if err != nil {
		return
	}
	reqJson, _ := sonic.GetFromString(rawJson)
	botCtx = service.Bot().CtxWithReqJson(botCtx, &reqJson)
	g.Log().Info(ctx, tokenName+" access successfully with "+rawJson)
	var retMsg string
	// 异步执行
	if req.Async {
		go service.Command().TryCommand(botCtx, req.Command)
		retMsg = "async"
	} else {
		var catch bool
		catch, retMsg = service.Command().TryCommand(botCtx, req.Command)
		if !catch {
			err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil),
				"command not found")
			return
		}
	}
	// 响应
	res = &v1.CommandRes{
		Message: retMsg,
	}
	// 检查是否需要同步消息
	if !req.MessageSync || req.Async {
		return
	}
	if req.GroupId == 0 || !service.Group().IsBinding(botCtx, req.GroupId) {
		err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil),
			"group not binding")
		return
	}
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(botCtx) {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
			"permission denied")
		return
	}
	// 限速 一分钟只能发送 5 条消息
	if limit, _ := service.Util().AutoLimit(ctx,
		"send_msg", gconv.String(req.GroupId), 5, time.Minute); limit {
		err = gerror.NewCode(gcode.New(http.StatusTooManyRequests, "", nil),
			http.StatusText(http.StatusTooManyRequests))
		return
	}
	// 发送消息
	err = service.Bot().SendMessage(botCtx, "", 0, req.GroupId, retMsg, true)
	if err != nil {
		err = gerror.NewCode(gcode.New(http.StatusInternalServerError, "", nil),
			err.Error())
		return
	}
	return
}
