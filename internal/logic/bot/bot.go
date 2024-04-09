package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gorilla/websocket"
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

var (
	connectionPool = sync.Map{}
)

type echoModel struct {
	LastContext  context.Context
	CallbackFunc func(ctx context.Context, rsyncCtx context.Context)
	TimeoutFunc  func(ctx context.Context)
}

func (s *sBot) CtxWithWebSocket(parent context.Context, conn *websocket.Conn) context.Context {
	return context.WithValue(parent, ctxKeyForWebSocket, conn)
}

func (s *sBot) webSocketFromCtx(ctx context.Context) *websocket.Conn {
	if v := ctx.Value(ctxKeyForWebSocket); v != nil {
		return v.(*websocket.Conn)
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

func (s *sBot) JoinConnectionPool(ctx context.Context, key int64) {
	connectionPool.Store(key, ctx)
}

func (s *sBot) LeaveConnectionPool(key int64) {
	connectionPool.Delete(key)
}

func (s *sBot) LoadConnectionPool(key int64) context.Context {
	if v, ok := connectionPool.Load(key); ok {
		return v.(context.Context)
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
		catch = true
		if echo.CallbackFunc == nil {
			if echo.TimeoutFunc != nil {
				echo.TimeoutFunc(echo.LastContext)
			}
			return
		}
		echo.CallbackFunc(echo.LastContext, ctx)
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

func (s *sBot) pushEchoCache(ctx context.Context, echoSign string,
	callbackFunc func(ctx context.Context, rsyncCtx context.Context),
	timeoutFunc func(ctx context.Context)) (err error) {
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
		v, e := gcache.Remove(ctx, echoKey)
		if e != nil {
			g.Log().Error(ctx, e)
			return
		}
		// 执行超时回调
		if v == nil {
			return
		}
		var echo *echoModel
		e = v.Scan(&echo)
		if e != nil {
			g.Log().Error(ctx, e)
			return
		}
		if echo == nil || echo.TimeoutFunc == nil {
			return
		}
		echo.TimeoutFunc(echo.LastContext)
	}()
	// 放入缓存
	err = gcache.Set(ctx, echoKey, echoModel{
		LastContext:  ctx,
		CallbackFunc: callbackFunc,
		TimeoutFunc:  timeoutFunc,
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
