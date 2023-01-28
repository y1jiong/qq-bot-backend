package module

import (
	"context"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"strings"
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
		// 没有关键字检查流程，跳过撤回功能
		return
	}
	// 获取聊天信息
	msg := service.Bot().GetMessage(ctx)
	shouldRevoke := false
	// 处理流程
	if _, ok := process[consts.BlacklistCmd]; ok && !shouldRevoke {
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
	catch = true
	return
}

func isInKeywordBlacklist(ctx context.Context, groupId int64, msg string) (yes bool) {
	blacklists := service.Group().GetKeywordBlacklists(ctx, groupId)
	for k := range blacklists {
		blacklist := service.List().GetList(ctx, k)
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
		whitelist := service.List().GetList(ctx, k)
		for kk := range whitelist {
			if strings.Contains(msg, kk) {
				yes = false
				return
			}
		}
	}
	return
}
