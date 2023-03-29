package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/controller"
)

var (
	Main = gcmd.Command{
		Name:          consts.ProjName,
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/ws", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(controller.Bot.Websocket)
			})
			s.Run()
			return
		},
	}
)
