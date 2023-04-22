package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"qq-bot-backend/internal/service"
	"time"
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
	echoPrefix         = "echo_"
	echoDuration       = 60 * time.Second
	echoTimeout        = echoDuration + 10*time.Second
)

type echoModel struct {
	LastContext  context.Context
	CallbackFunc func(ctx context.Context, rsyncCtx context.Context)
}

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
	if service.Cfg().IsEnabledDebug(ctx) && s.GetPostType(ctx) != "meta_event" {
		g.Log().Info(ctx, "\n", rawJson)
	}
	// 捕捉 echo
	if s.catchEcho(ctx) {
		return
	}
	// 下一步执行
	nextProcess(ctx)
}

func (s *sBot) catchEcho(ctx context.Context) (catch bool) {
	if echoSign := s.getEcho(ctx); echoSign != "" {
		echo, err := s.popEchoCache(ctx, echoSign)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		if echo == nil {
			return
		}
		echo.CallbackFunc(echo.LastContext, ctx)
		catch = true
	}
	return
}

func (s *sBot) DefaultEchoProcess(ctx context.Context, rsyncCtx context.Context) {
	switch s.GetEchoStatus(rsyncCtx) {
	case "async":
		s.SendPlainMsg(ctx, "已提交 async 处理")
	case "failed":
		s.SendPlainMsg(ctx, s.GetEchoFailedMsg(rsyncCtx))
	}
}

func (s *sBot) IsGroupOwnerOrAdmin(ctx context.Context) (yes bool) {
	role := s.reqJsonFromCtx(ctx).Get("sender").Get("role").MustString()
	if role == "owner" || role == "admin" {
		yes = true
	}
	return
}

func (s *sBot) pushEchoCache(ctx context.Context, echoSign string, callbackFunc func(ctx context.Context, rsyncCtx context.Context)) (err error) {
	echoKey := echoPrefix + echoSign
	// 检查超时
	go func() {
		time.Sleep(echoDuration)
		contain, e := gcache.Contains(ctx, echoKey)
		if e != nil {
			g.Log().Error(ctx, e)
			return
		}
		if !contain {
			return
		}
		_, e = gcache.Remove(ctx, echoKey)
		if e != nil {
			g.Log().Error(ctx, e)
		}
		s.SendPlainMsg(ctx, "echo 超时")
	}()
	// 放入缓存
	err = gcache.Set(ctx, echoKey, echoModel{
		LastContext:  ctx,
		CallbackFunc: callbackFunc,
	}, echoTimeout)
	return
}

func (s *sBot) popEchoCache(ctx context.Context, echoSign string) (echo *echoModel, err error) {
	echoKey := echoPrefix + echoSign
	contain, err := gcache.Contains(ctx, echoKey)
	if err != nil || !contain {
		return
	}
	v, err := gcache.Remove(ctx, echoKey)
	if err != nil {
		return
	}
	err = v.Scan(&echo)
	return
}
