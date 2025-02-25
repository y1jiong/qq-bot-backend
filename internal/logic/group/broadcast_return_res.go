package group

import (
	"context"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
)

func (s *sGroup) AcceptBroadcastReturnRes(ctx context.Context, groupId int64) (retMsg string) {
	if groupId == 0 {
		return
	}
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}

	_, err := dao.Group.Ctx(ctx).
		Data(dao.Group.Columns().AcceptBroadcast, true).
		Where(dao.Group.Columns().GroupId, groupId).
		Update()
	if err != nil {
		retMsg = "接受广播失败"
	}

	retMsg = "接受广播"
	return
}

func (s *sGroup) RejectBroadcastReturnRes(ctx context.Context, groupId int64) (retMsg string) {
	if groupId == 0 {
		return
	}
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}

	_, err := dao.Group.Ctx(ctx).
		Data(dao.Group.Columns().AcceptBroadcast, false).
		Where(dao.Group.Columns().GroupId, groupId).
		Update()
	if err != nil {
		retMsg = "拒绝广播失败"
	}

	retMsg = "拒绝广播"
	return
}
