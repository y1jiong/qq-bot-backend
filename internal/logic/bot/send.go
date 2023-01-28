package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
)

func (s *sBot) SendPlainMsg(ctx context.Context, msg string) {
	s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, true)
}

func (s *sBot) SendMsg(ctx context.Context, msg string) {
	s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, false)
}

func (s *sBot) SendMessage(ctx context.Context, messageType string, uid, gid int64, msg string, plain bool) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "send_msg")
	// 参数
	params := make(map[string]any)
	params["message_type"] = messageType
	params["message"] = msg
	if plain {
		// 以纯文本方式发送
		params["auto_escape"] = true
	}
	if uid == 0 && gid == 0 {
		return
	}
	if uid != 0 {
		params["user_id"] = uid
	}
	if gid != 0 {
		params["group_id"] = gid
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) ApproveAddGroup(ctx context.Context, flag, subType string, approve bool, reason string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "set_group_add_request")
	// 参数
	params := make(map[string]any)
	params["flag"] = flag
	params["sub_type"] = subType
	params["approve"] = approve
	// 当不予通过时，给出理由
	if !approve {
		params["reason"] = reason
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// 发送响应
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SetModel(ctx context.Context, model string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "_set_model_show")
	// echo beacon
	echoBeacon := guid.S()
	resJson.Set("echo", echoBeacon)
	// 参数
	params := make(map[string]any)
	params["model"] = model
	params["model_show"] = model
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// 发送响应
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	// echo
	echoPool[echoBeacon] = echoModel{
		Context:    ctx,
		SuccessMsg: "已更改机型为 '" + model + "'",
	}
}

func (s *sBot) RevokeMessage(ctx context.Context, msgId int64) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "delete_msg")
	// 参数
	params := make(map[string]any)
	params["message_id"] = msgId
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}
