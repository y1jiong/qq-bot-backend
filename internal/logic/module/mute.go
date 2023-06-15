package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sModule) AutoMute(ctx context.Context, kind string, groupId, userId int64,
	passTimes, baseMinutes, limitMinutes int, duration time.Duration) {
	// 缓存键名
	cacheKey := "MuteTimes:" + kind + "=" + gconv.String(userId)
	// 过期时间
	timesVar, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	var times int
	if timesVar == nil {
		// 设置缓存
		defaultTimes := 1
		err = gcache.Set(ctx, cacheKey, defaultTimes, duration)
		if err != nil {
			g.Log().Warning(ctx, err)
			return
		}
		times = defaultTimes - 1
	} else {
		// 更新缓存
		times = timesVar.Int()
		_, _, err = gcache.Update(ctx, cacheKey, times+1)
		if err != nil {
			g.Log().Warning(ctx, err)
			return
		}
	}
	if times < passTimes {
		return
	}
	// 最终禁言分钟数
	muteMinutes := 1
	// 执行幂次运算
	for i := 0; i < times; i++ {
		muteMinutes *= baseMinutes
		if limitMinutes > 0 && muteMinutes > limitMinutes {
			muteMinutes = limitMinutes
			break
		}
		// 不超过 30 天 30*24*60=43200
		if muteMinutes > 43199 {
			muteMinutes = 43199
			break
		}
	}
	// 禁言 BaseMuteMinutes^times 分钟
	service.Bot().MutePrototype(ctx, groupId, userId, muteMinutes*60)
}
