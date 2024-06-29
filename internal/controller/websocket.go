package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"net/http"
	"qq-bot-backend/internal/service"
	"strings"
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
	var (
		tokenName string
		botId     int64
	)
	{
		// 忽视前置的 Bearer 或 Token 进行鉴权
		authorizations := strings.Fields(r.Header.Get("Authorization"))
		if len(authorizations) < 2 {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
		token := authorizations[1]
		var pass bool
		pass, tokenName, _, botId = service.Token().IsCorrectToken(ctx, token)
		if service.Cfg().IsDebugEnabled(ctx) {
			// token debug 验证模式
			if !pass && token != service.Cfg().GetDebugToken(ctx) {
				r.Response.WriteHeader(http.StatusForbidden)
				return
			}
			if tokenName == "" {
				tokenName = "debug"
			}
		} else {
			// token 正常验证模式
			if !pass {
				r.Response.WriteHeader(http.StatusForbidden)
				return
			}
		}
		// 记录登录时间
		service.Token().UpdateLoginTime(ctx, token)
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
	// 加入连接池
	if botId != 0 {
		service.Bot().JoinConnectionPool(ctx, botId)
		g.Log().Info(ctx, tokenName+"("+gconv.String(botId)+") joined connection pool")
	}
	// 消息循环
	for {
		var wsReq []byte
		_, wsReq, err = conn.ReadMessage()
		if err != nil {
			// 离开连接池
			if botId != 0 {
				service.Bot().LeaveConnectionPool(botId)
				g.Log().Info(ctx, tokenName+"("+gconv.String(botId)+") left connection pool")
			}
			g.Log().Info(ctx, tokenName+" disconnected")
			break
		}
		// 异步处理 WebSocket 请求
		go service.Bot().Process(ctx, wsReq, service.Process().Process)
	}
}
