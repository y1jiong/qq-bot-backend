package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sModule) AutoMute(ctx context.Context, groupId, userId int64) {
	// 缓存键名
	cacheKey := "MuteTimes.QQ=" + gconv.String(userId)
	// 过期时间
	expirationDuration := 16 * time.Hour
	timesVar, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	if timesVar == nil {
		// 第一次不禁言
		err = gcache.Set(ctx, cacheKey, 1, expirationDuration)
		if err != nil {
			g.Log().Warning(ctx, err)
		}
		return
	}
	times := timesVar.Int()
	// 多次撤回
	err = gcache.Set(ctx, cacheKey, times+1, expirationDuration)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// 最终禁言分钟数
	muteMinutes := 1
	// 执行幂次运算
	for i := 0; i < times; i++ {
		muteMinutes *= consts.BaseMuteMinutes
		// 不超过 30 天 30*24*60=43200
		if muteMinutes > 43199 {
			muteMinutes = 43199
			break
		}
	}
	// 禁言 BaseMuteMinutes^times 分钟
	service.Bot().MutePrototype(ctx, groupId, userId, muteMinutes*60)
}
