package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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
	if service.Cfg().IsDebugEnabled(ctx) && s.GetPostType(ctx) != "meta_event" {
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

func (s *sBot) IsGroupOwnerOrAdmin(ctx context.Context) (yes bool) {
	role := s.reqJsonFromCtx(ctx).Get("sender").Get("role").MustString()
	if role == "owner" || role == "admin" {
		yes = true
	}
	return
}
