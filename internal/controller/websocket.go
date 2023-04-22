package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"qq-bot-backend/internal/service"
	"strings"
)

var (
	Bot = cBot{}
)

type cBot struct{}

func (c *cBot) Websocket(r *ghttp.Request) {
	ctx := r.Context()
	// 忽视前置的 Bearer 或 Token 进行鉴权
	authorizations := strings.Fields(r.Header.Get("Authorization"))
	if len(authorizations) < 2 {
		r.Response.WriteHeader(http.StatusForbidden)
		return
	}
	token := authorizations[1]
	var tokenName string
	if service.Cfg().IsEnabledDebug(ctx) {
		// token debug 验证模式
		var pass bool
		tokenName = "debug"
		pass, tokenName = service.Token().IsCorrectToken(ctx, token)
		// debug mode
		if !pass && token != service.Cfg().GetDebugToken(ctx) {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		// token 正常验证模式
		var pass bool
		pass, tokenName = service.Token().IsCorrectToken(ctx, token)
		if !pass {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
	}
	// 升级 WebSocket 协议
	ws, err := r.WebSocket()
	if err != nil {
		return
	}
	g.Log().Info(ctx, tokenName+" Connected")
	for {
		var wsReq []byte
		_, wsReq, err = ws.ReadMessage()
		if err != nil {
			g.Log().Info(ctx, tokenName+" Disconnected")
			return
		}
		// 异步处理 WebSocket 请求
		go service.Bot().Process(ctx, ws, wsReq, service.Process().Process)
	}
}
