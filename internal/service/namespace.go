// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	INamespace interface {
		GetForwardingToAliasList(ctx context.Context) (aliasList map[string]any)
		GetForwardingTo(ctx context.Context, alias string) (url, authorization string)
		IsForwardingMatchUserId(ctx context.Context, userId string) bool
		IsForwardingMatchGroupId(ctx context.Context, groupId string) bool
		AddForwardingToReturnRes(ctx context.Context, alias, url, authorization string) (retMsg string)
		RemoveForwardingToReturnRes(ctx context.Context, alias string) (retMsg string)
		AddForwardingMatchUserIdReturnRes(ctx context.Context, userId string) (retMsg string)
		AddForwardingMatchGroupIdReturnRes(ctx context.Context, groupId string) (retMsg string)
		RemoveForwardingMatchUserIdReturnRes(ctx context.Context, userId string) (retMsg string)
		RemoveForwardingMatchGroupIdReturnRes(ctx context.Context, groupId string) (retMsg string)
		ResetForwardingMatchUserIdReturnRes(ctx context.Context) (retMsg string)
		ResetForwardingMatchGroupIdReturnRes(ctx context.Context) (retMsg string)
		AddNamespaceList(ctx context.Context, namespace, listName string)
		RemoveNamespaceList(ctx context.Context, namespace, listName string)
		GetNamespaceLists(ctx context.Context, namespace string) (lists map[string]any)
		GetNamespaceListsIncludingGlobal(ctx context.Context, namespace string) (lists map[string]any)
		GetGlobalNamespaceLists(ctx context.Context) (lists map[string]any)
		IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) bool
		IsNamespaceOwnerOrAdminOrOperator(ctx context.Context, namespace string, userId int64) bool
		IsGlobalNamespace(namespace string) bool
		GetGlobalNamespace() string
		IsNamespacePropertyPublic(ctx context.Context, namespace string) bool
		AddNewNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		RemoveNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		QueryNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		QueryOwnNamespaceReturnRes(ctx context.Context) (retMsg string)
		AddNamespaceAdminReturnRes(ctx context.Context, namespace string, userId int64) (retMsg string)
		RemoveNamespaceAdminReturnRes(ctx context.Context, namespace string, userId int64) (retMsg string)
		ResetNamespaceAdminReturnRes(ctx context.Context, namespace string) (retMsg string)
		ChangeNamespaceOwnerReturnRes(ctx context.Context, namespace, ownerId string) (retMsg string)
		SetNamespacePropertyPublicReturnRes(ctx context.Context, namespace string, value bool) (retMsg string)
	}
)

var (
	localNamespace INamespace
)

func Namespace() INamespace {
	if localNamespace == nil {
		panic("implement not found for interface INamespace, forgot register?")
	}
	return localNamespace
}

func RegisterNamespace(i INamespace) {
	localNamespace = i
}
