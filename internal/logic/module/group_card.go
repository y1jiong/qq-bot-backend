package module

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func (s *sModule) TryLockCard(ctx context.Context) (catch bool) {
	// 获取当前 group card lock
	groupId := service.Bot().GetGroupId(ctx)
	locked := service.Group().IsCardLocked(ctx, groupId)
	if !locked {
		// 不需要锁定
		return
	}
	catch = true
	oldCard, newCard := service.Bot().GetCardOldNew(ctx)
	if oldCard == "" {
		// 无旧名片允许修改一次
		return
	}
	if oldCard != newCard {
		// 执行锁定
		service.Bot().SetGroupCard(ctx, groupId, service.Bot().GetUserId(ctx), oldCard)
		// 发送提示
		service.Bot().SendPlainMsg(ctx, "名片已锁定，请联系管理员修改")
	}
	return
}

func (s *sModule) TryAutoSetCard(ctx context.Context) (catch bool) {
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
	if _, ok := listMap[gconv.String(userId)].(string); !ok {
		// 不在 auto_set list 中
		return
	}
	// 执行设置
	service.Bot().SetGroupCard(ctx, groupId, userId, listMap[gconv.String(userId)].(string))
	return
}
