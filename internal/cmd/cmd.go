package cmd

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/controller"
	"qq-bot-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:          consts.ProjName,
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/v1/ws", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().Common)
				group.Bind(controller.Bot.Websocket)
			})
			s.Run()
			return
		},
	}
)
