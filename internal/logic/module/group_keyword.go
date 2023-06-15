package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
)

func (s *sModule) TryKeywordRecall(ctx context.Context) (catch bool) {
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
	shouldRecall := false
	// 命中规则
	hit := ""
	// 处理
	if _, ok := process[consts.BlacklistCmd]; ok {
		shouldRecall, hit = isInKeywordBlacklist(ctx, groupId, msg)
	}
	if _, ok := process[consts.WhitelistCmd]; ok && !shouldRecall {
		shouldRecall, hit = isNotInKeywordWhitelist(ctx, groupId, msg)
	}
	// 结果处理
	if !shouldRecall {
		// 不需要撤回
		return
	}
	// 撤回
	service.Bot().RecallMessage(ctx, service.Bot().GetMsgId(ctx))
	userId := service.Bot().GetUserId(ctx)
	// 打印撤回日志
	g.Log().Infof(ctx, "recall group(%v) user(%v) hit(%v) detail %v",
		groupId,
		userId,
		hit,
		msg)
	// 禁言
	s.AutoMute(ctx, "keyword", groupId, userId,
		1, 5, 0, gconv.Duration("16h"))
	catch = true
	return
}

func isInKeywordBlacklist(ctx context.Context, groupId int64, msg string) (in bool, hit string) {
	blacklists := service.Group().GetKeywordBlacklists(ctx, groupId)
	for k := range blacklists {
		blacklist := service.List().GetListData(ctx, k)
		if contains, hitStr, _ := service.Module().MultiContains(msg, blacklist); contains {
			in = true
			hit = hitStr
			return
		}
	}
	return
}

func isNotInKeywordWhitelist(ctx context.Context, groupId int64, msg string) (notIn bool, hit string) {
	// 默认不在白名单内
	notIn = true
	whitelists := service.Group().GetKeywordWhitelists(ctx, groupId)
	for k := range whitelists {
		whitelist := service.List().GetListData(ctx, k)
		if contains, hitStr, _ := service.Module().MultiContains(msg, whitelist); contains {
			notIn = false
			hit = hitStr
			return
		}
	}
	return
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
	listMap := service.List().GetListData(ctx, listName)
	// 获取聊天信息
	msg := service.Bot().GetMessage(ctx)
	// 匹配关键词
	if contains, _, value := service.Module().MultiContains(msg, listMap); contains && value != "" {
		// 匹配成功，回复
		pre := "[CQ:at,qq=" + gconv.String(service.Bot().GetUserId(ctx)) + "]" + value
		service.Bot().SendMsg(ctx, pre)
	}
	catch = true
	return
}
