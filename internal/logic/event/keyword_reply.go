package event

import (
	"context"
	"net/http"
	"net/url"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"qq-bot-backend/utility/codec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel/attribute"
)

const (
	messagePlaceholder  = "{message}"
	remainPlaceholder   = "{remain}"
	nicknamePlaceholder = "{nickname}"
	userIdPlaceholder   = "{userId}"
	groupIdPlaceholder  = "{groupId}"
)

var (
	cqAtPrefixRe    = regexp.MustCompile(`^\[CQ:at,qq=([^,\]]+)(?:,[^,=]+=[^,\]]*)*]\s*`)
	webhookPrefixRe = regexp.MustCompile(`^webhook(?::([A-Za-z]{3,7}))?(?:#([\s\S]+)#)?(?:<([\s\S]+)>)?(?:@(.+)@)?://(.+)$`)
	webhookHolderRe = regexp.MustCompile(`\{\{webhook(?::([A-Za-z]{3,7}))?(?:#([\s\S]+)#)?(?:<([\s\S]+)>)?(?:@(.+)@)?://(.+)}}`)
	commandPrefixRe = regexp.MustCompile(`^(?:command|cmd)://([\s\S]+)$`)
	rewritePrefixRe = regexp.MustCompile(`^rewrite://([\s\S]+)$`)
	placeholderRe   = regexp.MustCompile(`\{([^}\d\s]+)(\d+)?}`)
)

func decreasePlaceholderIndex(text string) string {
	arr := placeholderRe.FindAllStringSubmatch(text, -1)
	oldNew := make([]string, 0, len(arr)*2)
	for _, sub := range arr {
		if len(sub) < 3 {
			continue
		}
		num := gconv.Int(sub[2]) - 1
		if num <= 0 {
			oldNew = append(oldNew, sub[0], "{"+sub[1]+"}")
			continue
		}
		oldNew = append(oldNew, sub[0], "{"+sub[1]+gconv.String(num)+"}")
	}
	return strings.NewReplacer(oldNew...).Replace(text)
}

func (s *sEvent) TryKeywordReply(ctx context.Context) (caught bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryKeywordReply")
	defer span.End()

	// 获取基础信息
	msg := service.Bot().GetMessage(ctx)
	userId := service.Bot().GetUserId(ctx)
	// 匹配 @bot
	if cqAtPrefixRe.MatchString(msg) {
		sub := cqAtPrefixRe.FindStringSubmatch(msg)
		if sub[1] == gconv.String(service.Bot().GetSelfId(ctx)) {
			msg = strings.Replace(msg, sub[0], "", 1)
		}
	}
	// 匹配关键词
	found, hit, value := service.Util().FindBestKeywordMatch(ctx, msg, service.Namespace().GetGlobalNamespaceLists(ctx))
	if !found || value == "" {
		return
	}
	// 匹配成功，回复
	replyMsg := value
	noReplyPrefix := false
	switch {
	case webhookPrefixRe.MatchString(value):
		replyMsg, noReplyPrefix = s.keywordReplyWebhook(ctx,
			userId, 0, service.Bot().GetNickname(ctx),
			msg, hit, value,
		)
	case rewritePrefixRe.MatchString(value):
		caught = s.keywordReplyRewrite(ctx, s.TryKeywordReply, msg, hit, value)
		replyMsg = ""
	case commandPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyCommand(ctx, msg, hit, value)
	}
	// 内容为空，不回复
	if replyMsg == "" {
		return
	}
	// 限速
	const kind = "replyU"
	key := gconv.String(service.Bot().GetSelfId(ctx)) + "_" + gconv.String(userId)
	if limited, _ := utility.AutoLimit(ctx, kind, key, consts.MaxSendMessageCount, time.Minute); limited {
		g.Log().Notice(ctx, kind, key, "is limited")
		return
	}
	if !noReplyPrefix {
		if msgId := service.Bot().GetMsgId(ctx); msgId != 0 {
			replyMsg = "[CQ:reply,id=" + gconv.String(msgId) + "]" + replyMsg
		}
	}
	service.Bot().SendMsg(ctx, replyMsg)

	caught = true
	return
}

func (s *sEvent) keywordReplyWebhook(ctx context.Context,
	userId, groupId int64,
	nickname, message, hit, value string,
) (replyMsg string, noReplyPrefix bool) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.keywordReplyWebhook")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("keyword_reply_webhook.user_id", userId),
		attribute.Int64("keyword_reply_webhook.group_id", groupId),
		attribute.String("keyword_reply_webhook.nickname", nickname),
		attribute.String("keyword_reply_webhook.message", message),
		attribute.String("keyword_reply_webhook.hit", hit),
		attribute.String("keyword_reply_webhook.value", value),
	)

	// URL
	subMatch := webhookPrefixRe.FindStringSubmatch(codec.DecodeCQCode(value))
	method := strings.ToUpper(subMatch[1])
	if method == "" {
		method = http.MethodGet
	}
	headers := subMatch[2]
	payload := subMatch[3]
	bodyPath := strings.Split(subMatch[4], ".")
	urlLink := subMatch[5]
	// Arguments
	var err error
	message = codec.DecodeCQCode(message)
	hit = codec.DecodeCQCode(hit)
	remain := strings.Replace(message, hit, "", 1)
	// Headers
	if headers != "" {
		headers = strings.NewReplacer(
			"\\n", "\n",
			"\r", "\n",
			messagePlaceholder, message,
			remainPlaceholder, remain,
			nicknamePlaceholder, nickname,
			userIdPlaceholder, gconv.String(userId),
			groupIdPlaceholder, gconv.String(groupId),
		).Replace(headers)
	}
	// URL escape
	urlLink = strings.NewReplacer(
		messagePlaceholder, url.QueryEscape(message),
		remainPlaceholder, url.QueryEscape(remain),
		nicknamePlaceholder, url.QueryEscape(nickname),
		userIdPlaceholder, gconv.String(userId),
		groupIdPlaceholder, gconv.String(groupId),
	).Replace(urlLink)
	// Call webhook
	var body []byte
	var statusCode int
	var contentType string
	switch method {
	case http.MethodGet:
		statusCode, contentType, body, err = utility.SendWebhookRequest(ctx, headers, method, urlLink)
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		if webhookPrefixRe.MatchString(payload) {
			payload, _ = s.keywordReplyWebhook(ctx,
				userId, groupId, nickname,
				codec.EncodeCQCode(message), codec.EncodeCQCode(hit), codec.EncodeCQCode(payload),
			)
		}
		if webhookHolderRe.MatchString(payload) {
			type webhookResult struct {
				placeholder string
				result      string
			}

			whs := webhookHolderRe.FindAllString(payload, -1)
			seen := make(map[string]struct{})
			results := make(chan webhookResult, len(whs))

			go func() {
				wg := sync.WaitGroup{}
				for _, wh := range whs {
					if _, saw := seen[wh]; saw {
						continue
					}
					seen[wh] = struct{}{}

					wg.Go(func() {
						whUri := wh[len("{{") : len(wh)-len("}}")]
						ret, _ := s.keywordReplyWebhook(ctx,
							userId, groupId, nickname,
							codec.EncodeCQCode(message), codec.EncodeCQCode(hit), codec.EncodeCQCode(whUri),
						)
						results <- webhookResult{wh, ret}
					})
				}
				wg.Wait()
				close(results)
			}()

			oldNew := make([]string, 0, len(whs)*2)
			for result := range results {
				oldNew = append(oldNew, result.placeholder, result.result)
			}
			payload = strings.NewReplacer(oldNew...).Replace(payload)
		}

		// Payload
		oldNew := make([]string, 0, 5*2)
		switch payload {
		case messagePlaceholder:
			payload = message
		case remainPlaceholder:
			payload = remain
		case nicknamePlaceholder:
			payload = nickname
		default:
			msg, _ := sonic.MarshalString(message)
			r, _ := sonic.MarshalString(remain)
			nick, _ := sonic.MarshalString(nickname)
			oldNew = append(oldNew,
				messagePlaceholder, msg,
				remainPlaceholder, r,
				nicknamePlaceholder, nick,
			)
		}
		// 占位符替换
		oldNew = append(oldNew,
			userIdPlaceholder, gconv.String(userId),
			groupIdPlaceholder, gconv.String(groupId),
		)
		payload = strings.NewReplacer(oldNew...).Replace(payload)

		statusCode, contentType, body, err = utility.SendWebhookRequest(ctx, headers, method, urlLink, payload)
	default:
		return
	}
	if err != nil {
		g.Log().Warning(ctx, "webhook", statusCode, method, urlLink, message, err)
		return
	}
	// Log
	if statusCode != http.StatusOK {
		g.Log().Notice(ctx,
			nickname+"("+gconv.String(userId)+") in group("+gconv.String(service.Bot().GetGroupId(ctx))+
				") call webhook", statusCode, method, urlLink, message)
	}
	// 媒体文件
	{
		var mediumURL string
		// 如果是图片
		if strings.HasPrefix(contentType, "image/") {
			mediumURL, err = service.File().CacheFile(ctx, body, 5*time.Minute)
			if err != nil {
				replyMsg = "Image cache failed"
				return
			}
			replyMsg = "[CQ:image,file=" + codec.EncodeCQCode(mediumURL) + "]"
			return
		}
		// 如果是音频
		if strings.HasPrefix(contentType, "audio/") {
			mediumURL, err = service.File().CacheFile(ctx, body, 5*time.Minute)
			if err != nil {
				replyMsg = "Audio cache failed"
				return
			}
			replyMsg = "[CQ:record,file=" + codec.EncodeCQCode(mediumURL) + "]"
			noReplyPrefix = true
			return
		}
		// 如果是视频
		if strings.HasPrefix(contentType, "video/") {
			mediumURL, err = service.File().CacheFile(ctx, body, 5*time.Minute)
			if err != nil {
				replyMsg = "Video cache failed"
				return
			}
			replyMsg = "[CQ:video,file=" + codec.EncodeCQCode(mediumURL) + "]"
			noReplyPrefix = true
			return
		}
	}
	// 没有 bodyPath，直接返回 body
	if len(bodyPath) == 1 && bodyPath[0] == "" {
		replyMsg = string(body)
		return
	}
	// 默认视为 JSON 数据
	path := make([]any, len(bodyPath))
	for i, v := range bodyPath {
		index, e := strconv.Atoi(v)
		if e == nil {
			path[i] = index
			continue
		}
		path[i] = v
	}
	// 解析 body 获取数据
	node, err := sonic.Get(body, path...)
	if err != nil {
		replyMsg = "Wrong JSON path"
		return
	}
	if node.TypeSafe() != ast.V_STRING {
		if err = node.LoadAll(); err != nil {
			replyMsg = "Wrong JSON format"
			return
		}
		replyMsg, _ = node.Raw()
		return
	}
	replyMsg, _ = node.StrictString()
	return
}

func (s *sEvent) keywordReplyCommand(ctx context.Context, message, hit, text string) (replyMsg string) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.keywordReplyCommand")
	defer span.End()
	span.SetAttributes(
		attribute.String("keyword_reply_command.message", message),
		attribute.String("keyword_reply_command.hit", hit),
		attribute.String("keyword_reply_command.text", text),
	)

	// 解码提取
	subMatch := commandPrefixRe.FindStringSubmatch(codec.DecodeCQCode(text))
	// 占位符替换
	remain := strings.Replace(message, hit, "", 1)
	subMatch[1] = strings.NewReplacer(
		messagePlaceholder, message,
		remainPlaceholder, remain,
		// 为什么是 " &&"？因为 " &&" 后可能是换行符，需要替换为 " "
		" &&\r", " && ",
		" &&\n", " && ",
	).Replace(subMatch[1])
	// 转换占位符
	subMatch[1] = decreasePlaceholderIndex(subMatch[1])
	// 切分命令
	commands := strings.Split(subMatch[1], " && ")
	var replyBuilder strings.Builder
	for _, command := range commands {
		caught, tmp := service.Command().TryCommand(ctx, strings.TrimSpace(command))
		if !caught {
			return
		}
		if tmp != "" {
			replyBuilder.WriteString(tmp + "\n")
		}
	}
	replyMsg = strings.TrimSuffix(replyBuilder.String(), "\n")
	return
}

func (s *sEvent) keywordReplyRewrite(ctx context.Context,
	try func(context.Context) bool,
	message, hit, text string,
) (caught bool) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.keywordReplyRewrite")
	defer span.End()
	span.SetAttributes(
		attribute.String("keyword_reply_rewrite.message", message),
		attribute.String("keyword_reply_rewrite.hit", hit),
		attribute.String("keyword_reply_rewrite.text", text),
	)

	// 防止循环递归
	if err := service.Bot().SetHistory(ctx, hit); err != nil {
		// rewrite loop detected
		g.Log().Notice(ctx, "rewrite loop detected:", hit)
		return
	}
	// 解码提取
	subMatch := rewritePrefixRe.FindStringSubmatch(codec.DecodeCQCode(text))
	// 占位符替换
	remain := strings.Replace(message, hit, "", 1)
	subMatch[1] = strings.NewReplacer(
		messagePlaceholder, message,
		remainPlaceholder, remain,
		// 为什么是 " &"？因为 " &" 后可能是换行符，需要替换为 " "
		" &\r", " & ",
		" &\n", " & ",
	).Replace(subMatch[1])
	// 切分
	rewrites := strings.Split(subMatch[1], " & ")
	for _, rewrite := range rewrites {
		service.Bot().RewriteMessage(ctx, strings.TrimSpace(rewrite))
		// callback
		caught = try(ctx)
	}
	return
}
