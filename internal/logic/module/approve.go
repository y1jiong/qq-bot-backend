package module

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
)

func (s *sModule) TryApproveAddGroup(ctx context.Context) (catch bool) {
	comment := service.Bot().GetComment(ctx)
	// 获取当前 group approval 策略
	groupId := service.Bot().GetGroupId(ctx)
	process := service.Group().GetApprovalProcess(ctx, groupId)
	// 预处理
	if len(process) < 1 {
		// 没有入群审批策略，跳过审批功能
		return
	}
	// 默认通过审批
	pass := true
	// 局部变量
	userId := service.Bot().GetUserId(ctx)
	var mcUuid string
	// 处理流程
	if _, ok := process[consts.RegexpCmd]; ok && pass {
		// 正则表达式
		pass = isMatchRegexp(ctx, groupId, comment)
	}
	if _, ok := process[consts.McCmd]; ok && pass {
		// mc 正版验证
		pass, mcUuid = verifyMinecraftGenuine(ctx, comment)
	}
	if _, ok := process[consts.WhitelistCmd]; ok && pass {
		// 白名单
		pass = isInWhitelist(ctx, groupId, userId, mcUuid)
	}
	if _, ok := process[consts.BlacklistCmd]; ok && pass {
		// 黑名单
		pass = isNotInBlacklist(ctx, groupId, userId, mcUuid)
	}
	// 审批请求回执
	service.Bot().ApproveAddGroup(ctx,
		service.Bot().GetFlag(ctx),
		service.Bot().GetSubType(ctx),
		pass,
		"auto reject")
	// 打印通过的日志
	if pass {
		g.Log().Infof(ctx, "approve user(%v) join group(%v) with %v",
			userId,
			groupId,
			comment)
	}
	catch = true
	return
}

func isMatchRegexp(ctx context.Context, groupId int64, comment string) (yes bool) {
	re := service.Group().GetRegexp(ctx, groupId)
	// 匹配正则
	yes, err := regexp.MatchString(re, comment)
	if err != nil {
		g.Log().Warning(ctx, err)
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

func isInWhitelist(ctx context.Context, groupId, userId int64, extra string) (yes bool) {
	// 获取白名单组
	whitelists := service.Group().GetWhitelists(ctx, groupId)
	for k := range whitelists {
		// 获取其中一个白名单
		whitelist := service.List().GetList(ctx, k)
		if v, ok := whitelist[gconv.String(userId)]; ok {
			// userId 在白名单中
			if vv, okay := v.(string); okay {
				// 有额外验证信息
				if vv == extra {
					yes = true
					return
				}
			} else {
				// 没有额外验证信息
				yes = true
				return
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
					yes = true
					return
				}
			}
		}
	}
	return
}

func isNotInBlacklist(ctx context.Context, groupId, userId int64, extra string) (yes bool) {
	// 默认不在黑名单内
	yes = true
	// 获取黑名单组
	blacklists := service.Group().GetBlacklists(ctx, groupId)
	for k := range blacklists {
		// 获取其中一个黑名单
		blacklist := service.List().GetList(ctx, k)
		if v, ok := blacklist[gconv.String(userId)]; ok {
			// userId 在黑名单中
			if vv, okay := v.(string); okay {
				// 有额外验证信息
				if vv == extra {
					yes = false
					return
				}
			} else {
				// 没有额外验证信息
				yes = false
				return
			}
		}
		if extra == "" {
			// 没有额外验证信息则跳过反向验证
			continue
		}
		// 反向验证
		if v, ok := blacklist[extra]; ok {
			if vv, okay := v.(string); okay {
				if vv == gconv.String(userId) {
					yes = false
					return
				}
			}
		}
	}
	return
}
