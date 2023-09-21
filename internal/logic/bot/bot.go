package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"qq-bot-backend/internal/service"
	"sync"
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
	ctxKeyForWebSocketMutex = "ws.mutex"
	ctxKeyForWebSocket      = "ws"
	ctxKeyForReqJson        = "reqJson"
	echoPrefix              = "echo_"
	echoDuration            = 60 * time.Second
	echoTimeout             = echoDuration + 10*time.Second
)

type echoModel struct {
	LastContext  context.Context
	CallbackFunc func(ctx context.Context, rsyncCtx context.Context)
}

func (s *sBot) CtxWithWebSocket(parent context.Context, ws *ghttp.WebSocket) context.Context {
	return context.WithValue(parent, ctxKeyForWebSocket, ws)
}

func (s *sBot) webSocketFromCtx(ctx context.Context) *ghttp.WebSocket {
	if v := ctx.Value(ctxKeyForWebSocket); v != nil {
		return v.(*ghttp.WebSocket)
	}
	return nil
}

func (s *sBot) CtxNewWebSocketMutex(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKeyForWebSocketMutex, &sync.Mutex{})
}

func (s *sBot) webSocketMutexFromCtx(ctx context.Context) *sync.Mutex {
	if v := ctx.Value(ctxKeyForWebSocketMutex); v != nil {
		return v.(*sync.Mutex)
	}
	return nil
}

func (s *sBot) CtxWithReqJson(ctx context.Context, reqJson *ast.Node) context.Context {
	return context.WithValue(ctx, ctxKeyForReqJson, reqJson)
}

func (s *sBot) reqJsonFromCtx(ctx context.Context) *ast.Node {
	if v := ctx.Value(ctxKeyForReqJson); v != nil {
		return v.(*ast.Node)
	}
	return nil
}

func (s *sBot) writeMessage(ctx context.Context, messageType int, data []byte) error {
	mu := s.webSocketMutexFromCtx(ctx)
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	return s.webSocketFromCtx(ctx).WriteMessage(messageType, data)
}

func (s *sBot) Process(ctx context.Context, rawJson []byte, nextProcess func(ctx context.Context)) {
	// 检查 context 中是否携带 WebSocket 对象
	if s.webSocketFromCtx(ctx) == nil {
		panic("context does not include websocket")
	}
	// ctx 携带 reqJson
	reqJson, err := sonic.Get(rawJson)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	ctx = s.CtxWithReqJson(ctx, &reqJson)
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

func (s *sBot) defaultEchoProcess(rsyncCtx context.Context) (err error) {
	if s.getEchoStatus(rsyncCtx) != "ok" {
		switch s.getEchoStatus(rsyncCtx) {
		case "async":
			err = errors.New("已提交 async 处理")
		case "failed":
			err = errors.New(s.getEchoFailedMsg(rsyncCtx))
		}
	}
	return
}

func (s *sBot) IsGroupOwnerOrAdmin(ctx context.Context) (yes bool) {
	role, _ := s.reqJsonFromCtx(ctx).Get("sender").Get("role").StrictString()
	// lazy load user role
	if role == "" {
		member, err := s.GetGroupMemberInfo(ctx, s.GetGroupId(ctx), s.GetUserId(ctx))
		if err != nil {
			g.Log().Warning(ctx, err)
			return
		}
		role, err = member.Get("role").StrictString()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		params := []ast.Pair{
			{
				Key:   "role",
				Value: ast.NewString(role),
			},
		}
		_, err = s.reqJsonFromCtx(ctx).Set("sender", ast.NewObject(params))
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	return role == "owner" || role == "admin"
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
