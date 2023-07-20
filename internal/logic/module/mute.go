package module

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sModule) AutoMute(ctx context.Context, kind string, groupId, userId int64,
	limitTimes, baseMinutes, limitMinutes int, duration time.Duration) {
	limited, times := s.AutoLimit(ctx, kind, gconv.String(userId), limitTimes, duration)
	if !limited {
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
