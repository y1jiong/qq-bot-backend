package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"he3-bot/internal/service"
	"net/http"
	"strings"
)

var (
	ConnectClient = cConnectClient{}
)

type cConnectClient struct{}

func (c *cConnectClient) Connect(r *ghttp.Request) {
	ctx := r.Context()
	// 忽视前置的 Bearer 或 Token 进行鉴权
	authorization := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authorization) < 2 || authorization[1] != service.Cfg().GetAuthToken(ctx) {
		r.Response.WriteHeader(http.StatusForbidden)
		return
	}
	ws, err := r.WebSocket()
	if err != nil {
		return
	}
	glog.Info(ctx, "Connected")
	for {
		var msg []byte
		_, msg, err = ws.ReadMessage()
		if err != nil {
			glog.Info(ctx, "Disconnected")
			return
		}
		go service.Bot().Parse(ctx, ws, msg)
	}
}
