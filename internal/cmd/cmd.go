package cmd

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"he3-bot/internal/consts"
	"he3-bot/internal/controller"
	"he3-bot/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:        "he3-bot",
		Description: "Version: " + consts.Version + "\n" + "Build Time: " + consts.BuildTime,
		Brief:       "start bot",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/v1/ws", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().Common)
				group.Bind(controller.ConnectClient.Connect)
			})
			s.Run()
			return nil
		},
	}
)
