package bot

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"qq-bot-backend/utility/segment"
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
	ctxKeyWebSocketMutex = "bot_ws_mutex"
	ctxKeyWebSocket      = "bot_ws"
	ctxKeyReq            = "bot_req"

	cacheKeyMsgIdPrefix = "bot_msg_id_"

	messageLengthLimit = 4500
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

func (s *sBot) CtxWithReqNode(ctx context.Context, req *ast.Node) context.Context {
	return context.WithValue(ctx, ctxKeyReq, req)
}

func (s *sBot) reqNodeFromCtx(ctx context.Context) *ast.Node {
	if node, ok := ctx.Value(ctxKeyReq).(*ast.Node); ok {
		return node
	}
	return nil
}

func (s *sBot) CloneReqNode(ctx context.Context) *ast.Node {
	j, err := s.reqNodeFromCtx(ctx).MarshalJSON()
	if err != nil {
		return nil
	}
	node, err := sonic.Get(j)
	if err != nil {
		return nil
	}

	allowedFields := map[string]struct{}{
		"message":        {},
		"self_id":        {},
		"user_id":        {},
		"group_id":       {},
		"sub_type":       {},
		"post_type":      {},
		"raw_message":    {},
		"message_type":   {},
		"message_format": {},
	}

	fields, _ := node.Map()
	for field := range fields {
		if _, ok := allowedFields[field]; !ok {
			_, _ = node.Unset(field)
		}
	}

	return &node
}

func (s *sBot) buildMessageNode(userId int64, nickname, message string) map[string]any {
	return map[string]any{
		"type": "node",
		"data": map[string]any{
			"user_id":  userId,
			"nickname": nickname,
			"content":  segment.ParseMessage(message),
		},
	}
}

func (s *sBot) MessageToNodes(userId int64, nickname, message string) []map[string]any {
	runes := []rune(message)
	nodes := make([]map[string]any, 0, (len(runes)+messageLengthLimit-1)/messageLengthLimit)
	for len(runes) > messageLengthLimit {
		length := utility.FindNaturalBreakPoint(runes, messageLengthLimit, 100)
		nodes = append(nodes, s.buildMessageNode(userId, nickname, string(runes[:length])))
		runes = runes[length:]
	}
	return append(nodes, s.buildMessageNode(userId, nickname, string(runes)))
}

func (s *sBot) writeMessage(ctx context.Context, messageType int, data []byte) error {
	if mu := s.webSocketMutexFromCtx(ctx); mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	return s.webSocketFromCtx(ctx).WriteMessage(messageType, data)
}

func (s *sBot) Process(ctx context.Context, rawJSON []byte, nextProcess func(ctx context.Context)) {
	// 检查 context 中是否携带 WebSocket 对象
	if s.webSocketFromCtx(ctx) == nil {
		panic("context does not include websocket")
	}
	// ctx 携带 reqNode
	reqNode, err := sonic.Get(rawJSON)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	ctx = s.CtxWithReqNode(ctx, &reqNode)
	// message segment
	s.tryMessageSegmentToString(ctx)
	// debug mode
	if service.Cfg().IsDebugEnabled(ctx) && s.GetPostType(ctx) != "meta_event" {
		g.Log().Debug(ctx, "\n", rawJSON)
	}
	// 捕捉 echo
	if s.catchEcho(ctx) {
		return
	}
	// 下一步执行
	nextProcess(ctx)
}
