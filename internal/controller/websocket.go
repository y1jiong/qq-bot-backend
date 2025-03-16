package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"qq-bot-backend/internal/service"
	"strings"
	"sync"
)

var (
	Bot = cBot{}
)

type cBot struct{}

var (
	// wsUpGrader is the default up-grader configuration for websocket.
	wsUpGrader = websocket.Upgrader{
		// It does not check the origin in default, the application can do it itself.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (c *cBot) Websocket(r *ghttp.Request) {
	ctx := r.Context()
	ctx, span := gtrace.NewSpan(ctx, "controller.Bot.Websocket")
	spanEnd := sync.OnceFunc(func() { span.End() })
	defer func() {
		spanEnd()
	}()

	var (
		tokenName string
		botId     int64
	)
	{
		// 忽视前置的 Bearer 或 Token 进行鉴权
		authorizations := strings.Fields(r.Header.Get("Authorization"))
		if len(authorizations) < 2 {
			r.Response.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := authorizations[1]
		var pass bool
		pass, tokenName, _, botId = service.Token().IsCorrectToken(ctx, token)
		if !pass {
			// token debug 验证模式
			if service.Cfg().IsDebugEnabled(ctx) && token == service.Cfg().GetDebugToken(ctx) {
				pass = true
				if tokenName == "" {
					tokenName = "debug"
				}
			}
			if !pass {
				r.Response.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		// 试图从 Header 获取 botId
		if botId == 0 {
			botId = gconv.Int64(r.Header.Get("X-Self-Id"))
		}
		// 记录登录时间
		service.Token().UpdateLoginTime(ctx, token)
		span.SetAttributes(attribute.String("websocket.token_name", tokenName))
	}
	// 升级 WebSocket 协议
	conn, err := wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	g.Log().Info(ctx, tokenName+" connected")
	// context 携带 WebSocket 对象
	ctx = service.Bot().CtxWithWebSocket(ctx, conn)
	// 并发 ws 写锁
	ctx = service.Bot().CtxNewWebSocketMutex(ctx)

	var (
		waitJoin       = make(chan struct{})
		joinConnection = sync.OnceFunc(func() { close(waitJoin) })
	)
	go func() {
		<-waitJoin

		// 加入连接
		if botId == 0 {
			var err error
			botId, _, err = service.Bot().GetLoginInfo(ctx)
			if err != nil {
				g.Log().Warning(ctx, err)
			}
		}
		if botId != 0 {
			service.Bot().JoinConnection(ctx, botId)
			g.Log().Info(ctx, tokenName+"("+gconv.String(botId)+") joined connection")
		}
	}()

	spanEnd()

	// 消息循环
	for {
		joinConnection()

		var bytes []byte
		_, bytes, err = conn.ReadMessage()
		if err != nil {
			// 离开连接
			if botId != 0 {
				service.Bot().LeaveConnection(botId)
				g.Log().Info(ctx, tokenName+"("+gconv.String(botId)+") left connection")
			}
			g.Log().Info(ctx, tokenName+" disconnected")
			break
		}
		// new trace
		ctx := trace.ContextWithSpanContext(ctx, trace.SpanContext{})
		// 异步处理 WebSocket 请求
		go service.Bot().Process(ctx, bytes, service.Process().Process)
	}
}
