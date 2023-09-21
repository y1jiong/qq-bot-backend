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
	// 验证消息是否过期
	msgTime := gtime.New(time.Unix(req.Timestamp, 0))
	if gtime.Now().Sub(msgTime) > 5*time.Second {
		err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil), "message expired")
		return
	}
	// 验证 token
	pass, tokenName, ownerId, botId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil), "permission denied")
		return
	}
	// 验证签名
	{
		// 以 token+command+group_id+timestamp+message_sync 为原文，以 token_name 为 key 的 HmacSha1 值的 base64 值
		s := req.Token + req.Command + gconv.String(req.GroupId) +
			gconv.String(req.Timestamp) + gconv.String(req.MessageSync)
		// HmacSha1
		hmacSha1 := hmac.New(sha1.New, []byte(tokenName))
		hmacSha1.Write([]byte(s))
		macBase64 := gbase64.Encode(hmacSha1.Sum(nil))
		if !hmac.Equal(macBase64, []byte(req.Signature)) {
			err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil), "signature error")
			return
		}
	}
	// 记录登录时间
	service.Token().UpdateLoginTime(ctx, req.Token)
	// 初始化内部请求
	innerReq := struct {
		Message string `json:"message"`
		UserId  int64  `json:"user_id"`
		GroupId int64  `json:"group_id"`
	}{
		Message: req.Command,
		UserId:  ownerId,
		GroupId: req.GroupId,
	}
	rawJson, err := sonic.ConfigStd.Marshal(innerReq)
	if err != nil {
		return
	}
	reqJson, _ := sonic.Get(rawJson)
	// 加载 botId 对应的 botCtx
	botCtx := service.Bot().LoadConnectionPool(botId)
	if botCtx == nil {
		err = gerror.NewCode(gcode.New(http.StatusInternalServerError, "", nil),
			"bot not connected")
		return
	}
	g.Log().Info(ctx, tokenName+" access successfully with "+string(rawJson))
	botCtx = service.Bot().CtxWithReqJson(botCtx, &reqJson)
	// 处理命令
	catch, retMsg := service.Command().TryCommand(botCtx)
	if !catch {
		err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil), "command not found")
		return
	}
	// 响应
	res = &v1.CommandRes{
		Message: retMsg,
	}
	// 检查是否需要同步消息
	if !req.MessageSync {
		return
	}
	if req.GroupId == 0 || !service.Group().IsBinding(botCtx, req.GroupId) {
		err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil),
			"group not binding")
		return
	}
	if !service.Bot().IsGroupOwnerOrAdmin(botCtx) {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil),
			"permission denied")
		return
	}
	// 发送消息
	service.Bot().SendMessage(botCtx, "", 0, req.GroupId, retMsg, true)
	return
}
