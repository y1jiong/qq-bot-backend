// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IGroup interface {
		GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any)
		GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetApprovalRegexp(ctx context.Context, groupId int64) (exp string)
		GetApprovalNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64)
		IsEnabledApprovalAutoPass(ctx context.Context, groupId int64) (enabled bool)
		IsEnabledApprovalAutoReject(ctx context.Context, groupId int64) (enabled bool)
		AddApprovalProcessReturnRes(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveApprovalProcessReturnRes(ctx context.Context, groupId int64, processName string, args ...string)
		GetCardAutoSetList(ctx context.Context, groupId int64) (listName string)
		IsCardLocked(ctx context.Context, groupId int64) (locked bool)
		SetAutoSetListReturnRes(ctx context.Context, groupId int64, listName string)
		RemoveAutoSetListReturnRes(ctx context.Context, groupId int64)
		CheckCardWithRegexpReturnRes(ctx context.Context, groupId int64, listName, exp string)
		CheckCardByListReturnRes(ctx context.Context, groupId int64, toList, fromList string)
		LockCardReturnRes(ctx context.Context, groupId int64)
		UnlockCardReturnRes(ctx context.Context, groupId int64)
		ExportGroupMemberListReturnRes(ctx context.Context, groupId int64, listName string)
		BindNamespaceReturnRes(ctx context.Context, groupId int64, namespace string)
		UnbindReturnRes(ctx context.Context, groupId int64)
		QueryGroupReturnRes(ctx context.Context, groupId int64)
		KickFromListReturnRes(ctx context.Context, groupId int64, listName string)
		KeepFromListReturnRes(ctx context.Context, groupId int64, listName string)
		CheckExistReturnRes(ctx context.Context)
		GetKeywordProcess(ctx context.Context, groupId int64) (process map[string]any)
		GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetKeywordReplyList(ctx context.Context, groupId int64) (listName string)
		AddKeywordProcessReturnRes(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveKeywordProcessReturnRes(ctx context.Context, groupId int64, processName string, args ...string)
		GetLogLeaveList(ctx context.Context, groupId int64) (listName string)
		SetLogLeaveListReturnRes(ctx context.Context, groupId int64, listName string)
		RemoveLogLeaveListReturnRes(ctx context.Context, groupId int64)
		IsEnabledAntiRecall(ctx context.Context, groupId int64) (enabled bool)
		GetMessageNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64)
		SetAntiRecallReturnRes(ctx context.Context, groupId int64, enable bool)
		SetMessageNotificationReturnRes(ctx context.Context, groupId int64, notificationGroupId int64)
		RemoveMessageNotificationReturnRes(ctx context.Context, groupId int64)
	}
)

var (
	localGroup IGroup
)

func Group() IGroup {
	if localGroup == nil {
		panic("implement not found for interface IGroup, forgot register?")
	}
	return localGroup
}

func RegisterGroup(i IGroup) {
	localGroup = i
}
