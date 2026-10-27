package cmd

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/controller"
	"qq-bot-backend/internal/controller/command"
	"qq-bot-backend/internal/controller/message"
	"qq-bot-backend/internal/service"

	"github.com/gogf/gf/contrib/trace/otlpgrpc/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:          consts.ProjName,
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return runServer(ctx)
		},
	}
)

// runServer 启动 crontab 与 HTTP server，阻塞到 ctx 被取消（服务模式）或进程收到退出信号。
func runServer(ctx context.Context) (err error) {
	// OpenTelemetry
	if endpoint, traceToken := service.Cfg().GetOTel(ctx); endpoint != "" {
		var shutdown func(context.Context)
		shutdown, err = otlpgrpc.Init(consts.ProjName, endpoint, traceToken)
		if err != nil {
			return
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
	// 服务管理器停止服务时 kardianos 会取消 ctx，这里主动关闭 HTTP server，
	// 让 s.Run() 返回，Stop 才能在 stopTimeout 内等到业务函数退出
	if done := ctx.Done(); done != nil {
		go func() {
			<-done
			_ = s.Shutdown()
		}()
	}
	s.Run()

	return
}
