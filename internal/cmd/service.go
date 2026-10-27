package cmd

import (
	"context"
	"os"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/svc"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	// Run 由操作系统服务管理器调用（安装后的启动命令为 "<binary> run"），不要手动执行
	Run = gcmd.Command{
		Name:          "run",
		Brief:         "run in foreground, the way the OS service manager starts it",
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) error {
			s, err := newService()
			if err != nil {
				return err
			}
			return s.Run()
		},
	}
)

// newService 构造服务定义。WorkingDirectory 记录安装时的当前目录，
// 保证服务管理器启动进程后 GoFrame 仍能读到 manifest/config 下的配置。
func newService() (*svc.Svc, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return svc.New(
		consts.ProjName,
		"QQ Bot Backend",
		"QQ Bot Backend Service",
		wd,
		runAsService,
	)
}

// runAsService 是 kardianos service 的业务入口；ctx 取消即退出，
// 是否重启交给服务管理器（Restart/OnFailure）决定
func runAsService(ctx context.Context) {
	if err := runServer(ctx); err != nil {
		g.Log().Fatal(ctx, err)
	}
}
