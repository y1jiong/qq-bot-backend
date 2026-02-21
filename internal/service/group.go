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
		GetApprovalPolicy(ctx context.Context, groupId int64) (policy map[string]any)
		GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetApprovalRegexp(ctx context.Context, groupId int64) (exp string)
		GetApprovalNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64)
		GetApprovalLevel(ctx context.Context, groupId int64) (level int64)
		GetApprovalReason(ctx context.Context, groupId int64) (reason string)
		IsApprovalNotifyOnlyEnabled(ctx context.Context, groupId int64) bool
		IsApprovalAutoPassEnabled(ctx context.Context, groupId int64) bool
		IsApprovalAutoRejectEnabled(ctx context.Context, groupId int64) bool
		AddApprovalPolicyReturnRes(ctx context.Context, groupId int64, policyName string, args ...string) (retMsg string)
		RemoveApprovalPolicyReturnRes(ctx context.Context, groupId int64, policyName string, args ...string) (retMsg string)
		AcceptBroadcastReturnRes(ctx context.Context, groupId int64) (retMsg string)
		RejectBroadcastReturnRes(ctx context.Context, groupId int64) (retMsg string)
		GetCardAutoSetList(ctx context.Context, groupId int64) (listName string)
		IsCardLocked(ctx context.Context, groupId int64) bool
		SetAutoSetListReturnRes(ctx context.Context, groupId int64, listName string) (retMsg string)
		RemoveAutoSetListReturnRes(ctx context.Context, groupId int64) (retMsg string)
		CheckCardWithRegexpReturnRes(ctx context.Context, groupId int64, listName string, exp string) (retMsg string)
		CheckCardByListReturnRes(ctx context.Context, groupId int64, toList string, fromList string) (retMsg string)
		LockCardReturnRes(ctx context.Context, groupId int64) (retMsg string)
		UnlockCardReturnRes(ctx context.Context, groupId int64) (retMsg string)
		ExportGroupMemberListReturnRes(ctx context.Context, groupId int64, listName string) (retMsg string)
		IsBinding(ctx context.Context, groupId int64) bool
		GetNamespace(ctx context.Context, groupId int64) string
		GetKeywordPolicy(ctx context.Context, groupId int64) (policy map[string]any)
		GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetKeywordReplyLists(ctx context.Context, groupId int64) (replyLists map[string]any)
		AddKeywordPolicyReturnRes(ctx context.Context, groupId int64, policyName string, args ...string) (retMsg string)
		RemoveKeywordPolicyReturnRes(ctx context.Context, groupId int64, policyName string, args ...string) (retMsg string)
		GetLogLeaveList(ctx context.Context, groupId int64) (listName string)
		GetLogApprovalList(ctx context.Context, groupId int64) (listName string)
		SetLogLeaveListReturnRes(ctx context.Context, groupId int64, listName string) (retMsg string)
		RemoveLogLeaveListReturnRes(ctx context.Context, groupId int64) (retMsg string)
		SetLogApprovalListReturnRes(ctx context.Context, groupId int64, listName string) (retMsg string)
		RemoveLogApprovalListReturnRes(ctx context.Context, groupId int64) (retMsg string)
		IsAntiRecallEnabled(ctx context.Context, groupId int64) bool
		GetMessageNotificationGroupId(ctx context.Context, groupId int64) (notificationGroupId int64)
		IsOnlyAntiRecallMemberSet(ctx context.Context, groupId int64) bool
		SetAntiRecallReturnRes(ctx context.Context, groupId int64, enable bool) (retMsg string)
		SetMessageNotificationReturnRes(ctx context.Context, groupId int64, notificationGroupId int64) (retMsg string)
		RemoveMessageNotificationReturnRes(ctx context.Context, groupId int64) (retMsg string)
		SetOnlyAntiRecallMemberReturnRes(ctx context.Context, groupId int64, enable bool) (retMsg string)
		BindNamespaceReturnRes(ctx context.Context, groupId int64, namespace string) (retMsg string)
		UnbindReturnRes(ctx context.Context, groupId int64) (retMsg string)
		CloneReturnRes(ctx context.Context, groupId int64, srcGroupId int64) (retMsg string)
		QueryGroupReturnRes(ctx context.Context, groupId int64) (retMsg string)
		KickFromListReturnRes(ctx context.Context, groupId int64, listName string) (retMsg string)
		KeepFromListReturnRes(ctx context.Context, groupId int64, listName string) (retMsg string)
		CheckExistReturnRes(ctx context.Context) (retMsg string)
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
