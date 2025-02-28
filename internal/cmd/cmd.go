package cmd

import (
	"context"
	"github.com/gogf/gf/contrib/trace/otlpgrpc/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/controller"
	"qq-bot-backend/internal/controller/command"
	"qq-bot-backend/internal/controller/message"
	"qq-bot-backend/internal/service"
)

var (
	Main = gcmd.Command{
		Name:          consts.ProjName,
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// OpenTelemetry
			if endpoint, traceToken := service.Cfg().GetOTelConfig(ctx); endpoint != "" {
				shutdown, err := otlpgrpc.Init(consts.ProjName, endpoint, traceToken)
				if err != nil {
					return err
				}
				defer shutdown(ctx)
			}

			// crontab
			go service.Crontab().Run(ctx)

			// http server
			s := g.Server()
			s.Group("/ws", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(controller.Bot.Websocket)
			})
			s.Group("/file/{id}", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().ErrCodeToHttpStatus)
				group.Bind(controller.File.GetCachedFileById)
			})
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareCORS, ghttp.MiddlewareHandlerResponse,
					service.Middleware().ErrCodeToHttpStatus, service.Middleware().RateLimit,
				)
				group.Group("/v1", func(group *ghttp.RouterGroup) {
					group.Bind(
						command.NewV1(),
						message.NewV1(),
					)
				})
			})
			s.Run()

			return
		},
	}
)
