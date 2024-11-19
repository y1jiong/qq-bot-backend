package bot

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcache"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	echoPrefix   = "echo_"
	echoDuration = 60 * time.Second
	echoTimeout  = echoDuration + 10*time.Second
)

type echoModel struct {
	LastContext  context.Context
	CallbackFunc func(ctx context.Context, asyncCtx context.Context)
	TimeoutFunc  func(ctx context.Context)
}

func getEchoCacheKey(echoSign string) string {
	return echoPrefix + echoSign
}

func (s *sBot) catchEcho(ctx context.Context) (catch bool) {
	if echoSign := s.getEcho(ctx); echoSign != "" {
		echo, err := s.popEchoCache(ctx, echoSign)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		if echo == nil {
			return
		}
		catch = true
		if echo.CallbackFunc == nil {
			return
		}

		var span trace.Span
		ctx, span = gtrace.NewSpan(s.extractEchoSign(ctx, echoSign), "bot.echo")
		defer span.End()
		echo.CallbackFunc(echo.LastContext, ctx)
	}
	return
}

func (s *sBot) defaultEchoHandler(asyncCtx context.Context) error {
	if s.getEchoStatus(asyncCtx) != "ok" {
		switch s.getEchoStatus(asyncCtx) {
		case "async":
			return errors.New("已提交 async 处理")
		case "failed":
			return errors.New(s.getEchoFailedMsg(asyncCtx))
		}
	}
	return nil
}

func (s *sBot) pushEchoCache(ctx context.Context, echoSign string,
	callbackFunc func(ctx context.Context, asyncCtx context.Context),
	timeoutFunc func(ctx context.Context)) error {
	if callbackFunc == nil || timeoutFunc == nil {
		return errors.New("callbackFunc or timeoutFunc must not be nil")
	}
	echoKey := getEchoCacheKey(echoSign)
	// 放入缓存
	if err := gcache.Set(ctx, echoKey, &echoModel{
		LastContext:  ctx,
		CallbackFunc: callbackFunc,
		TimeoutFunc:  timeoutFunc,
	}, echoTimeout); err != nil {
		return err
	}
	// 检查超时
	go func(ctx context.Context, echoKey string) {
		time.Sleep(echoDuration)
		contain, err := gcache.Contains(ctx, echoKey)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		if !contain {
			return
		}
		v, err := gcache.Remove(ctx, echoKey)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		if v == nil {
			return
		}
		echo, ok := v.Val().(*echoModel)
		if !ok {
			return
		}
		if echo == nil || echo.TimeoutFunc == nil {
			return
		}
		// 执行超时回调
		echo.TimeoutFunc(echo.LastContext)
	}(ctx, echoKey)
	return nil
}

func (s *sBot) popEchoCache(ctx context.Context, echoSign string) (echo *echoModel, err error) {
	echoKey := getEchoCacheKey(echoSign)
	contain, err := gcache.Contains(ctx, echoKey)
	if err != nil || !contain {
		return
	}
	v, err := gcache.Remove(ctx, echoKey)
	if err != nil {
		return
	}
	echo, ok := v.Val().(*echoModel)
	if !ok {
		return nil, errors.New("echo model type error")
	}
	return
}
