package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"strings"
	"time"
)

func (s *sModule) TryKeywordRevoke(ctx context.Context) (catch bool) {
	if service.Bot().IsGroupOwnerOrAdmin(ctx) {
		// owner or admin 不需要撤回
		return
	}
	// 获取当前 group keyword 策略
	groupId := service.Bot().GetGroupId(ctx)
	process := service.Group().GetKeywordProcess(ctx, groupId)
	// 预处理
	if len(process) < 1 {
		// 没有关键词检查策略，跳过撤回功能
		return
	}
	// 获取聊天信息
	msg := service.Bot().GetMessage(ctx)
	shouldRevoke := false
	// 处理
	if _, ok := process[consts.BlacklistCmd]; ok {
		shouldRevoke = isInKeywordBlacklist(ctx, groupId, msg)
	}
	if _, ok := process[consts.WhitelistCmd]; ok && !shouldRevoke {
		shouldRevoke = isNotInKeywordWhitelist(ctx, groupId, msg)
	}
	// 结果处理
	if !shouldRevoke {
		// 不需要撤回
		return
	}
	// 撤回
	service.Bot().RevokeMessage(ctx, service.Bot().GetMsgId(ctx))
	// 打印撤回日志
	g.Log().Infof(ctx, "revoke group(%v) user(%v) because %v",
		groupId,
		service.Bot().GetUserId(ctx),
		msg)
	// 禁言
	doMute(ctx)
	catch = true
	return
}

func isInKeywordBlacklist(ctx context.Context, groupId int64, msg string) (yes bool) {
	blacklists := service.Group().GetKeywordBlacklists(ctx, groupId)
	for k := range blacklists {
		blacklist := service.List().GetListData(ctx, k)
		for kk := range blacklist {
			if strings.Contains(msg, kk) {
				yes = true
				return
			}
		}
	}
	return
}

func isNotInKeywordWhitelist(ctx context.Context, groupId int64, msg string) (yes bool) {
	// 默认不在白名单内
	yes = true
	whitelists := service.Group().GetKeywordWhitelists(ctx, groupId)
	for k := range whitelists {
		whitelist := service.List().GetListData(ctx, k)
		for kk := range whitelist {
			if strings.Contains(msg, kk) {
				yes = false
				return
			}
		}
	}
	return
}

func doMute(ctx context.Context) {
	userId := service.Bot().GetUserId(ctx)
	// 缓存键名
	cacheKey := "RevokeTimes.QQ=" + gconv.String(userId)
	// 过期时间
	expirationDuration := 16 * time.Hour
	timesVar, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	if timesVar == nil {
		// 第一次撤回不禁言
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
	service.Bot().Mute(ctx, muteMinutes*60)
}

func (s *sModule) TryKeywordReply(ctx context.Context) (catch bool) {
	// 获取当前 group reply list
	groupId := service.Bot().GetGroupId(ctx)
	listName := service.Group().GetKeywordReplyList(ctx, groupId)
	if listName == "" {
		// 没有设置回复列表，跳过回复功能
		return
	}
	// 获取 list
	list := service.List().GetListData(ctx, listName)
	// 获取聊天信息
	msg := service.Bot().GetMessage(ctx)
	// 匹配关键词
	for k, v := range list {
		if vv, ok := v.(string); ok && strings.Contains(msg, k) {
			// 匹配成功，回复
			pre := "[CQ:at,qq=" + gconv.String(service.Bot().GetUserId(ctx)) + "]" + vv
			service.Bot().SendMsg(ctx, pre)
			return
		}
	}
	catch = true
	return
}
