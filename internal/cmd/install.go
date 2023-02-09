package cmd

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"os"
	"qq-bot-backend/internal/consts"
	"runtime"
)

const (
	installPath = "/etc/systemd/system/" + consts.ProjName + ".service"
)

var (
	Install = gcmd.Command{
		Name:          "install",
		Brief:         "install service",
		CaseSensitive: true,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			if isWindows() {
				err = errors.New("windows 暂不支持安装到系统")
				return
			}
			// 注册系统服务
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			serviceContent := []byte(
				"[Unit]\n" +
					"Description=" + consts.ProjName + " Service\n" +
					"After=network.target\n\n" +
					"[Service]\n" +
					"Type=simple\n" +
					"WorkingDirectory=" + wd +
					"\nExecStart=" + wd + "/" + consts.ProjName + " --gf.gerror.brief=true" +
					"\nRestart=on-failure\n" +
					"RestartSec=2\n\n" +
					"[Install]\n" +
					"WantedBy=multi-user.target\n")
			err = os.WriteFile(installPath, serviceContent, 0600)
			if err != nil {
				return err
			}
			g.Log().Notice(ctx, "安装服务成功\n可以使用 systemctl 管理 "+consts.ProjName+" 服务了")
			return
		},
	}
	Uninstall = gcmd.Command{
		Name:  "uninstall",
		Brief: "uninstall service",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			if isWindows() {
				err = errors.New("windows 暂不支持安装到系统")
				return
			}
			err = os.Remove(installPath)
			if err != nil {
				return err
			}
			g.Log().Notice(ctx, "卸载服务成功")
			return
		},
	}
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}
