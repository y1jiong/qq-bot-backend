// Package svc 通过 github.com/kardianos/service 管理 qq-bot-backend 的系统服务
// （Linux systemd、macOS launchd、Windows SCM 由该库按平台适配）。
package svc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kardianos/service"
)

// stopTimeout 限制 Stop 等待业务函数返回的时间。
// systemd 默认 TimeoutStopSec 为 90s、Windows SCM 约 30s 后强制结束进程，
// 等待更久只会拖到被强杀，反而不干净。
var stopTimeout = 10 * time.Second

// Svc 包装一个 kardianos service 实例。
type Svc struct {
	inner service.Service
}

type program struct {
	runFunc func(context.Context)
	cancel  context.CancelFunc
	done    chan struct{}
}

func (p *program) Start(_ service.Service) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.done = make(chan struct{})

	go func() {
		defer close(p.done)
		p.runFunc(ctx)
	}()
	return nil
}

func (p *program) Stop(_ service.Service) error {
	// 业务函数是长驻循环（HTTP server），只有 ctx 被取消才会返回。
	// 这里等它退出，但设上限：无限阻塞会让服务管理器等到超时后发 SIGKILL。
	if p.cancel != nil {
		p.cancel()
	}
	if p.done == nil {
		return nil
	}

	select {
	case <-p.done:
	case <-time.After(stopTimeout):
		fmt.Fprintf(os.Stderr, "service did not stop within %v, exiting anyway\n", stopTimeout)
	}
	return nil
}

// config 描述交给操作系统服务管理器的服务信息。
func config(name, displayName, description, workingDirectory string) *service.Config {
	return &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		// 服务管理器执行：<binary> run；main() 识别到 run 子命令后进入服务模式
		Arguments: []string{"run"},
		// 记录安装时的工作目录，服务管理器从别处启动时 GoFrame 仍能读到配置
		WorkingDirectory: workingDirectory,
		// 仅 systemd 使用：直接写入 unit 的 [Unit] 段
		Dependencies: []string{
			"Wants=network-online.target",
			"After=network-online.target",
		},
		Option: service.KeyValue{
			// 进程异常退出时由服务管理器拉起；正常 stop 不会触发重启
			"Restart":     "on-failure",
			"LimitNOFILE": 65535,
			// Windows：服务未上报停止状态即视为失败，120s 后重启；
			// 连续运行一天无失败则重置失败计数。其他平台忽略这些键。
			"OnFailure":              "restart",
			"OnFailureDelayDuration": "120s",
			"OnFailureResetPeriod":   86400,
		},
	}
}

// New 创建服务。workingDirectory 是安装时的工作目录；runFunc 在服务模式下运行，
// 必须响应 ctx 取消并尽快返回，Stop 会取消 ctx 并最多等待 stopTimeout。
func New(name, displayName, description, workingDirectory string, runFunc func(context.Context)) (*Svc, error) {
	s, err := service.New(&program{runFunc: runFunc}, config(name, displayName, description, workingDirectory))
	if err != nil {
		return nil, err
	}
	return &Svc{inner: s}, nil
}

// InstallAndStart 安装服务并立即启动。
func (s *Svc) InstallAndStart() error {
	if err := s.inner.Install(); err != nil {
		return err
	}
	return s.inner.Start()
}

// StopAndUninstall 先停止服务再卸载；服务未安装或未在运行时 Stop 的报错忽略，
// 避免因为“本来就没在跑”而卸载不掉。
func (s *Svc) StopAndUninstall() error {
	_ = s.inner.Stop()
	return s.inner.Uninstall()
}

// Run 以服务模式运行，由操作系统服务管理器调用，阻塞到服务被停止。
func (s *Svc) Run() error {
	return s.inner.Run()
}
