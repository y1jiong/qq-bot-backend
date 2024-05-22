package util

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

func (s *sUtil) AutoLimit(ctx context.Context, kind, key string, limitTimes int, duration time.Duration) (
	limited bool, times int) {
	// 缓存键名
	cacheKey := "LimitTimes" + kind + key
	// 过期时间
	timesVar, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if timesVar == nil {
		// 设置缓存
		defaultTimes := 1
		err = gcache.Set(ctx, cacheKey, defaultTimes, duration)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		times = defaultTimes - 1
	} else {
		// 更新缓存
		times = timesVar.Int()
		_, _, err = gcache.Update(ctx, cacheKey, times+1)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	if times < limitTimes {
		return
	}
	limited = true
	return
}
