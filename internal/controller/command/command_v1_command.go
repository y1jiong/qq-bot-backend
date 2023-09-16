package command

import (
	"context"
	"encoding/json"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/crypto/gsha1"
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
	// 验证 token
	pass, tokenName, ownerId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(gcode.New(http.StatusForbidden, "", nil), "permission denied")
		return
	}
	// 验证消息是否过期
	msgTime := gtime.New(time.Unix(req.Timestamp, 0))
	if gtime.Now().Sub(msgTime) > 5*time.Second {
		err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil), "message expired")
		return
	}
	// 验证签名
	{
		// tokenName+token+timestamp+command+groupId 的 sha1 值的 base64 值
		s := tokenName + req.Token + gconv.String(req.Timestamp) + req.Command + gconv.String(req.GroupId)
		if gbase64.EncodeString(gsha1.Encrypt(s)) != req.Signature {
			err = gerror.NewCode(gcode.New(http.StatusBadRequest, "", nil), "signature error")
			return
		}
	}
	// 记录登录时间
	service.Token().UpdateLoginTime(ctx, req.Token)
	g.Log().Info(ctx, tokenName+" access successfully")
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
	rawJson, err := json.Marshal(innerReq)
	if err != nil {
		return
	}
	reqJson, err := sj.NewJson(rawJson)
	if err != nil {
		return
	}
	ctx = service.Bot().CtxWithReqJson(ctx, reqJson)
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
