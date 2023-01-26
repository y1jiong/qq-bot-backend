package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
	"qq-bot-backend/internal/service"
)

type sBot struct{}

func New() *sBot {
	return &sBot{}
}

func init() {
	service.RegisterBot(New())
}

const (
	ctxKeyForWebSocket = "ws"
	ctxKeyForReqJson   = "reqJson"
)

type echoModel struct {
	Context    context.Context
	SuccessMsg string
}

var (
	echoPool = make(map[string]echoModel)
)

func (s *sBot) ctxWithWebSocket(parent context.Context, ws *ghttp.WebSocket) context.Context {
	return context.WithValue(parent, ctxKeyForWebSocket, ws)
}

func (s *sBot) webSocketFromCtx(ctx context.Context) *ghttp.WebSocket {
	if v := ctx.Value(ctxKeyForWebSocket); v != nil {
		return v.(*ghttp.WebSocket)
	}
	return nil
}

func (s *sBot) ctxWithReqJson(ctx context.Context, reqJson *sj.Json) context.Context {
	return context.WithValue(ctx, ctxKeyForReqJson, reqJson)
}

func (s *sBot) reqJsonFromCtx(ctx context.Context) *sj.Json {
	if v := ctx.Value(ctxKeyForReqJson); v != nil {
		return v.(*sj.Json)
	}
	return nil
}

func (s *sBot) Process(ctx context.Context, ws *ghttp.WebSocket, rawJson []byte, nextProcess func(ctx context.Context)) {
	reqJson, err := sj.NewJson(rawJson)
	if err != nil {
		return
	}
	// ctx 携带 websocket
	ctx = s.ctxWithWebSocket(ctx, ws)
	// ctx 携带 reqJson
	ctx = s.ctxWithReqJson(ctx, reqJson)
	// debug mode
	postType := s.GetPostType(ctx)
	if service.Cfg().IsDebugEnabled(ctx) && postType != "meta_event" {
		g.Log().Info(ctx, "\n", rawJson)
	}
	// 下一步执行
	nextProcess(ctx)
}

func (s *sBot) CatchEcho(ctx context.Context) (catch bool) {
	if echo := s.getEcho(ctx); echo != "" {
		if lastReq, ok := echoPool[echo]; ok {
			switch s.getEchoStatus(ctx) {
			case "ok":
				s.SendPlainMsg(lastReq.Context, lastReq.SuccessMsg)
			case "async":
				s.SendPlainMsg(lastReq.Context, "已提交 async 处理")
			case "failed":
				s.SendPlainMsg(lastReq.Context, s.getEchoFailedMsg(ctx))
			}
			// 用后即删，以防内存泄露
			delete(echoPool, echo)
			catch = true
			return
		}
	}
	return
}

func (s *sBot) getEcho(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("echo").MustString()
}

func (s *sBot) getEchoStatus(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("status").MustString()
}

func (s *sBot) getEchoFailedMsg(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("wording").MustString()
}

func (s *sBot) GetPostType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("post_type").MustString()
}

func (s *sBot) GetMsgType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("message_type").MustString()
}

func (s *sBot) GetRequestType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("request_type").MustString()
}

func (s *sBot) GetNoticeType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("notice_type").MustString()
}

func (s *sBot) GetSubType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("sub_type").MustString()
}

func (s *sBot) GetMessage(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("message").MustString()
}

func (s *sBot) GetUserId(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("user_id").MustInt64()
}

func (s *sBot) GetGroupId(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("group_id").MustInt64()
}

func (s *sBot) GetComment(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("comment").MustString()
}

func (s *sBot) GetFlag(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("flag").MustString()
}

func (s *sBot) IsGroupOwnerOrAdmin(ctx context.Context) (yes bool) {
	role := s.reqJsonFromCtx(ctx).Get("sender").Get("role").MustString()
	if role == "owner" || role == "admin" {
		yes = true
	}
	return
}

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
