package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"strings"
	"time"
)

func (s *sEvent) TryLockCard(ctx context.Context) (caught bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryLockCard")
	defer span.End()

	// 获取基础信息
	userId := service.Bot().GetUserId(ctx)
	groupId := service.Bot().GetGroupId(ctx)

	// 获取当前 group card lock
	if !service.Group().IsCardLocked(ctx, groupId) {
		// 不需要锁定
		return
	}
	caught = true

	oldCard, newCard := service.Bot().GetCardOldNew(ctx)
	if strings.TrimSpace(oldCard) == "" || oldCard == newCard {
		// 无旧名片允许修改一次
		return
	}

	// 防止重复修改群名片
	cacheKey := "LockCard_" + gconv.String(groupId) + "_" + gconv.String(userId)
	cardVar, err := gcache.GetOrSet(ctx, cacheKey, oldCard, time.Hour)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	lastCard := cardVar.String()
	if newCard == lastCard {
		// 名片未改变
		if _, err = gcache.Remove(ctx, cacheKey); err != nil {
			g.Log().Warning(ctx, err)
			return
		}
		return
	}
	oldCard = lastCard

	// 执行锁定
	if err = service.Bot().SetGroupCard(ctx, groupId, service.Bot().GetUserId(ctx), oldCard); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// 发送提示
	service.Bot().SendMsg(ctx, "[CQ:at,qq="+gconv.String(userId)+"]名片已锁定，请联系管理员修改")
	return
}

func (s *sEvent) TryAutoSetCard(ctx context.Context) (caught bool) {
	// 获取当前 group card auto_set list
	groupId := service.Bot().GetGroupId(ctx)
	listName := service.Group().GetCardAutoSetList(ctx, groupId)
	// 预处理
	if listName == "" {
		// 没有设置 auto_set list
		return
	}
	// 处理
	caught = true
	listMap := service.List().GetListData(ctx, listName)
	userId := service.Bot().GetUserId(ctx)
	if card, ok := listMap[gconv.String(userId)].(string); ok {
		// 执行设置
		if err := service.Bot().SetGroupCard(ctx, groupId, userId, card); err != nil {
			g.Log().Warning(ctx, err)
		}
	}
	return
}
