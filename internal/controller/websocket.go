package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"qq-bot-backend/internal/service"
	"strings"
)

var (
	ConnectClient = cConnectClient{}
)

type cConnectClient struct{}

func (c *cConnectClient) Connect(r *ghttp.Request) {
	ctx := r.Context()
	// 忽视前置的 Bearer 或 Token 进行鉴权
	authorizations := strings.Fields(r.Header.Get("Authorization"))
	if len(authorizations) < 2 {
		r.Response.WriteHeader(http.StatusForbidden)
		return
	}
	token := authorizations[1]
	// token 验证模式
	if service.Cfg().IsDebugEnabled(ctx) {
		// debug mode
		if token != service.Cfg().GetDebugToken(ctx) && !service.Token().IsCorrectToken(ctx, token) {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		if !service.Token().IsCorrectToken(ctx, token) {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
	}
	// 升级 WebSocket 协议
	ws, err := r.WebSocket()
	if err != nil {
		return
	}
	g.Log().Info(ctx, "Connected")
	for {
		var wsReq []byte
		_, wsReq, err = ws.ReadMessage()
		if err != nil {
			g.Log().Info(ctx, "Disconnected")
			return
		}
		// 异步处理 WebSocket 请求
		// ctx 携带 WebSocket 对象
		go service.Bot().Process(service.Bot().CtxWithWebSocket(ctx, ws), wsReq)
	}
}
