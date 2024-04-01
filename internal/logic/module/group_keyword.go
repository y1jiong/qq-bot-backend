package module

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"net/url"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"strings"
	"time"
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
		shouldRecall, hit = isOnKeywordLists(ctx, msg, service.Group().GetKeywordBlacklists(ctx, groupId))
	}
	if _, ok := process[consts.WhitelistCmd]; ok && shouldRecall {
		in, _ := isOnKeywordLists(ctx, msg, service.Group().GetKeywordWhitelists(ctx, groupId))
		shouldRecall = !in
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
	logMsg := fmt.Sprintf("recall group(%v) user(%v) hit(%v) detail %v",
		groupId,
		userId,
		hit,
		msg)
	g.Log().Info(ctx, logMsg)
	// 通知
	notificationGroupId := service.Group().GetMessageNotificationGroupId(ctx, groupId)
	if notificationGroupId > 0 {
		service.Bot().SendMessage(ctx,
			"group", 0, notificationGroupId, logMsg, true)
	}
	// 禁言
	s.AutoMute(ctx, "keyword", groupId, userId,
		1, 5, 0, gconv.Duration("16h"))
	catch = true
	return
}

func isOnKeywordLists(ctx context.Context, msg string, lists map[string]any) (in bool, hit string) {
	for k := range lists {
		blacklist := service.List().GetListData(ctx, k)
		if contains, hitStr, _ := service.Module().MultiContains(msg, blacklist); contains {
			in = true
			hit = hitStr
			return
		}
	}
	return
}

func (s *sModule) TryKeywordReply(ctx context.Context) (catch bool) {
	userId := service.Bot().GetUserId(ctx)
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
	contains, hit, value := service.Module().MultiContains(msg, listMap)
	if !contains || value == "" {
		return
	}
	// 限速
	kind := "replyG"
	gid := gconv.String(groupId)
	if limited, _ := s.AutoLimit(ctx, kind, gid, 5, time.Minute); limited {
		g.Log().Info(ctx, kind, gid, "is limited")
		return
	}
	// 匹配成功，回复
	replyMsg := value
	if webhookPrefixRe.MatchString(value) {
		// 必须以 hit 开头
		if !strings.HasPrefix(msg, hit) {
			return
		}
		// Url
		subMatch := webhookPrefixRe.FindStringSubmatch(service.Codec().DecodeCqCode(value))
		method := strings.ToUpper(subMatch[1])
		if method == "" {
			method = http.MethodGet
		}
		payload := subMatch[2]
		urlLink := subMatch[3]
		// Arguments
		var err error
		msg = service.Codec().DecodeCqCode(msg)
		hit = service.Codec().DecodeCqCode(hit)
		remain := strings.Replace(msg, hit, "", 1)
		urlLink = strings.ReplaceAll(urlLink, "{message}", url.QueryEscape(msg))
		urlLink = strings.ReplaceAll(urlLink, "{userId}", gconv.String(userId))
		urlLink = strings.ReplaceAll(urlLink, "{groupId}", gconv.String(groupId))
		urlLink = strings.ReplaceAll(urlLink, "{remain}", url.QueryEscape(remain))
		// Log
		g.Log().Info(ctx,
			"user("+gconv.String(userId)+") in group("+gconv.String(service.Bot().GetGroupId(ctx))+
				") call webhook", method, urlLink)
		// Log end
		switch method {
		case http.MethodGet:
			// Webhook
			replyMsg, err = service.Bot().SendGetWebhook(ctx, urlLink)
		case http.MethodPost:
			payload = strings.ReplaceAll(payload, "{message}", url.QueryEscape(msg))
			payload = strings.ReplaceAll(payload, "{userId}", gconv.String(userId))
			payload = strings.ReplaceAll(payload, "{groupId}", gconv.String(groupId))
			payload = strings.ReplaceAll(payload, "{remain}", url.QueryEscape(remain))
			// Webhook
			replyMsg, err = service.Bot().SendPostWebhook(ctx, urlLink, payload)
		}
		if err != nil {
			g.Log().Notice(ctx, "webhook", method, urlLink, err)
			return
		}
		// 内容为空，不回复
		if replyMsg == "" {
			return
		}
	}
	pre := "[CQ:reply,id=" + gconv.String(service.Bot().GetMsgId(ctx)) + "]" + replyMsg
	service.Bot().SendMsg(ctx, pre)
	catch = true
	return
}
