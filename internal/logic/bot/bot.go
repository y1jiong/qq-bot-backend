package bot

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
	"qq-bot-backend/internal/service"
	"sync"
)

type sBot struct{}

func New() *sBot {
	return &sBot{}
}

func init() {
	service.RegisterBot(New())
}

const (
	ctxKeyWebSocketMutex = "ws.mutex"
	ctxKeyWebSocket      = "ws"
	ctxKeyReqJson        = "reqJson"

	cacheKeyMsgIdPrefix = "bot.message_id_"
)

func (s *sBot) CtxWithWebSocket(parent context.Context, conn *websocket.Conn) context.Context {
	return context.WithValue(parent, ctxKeyWebSocket, conn)
}

func (s *sBot) webSocketFromCtx(ctx context.Context) *websocket.Conn {
	if conn, ok := ctx.Value(ctxKeyWebSocket).(*websocket.Conn); ok {
		return conn
	}
	return nil
}

func (s *sBot) CtxNewWebSocketMutex(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKeyWebSocketMutex, &sync.Mutex{})
}

func (s *sBot) webSocketMutexFromCtx(ctx context.Context) *sync.Mutex {
	if mu, ok := ctx.Value(ctxKeyWebSocketMutex).(*sync.Mutex); ok {
		return mu
	}
	return nil
}

func (s *sBot) CtxWithReqJson(ctx context.Context, reqJson *ast.Node) context.Context {
	return context.WithValue(ctx, ctxKeyReqJson, reqJson)
}

func (s *sBot) reqJsonFromCtx(ctx context.Context) *ast.Node {
	if node, ok := ctx.Value(ctxKeyReqJson).(*ast.Node); ok {
		return node
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
	// message segment
	s.tryMessageSegmentToString(ctx)
	// debug mode
	if service.Cfg().IsDebugEnabled(ctx) && s.GetPostType(ctx) != "meta_event" {
		g.Log().Debug(ctx, "\n", rawJson)
	}
	// 捕捉 echo
	if s.catchEcho(ctx) {
		return
	}
	// 下一步执行
	nextProcess(ctx)
}
