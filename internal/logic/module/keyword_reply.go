package module

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"net/url"
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"
	"time"
)

var (
	webhookPrefixRe = regexp.MustCompile(`^webhook(?::([A-Za-z]{3,7}))?(?:<(.+)>)?(?:@(.+)@)?://(.+)$`)
	commandPrefixRe = regexp.MustCompile(`^(?:command|cmd)://(.+)$`)
)

func (s *sModule) TryKeywordReply(ctx context.Context) (catch bool) {
	// 获取基础信息
	msg := service.Bot().GetMessage(ctx)
	userId := service.Bot().GetUserId(ctx)
	// 匹配关键词
	contains, hit, value := s.isOnKeywordLists(ctx, msg, service.Namespace().GetPublicNamespaceLists(ctx))
	if !contains || value == "" {
		return
	}
	// 限速
	kind := "replyU"
	uid := gconv.String(userId)
	if limited, _ := s.AutoLimit(ctx, kind, uid, 5, time.Minute); limited {
		g.Log().Info(ctx, kind, uid, "is limited")
		return
	}
	// 匹配成功，回复
	replyMsg := value
	switch {
	case webhookPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyWebhook(ctx, userId, 0, msg, hit, value)
	case commandPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyCommand(ctx, msg, hit, value)
	}
	// 内容为空，不回复
	if replyMsg == "" {
		return
	}
	pre := "[CQ:reply,id=" + gconv.String(service.Bot().GetMsgId(ctx)) + "]" + replyMsg
	service.Bot().SendMsg(ctx, pre)
	catch = true
	return
}

func (s *sModule) keywordReplyWebhook(ctx context.Context, userId, groupId int64,
	message, hit, value string) (replyMsg string) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}
	// Url
	subMatch := webhookPrefixRe.FindStringSubmatch(service.Codec().DecodeCqCode(value))
	method := strings.ToUpper(subMatch[1])
	if method == "" {
		method = http.MethodGet
	}
	payload := subMatch[2]
	bodyPath := strings.Split(subMatch[3], ".")
	urlLink := subMatch[4]
	// Arguments
	var err error
	message = service.Codec().DecodeCqCode(message)
	hit = service.Codec().DecodeCqCode(hit)
	remain := strings.Replace(message, hit, "", 1)
	urlLink = strings.ReplaceAll(urlLink, "{message}", url.QueryEscape(message))
	urlLink = strings.ReplaceAll(urlLink, "{userId}", gconv.String(userId))
	urlLink = strings.ReplaceAll(urlLink, "{groupId}", gconv.String(groupId))
	urlLink = strings.ReplaceAll(urlLink, "{remain}", url.QueryEscape(remain))
	// Log
	g.Log().Info(ctx,
		"user("+gconv.String(userId)+") in group("+gconv.String(service.Bot().GetGroupId(ctx))+
			") call webhook", method, urlLink)
	// Log end
	var body string
	// Webhook
	switch method {
	case http.MethodGet:
		_, body, err = s.WebhookGetHeadConnectOptionsTrace(ctx, method, urlLink)
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		payload = strings.ReplaceAll(payload, "{message}", message)
		payload = strings.ReplaceAll(payload, "{userId}", gconv.String(userId))
		payload = strings.ReplaceAll(payload, "{groupId}", gconv.String(groupId))
		payload = strings.ReplaceAll(payload, "{remain}", remain)
		_, body, err = s.WebhookPostPutPatchDelete(ctx, method, urlLink, payload)
	default:
		return
	}
	if err != nil {
		g.Log().Notice(ctx, "webhook", method, urlLink, err)
		return
	}
	// 没有 bodyPath，直接返回 body
	if len(bodyPath) == 1 && bodyPath[0] == "" {
		replyMsg = body
		return
	}
	// 解析 body 获取数据
	path := make([]any, len(bodyPath))
	for i, v := range bodyPath {
		path[i] = v
	}
	node, err := sonic.Get([]byte(body), path...)
	if err != nil {
		replyMsg = "Wrong json path"
		return
	}
	if node.Type() != ast.V_STRING {
		tmp, _ := node.MarshalJSON()
		replyMsg = string(tmp)
		return
	}
	replyMsg, _ = node.StrictString()
	return
}

func (s *sModule) keywordReplyCommand(ctx context.Context, message, hit, text string) (replyMsg string) {
	// 必须全字匹配
	if message != hit {
		return
	}
	// 解码
	subMatch := commandPrefixRe.FindStringSubmatch(service.Codec().DecodeCqCode(text))
	// 切分命令
	commands := strings.Split(subMatch[1], " && ")
	var replyBuilder strings.Builder
	for _, command := range commands {
		catch, tmp := service.Command().TryCommand(ctx, command)
		if !catch {
			return
		}
		replyBuilder.WriteString(tmp + "\n")
	}
	replyMsg = strings.TrimRight(replyBuilder.String(), "\n")
	return
}
