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
	pass, tokenName, ownerId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil), "permission denied")
		return
	}
	// 验证签名
	{
		// 以 token+command+group_id+timestamp 为原文，以 token_name 为 key 的 HmacSha1 值的 base64 值
		s := req.Token + req.Command + gconv.String(req.GroupId) + gconv.String(req.Timestamp)
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
	g.Log().Info(ctx, tokenName+" access successfully with "+string(rawJson))
	if err != nil {
		return
	}
	reqJson, _ := sonic.Get(rawJson)
	ctx = service.Bot().CtxWithReqJson(ctx, &reqJson)
	// 处理命令
	catch, retMsg := service.Command().TryCommand(ctx)
	if !catch {
		err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil), "command not found")
		return
	}
	// 响应
	res = &v1.CommandRes{
		Message: retMsg,
	}
	return
}
