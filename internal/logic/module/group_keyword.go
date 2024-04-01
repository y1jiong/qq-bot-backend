package module

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"strings"
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
	if limited, _ := s.AutoLimit(ctx, kind, gid, 2, gconv.Duration("1m")); limited {
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
		method := strings.ToLower(webhookPrefixRe.FindStringSubmatch(value)[1])
		url := webhookPrefixRe.FindStringSubmatch(value)[2]
		// Log
		g.Log().Info(ctx,
			"user("+gconv.String(userId)+") in group("+gconv.String(service.Bot().GetGroupId(ctx))+") call webhook "+url)
		// Arguments
		var err error
		remain := strings.Replace(msg, hit, "", 1)
		switch method {
		case "get", "":
			url = strings.ReplaceAll(url, "{message}", msg)
			url = strings.ReplaceAll(url, "{userId}", gconv.String(userId))
			url = strings.ReplaceAll(url, "{groupId}", gconv.String(groupId))
			url = strings.ReplaceAll(url, "{remain}", remain)
			// Webhook
			replyMsg, err = service.Bot().SendGetWebhook(ctx, url)
		case "post":
			payload := struct {
				GroupId int64  `json:"group_id"`
				UserId  int64  `json:"user_id"`
				Message string `json:"message"`
				Remain  string `json:"remain"`
			}{
				GroupId: groupId,
				UserId:  userId,
				Message: msg,
				Remain:  remain,
			}
			var payloadJson []byte
			payloadJson, err = sonic.ConfigStd.Marshal(payload)
			if err != nil {
				g.Log().Error(ctx, err)
				return
			}
			// Webhook
			replyMsg, err = service.Bot().SendPostWebhook(ctx, url, payloadJson)
		}
		if err != nil {
			g.Log().Notice(ctx, "webhook", url, err)
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
