package event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
)

func (s *sEvent) TryApproveAddGroup(ctx context.Context) (catch bool) {
	// 获取当前 group approval 策略
	groupId := service.Bot().GetGroupId(ctx)
	policy := service.Group().GetApprovalPolicy(ctx, groupId)
	// 预处理
	if len(policy) == 0 {
		// 没有入群审核策略，跳过审核功能
		return
	}
	// 默认通过审核
	pass := true
	// 局部变量
	comment := service.Bot().GetComment(ctx)
	userId := service.Bot().GetUserId(ctx)
	var extra, blackReason string
	isOnBlacklist := false
	// 处理
	if _, ok := policy[consts.McCmd]; ok {
		// mc 正版验证
		pass, extra = verifyMinecraftGenuine(ctx, comment)
	}
	if _, ok := policy[consts.RegexpCmd]; ok && pass {
		// 正则表达式
		pass, extra = isMatchRegexp(ctx, groupId, comment)
	}
	if _, ok := policy[consts.WhitelistCmd]; ok && pass {
		// 白名单
		pass = isOnApprovalWhitelist(ctx, groupId, userId, extra)
	}
	if _, ok := policy[consts.BlacklistCmd]; ok && pass {
		// 黑名单
		pass, blackReason = isNotOnApprovalBlacklist(ctx, groupId, userId)
		isOnBlacklist = !pass
	}
	// 回执与日志
	var logMsg string
	if !service.Group().IsApprovalNotifyOnlyEnabled(ctx, groupId) &&
		((!pass && service.Group().IsApprovalAutoRejectEnabled(ctx, groupId)) ||
			(pass && service.Group().IsApprovalAutoPassEnabled(ctx, groupId)) ||
			isOnBlacklist) {
		if isOnBlacklist {
			// 黑名单拒绝
			pass = false
		}
		// 在开启自动通过、自动拒绝和黑名单的条件下发送审核回执
		// 审核请求回执
		service.Bot().ApproveJoinGroup(ctx,
			service.Bot().GetFlag(ctx),
			service.Bot().GetSubType(ctx),
			pass,
			"")
		// 打印审核日志
		if pass {
			logMsg = fmt.Sprintf("approve user(%v) join group(%v) with %v",
				userId,
				groupId,
				comment)
		} else {
			logMsg = fmt.Sprintf("REJECT user(%v) join group(%v) with %v",
				userId,
				groupId,
				comment)
		}
	} else if pass {
		// 打印跳过同意日志
		logMsg = fmt.Sprintf("skip processing approve user(%v) join group(%v) with %v",
			userId,
			groupId,
			comment)
	} else if !pass {
		// 打印跳过拒绝日志
		logMsg = fmt.Sprintf("skip processing REJECT user(%v) join group(%v) with %v",
			userId,
			groupId,
			comment)
	}
	if isOnBlacklist {
		logMsg = "[hit blacklist]" + blackReason + "\n" + logMsg
	}
	g.Log().Info(ctx, logMsg)
	// 通知
	notificationGroupId := service.Group().GetApprovalNotificationGroupId(ctx, groupId)
	if notificationGroupId != 0 {
		_, _ = service.Bot().SendMessage(ctx,
			"group", 0, notificationGroupId, logMsg, true)
	}
	catch = true
	return
}

func isMatchRegexp(ctx context.Context, groupId int64, comment string) (match bool, matched string) {
	exp := service.Group().GetApprovalRegexp(ctx, groupId)
	// 匹配正则
	re, err := regexp.Compile(exp)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	ans := re.FindStringSubmatch(comment)
	switch len(ans) {
	case 0:
	case 1:
		matched = ans[0]
		match = true
	default:
		// 读取第一个子表达式
		matched = ans[1]
		match = true
	}
	return
}

func verifyMinecraftGenuine(ctx context.Context, comment string) (genuine bool, uuid string) {
	// Minecraft 正版验证
	genuine, _, uuid, err := service.ThirdParty().QueryMinecraftGenuineUser(ctx, comment)
	if err != nil {
		g.Log().Notice(ctx, err)
	}
	return
}

func isOnApprovalWhitelist(ctx context.Context, groupId, userId int64, extra string) bool {
	// 获取白名单组
	whitelists := service.Group().GetApprovalWhitelists(ctx, groupId)
	for k := range whitelists {
		// 获取其中一个白名单
		whitelist := service.List().GetListData(ctx, k)
		if v, ok := whitelist[gconv.String(userId)]; ok {
			// userId 在白名单中
			if vv, okay := v.(string); okay {
				// 有额外验证信息
				if vv == extra {
					return true
				}
			} else {
				// 没有额外验证信息
				return true
			}
		}
		if extra == "" {
			// 没有额外验证信息则跳过反向验证
			continue
		}
		// 反向验证
		if v, ok := whitelist[extra]; ok {
			if vv, okay := v.(string); okay {
				if vv == gconv.String(userId) {
					return true
				}
			}
		}
	}
	return false
}

func isNotOnApprovalBlacklist(ctx context.Context, groupId, userId int64) (bool, string) {
	// 默认不在黑名单内
	// 获取黑名单组
	blacklists := service.Group().GetApprovalBlacklists(ctx, groupId)
	for k := range blacklists {
		// 获取其中一个黑名单
		blacklist := service.List().GetListData(ctx, k)
		if v, ok := blacklist[gconv.String(userId)]; ok {
			// userId 在黑名单中
			if vv, okay := v.(string); okay {
				// 有黑名单原因
				return false, vv
			} else {
				// 没有黑名单原因
				return false, ""
			}
		}
	}
	return true, ""
}
