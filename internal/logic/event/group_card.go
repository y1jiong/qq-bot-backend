package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"time"
)

func (s *sEvent) TryLockCard(ctx context.Context) (catch bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryLockCard")
	defer span.End()

	// 获取基础信息
	userId := service.Bot().GetUserId(ctx)
	groupId := service.Bot().GetGroupId(ctx)
	// 获取当前 group card lock
	locked := service.Group().IsCardLocked(ctx, groupId)
	if !locked {
		// 不需要锁定
		return
	}

	catch = true

	oldCard, newCard := service.Bot().GetCardOldNew(ctx)
	if oldCard == "" || oldCard == newCard {
		// 无旧名片允许修改一次
		return
	}
	// 防止重复修改群名片
	cacheKey := "LockCard_" + gconv.String(groupId) + "_" + gconv.String(userId)
	cardVar, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	if cardVar != nil {
		cardVarStr := cardVar.String()
		if newCard == cardVarStr {
			// 名片未改变
			_, err = gcache.Remove(ctx, cacheKey)
			if err != nil {
				g.Log().Warning(ctx, err)
				return
			}
			return
		}
		oldCard = cardVarStr
	} else {
		// 设置缓存
		if err = gcache.Set(ctx, cacheKey, oldCard, time.Hour); err != nil {
			g.Log().Warning(ctx, err)
			return
		}
	}
	// 执行锁定
	service.Bot().SetGroupCard(ctx, groupId, service.Bot().GetUserId(ctx), oldCard)
	// 发送提示
	service.Bot().SendMsg(ctx, "[CQ:at,qq="+gconv.String(userId)+"]名片已锁定，请联系管理员修改")
	return
}

func (s *sEvent) TryAutoSetCard(ctx context.Context) (catch bool) {
	// 获取当前 group card auto_set list
	groupId := service.Bot().GetGroupId(ctx)
	listName := service.Group().GetCardAutoSetList(ctx, groupId)
	// 预处理
	if listName == "" {
		// 没有设置 auto_set list
		return
	}
	// 处理
	catch = true
	listMap := service.List().GetListData(ctx, listName)
	userId := service.Bot().GetUserId(ctx)
	if card, ok := listMap[gconv.String(userId)].(string); ok {
		// 执行设置
		service.Bot().SetGroupCard(ctx, groupId, userId, card)
	}
	return
}
