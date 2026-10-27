package svc

import (
	"context"
	"errors"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/kardianos/service"
)

// Stop 必须取消业务 ctx 并等业务函数返回后再返回：无限阻塞会让 systemd 熬到
// TimeoutStopSec（SIGKILL），也会让 Windows SCM 判定服务无响应。
func TestStopCancelsRunFunc(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		started := make(chan struct{})
		returned := make(chan struct{})

		p := &program{runFunc: func(ctx context.Context) {
			close(started)
			defer close(returned)
			<-ctx.Done()
		}}

		if err := p.Start(nil); err != nil {
			t.Fatalf("Start: %v", err)
		}
		<-started

		done := make(chan error, 1)
		go func() { done <- p.Stop(nil) }()

		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("Stop: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Stop blocked; it must return after cancelling the run function's context")
		}

		select {
		case <-returned:
		case <-time.After(time.Second):
			t.Fatal("Stop returned before the run function returned")
		}
	})
}

// 业务函数无视 ctx 时，Stop 也不能被永久拖住。
func TestStopIsBounded(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		orig := stopTimeout
		stopTimeout = 50 * time.Millisecond
		defer func() { stopTimeout = orig }()

		blocked := make(chan struct{})
		p := &program{runFunc: func(ctx context.Context) { <-blocked }}
		if err := p.Start(nil); err != nil {
			t.Fatalf("Start: %v", err)
		}
		defer close(blocked)

		start := time.Now()
		if err := p.Stop(nil); err != nil {
			t.Fatalf("Stop: %v", err)
		}
		if elapsed := time.Since(start); elapsed != stopTimeout {
			t.Fatalf("Stop waited %v, want exactly %v", elapsed, stopTimeout)
		}
	})
}

// Stop 可以在没有 Start 的情况下被调用（例如先卸载从未安装过的服务）。
func TestStopWithoutStart(t *testing.T) {
	p := &program{}
	if err := p.Stop(nil); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

// 服务管理器调用的命令行必须是 "<binary> run"，且崩溃重启相关的键名与库读取的一致。
func TestConfig(t *testing.T) {
	cfg := config("n", "display", "description", "/wd")

	if want := []string{"run"}; !slices.Equal(cfg.Arguments, want) {
		t.Errorf("Arguments = %v, want %v", cfg.Arguments, want)
	}
	if cfg.WorkingDirectory != "/wd" {
		t.Errorf("WorkingDirectory = %q, want /wd", cfg.WorkingDirectory)
	}
	for key, want := range map[string]any{
		"Restart":                "on-failure",
		"OnFailure":              "restart",
		"OnFailureDelayDuration": "120s",
		"OnFailureResetPeriod":   86400,
	} {
		if got := cfg.Option[key]; got != want {
			t.Errorf("Option[%q] = %v, want %v", key, got, want)
		}
	}
}

// fakeService 记录调用顺序，用来验证 install/uninstall 的组合语义。
type fakeService struct {
	calls   []string
	stopErr error
}

func (f *fakeService) Run() error     { f.calls = append(f.calls, "run"); return nil }
func (f *fakeService) Start() error   { f.calls = append(f.calls, "start"); return nil }
func (f *fakeService) Stop() error    { f.calls = append(f.calls, "stop"); return f.stopErr }
func (f *fakeService) Restart() error { f.calls = append(f.calls, "restart"); return nil }
func (f *fakeService) Install() error { f.calls = append(f.calls, "install"); return nil }
func (f *fakeService) Uninstall() error {
	f.calls = append(f.calls, "uninstall")
	return nil
}
func (f *fakeService) Logger(_ chan<- error) (service.Logger, error)       { return nil, nil }
func (f *fakeService) SystemLogger(_ chan<- error) (service.Logger, error) { return nil, nil }
func (f *fakeService) String() string                                      { return "fake" }
func (f *fakeService) Platform() string                                    { return "fake" }
func (f *fakeService) Status() (service.Status, error)                     { return service.StatusUnknown, nil }

// install 必须同时把服务启动起来。
func TestInstallAndStart(t *testing.T) {
	f := &fakeService{}
	s := &Svc{inner: f}

	if err := s.InstallAndStart(); err != nil {
		t.Fatalf("InstallAndStart: %v", err)
	}
	if want := []string{"install", "start"}; !slices.Equal(f.calls, want) {
		t.Errorf("calls = %v, want %v", f.calls, want)
	}
}

// uninstall 必须先 stop；服务本来就没在跑（stop 报错）也不能阻塞卸载。
func TestStopAndUninstallStopsFirst(t *testing.T) {
	f := &fakeService{stopErr: errors.New("service is not running")}
	s := &Svc{inner: f}

	if err := s.StopAndUninstall(); err != nil {
		t.Fatalf("StopAndUninstall: %v", err)
	}
	if want := []string{"stop", "uninstall"}; !slices.Equal(f.calls, want) {
		t.Errorf("calls = %v, want %v", f.calls, want)
	}
}
