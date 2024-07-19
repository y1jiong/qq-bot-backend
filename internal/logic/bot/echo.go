package bot

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

const (
	echoPrefix   = "echo_"
	echoDuration = 60 * time.Second
	echoTimeout  = echoDuration + 10*time.Second
)

type echoModel struct {
	LastContext  context.Context
	CallbackFunc func(ctx context.Context, rsyncCtx context.Context)
	TimeoutFunc  func(ctx context.Context)
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
		echo.CallbackFunc(echo.LastContext, ctx)
	}
	return
}

func (s *sBot) defaultEchoProcess(rsyncCtx context.Context) error {
	if s.getEchoStatus(rsyncCtx) != "ok" {
		switch s.getEchoStatus(rsyncCtx) {
		case "async":
			return errors.New("已提交 async 处理")
		case "failed":
			return errors.New(s.getEchoFailedMsg(rsyncCtx))
		}
	}
	return nil
}

func (s *sBot) pushEchoCache(ctx context.Context, echoSign string,
	callbackFunc func(ctx context.Context, rsyncCtx context.Context),
	timeoutFunc func(ctx context.Context)) error {
	echoKey := echoPrefix + echoSign
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
		echo := v.Val().(*echoModel)
		if echo == nil || echo.TimeoutFunc == nil {
			return
		}
		// 执行超时回调
		echo.TimeoutFunc(echo.LastContext)
	}(ctx, echoKey)
	return nil
}

func (s *sBot) popEchoCache(ctx context.Context, echoSign string) (echo *echoModel, err error) {
	echoKey := echoPrefix + echoSign
	contain, err := gcache.Contains(ctx, echoKey)
	if err != nil || !contain {
		return
	}
	v, err := gcache.Remove(ctx, echoKey)
	if err != nil {
		return
	}
	echo = v.Val().(*echoModel)
	return
}
