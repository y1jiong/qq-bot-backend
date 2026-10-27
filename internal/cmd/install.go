package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Install = gcmd.Command{
		Name:          "install",
		Brief:         "install and start service",
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s, err := newService()
			if err != nil {
				return
			}
			if err = s.InstallAndStart(); err != nil {
				return
			}
			g.Log().Notice(ctx, "安装并启动服务成功")
			return
		},
	}
	Uninstall = gcmd.Command{
		Name:          "uninstall",
		Brief:         "stop and uninstall service",
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s, err := newService()
			if err != nil {
				return
			}
			if err = s.StopAndUninstall(); err != nil {
				return
			}
			g.Log().Notice(ctx, "卸载服务成功")
			return
		},
	}
)
