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
		AddNamespaceList(ctx context.Context, namespace, listName string)
		RemoveNamespaceList(ctx context.Context, namespace, listName string)
		GetNamespaceLists(ctx context.Context, namespace string) (lists map[string]any)
		GetNamespaceListsIncludingShared(ctx context.Context, namespace string) (lists map[string]any)
		GetSharedNamespaceLists(ctx context.Context) (lists map[string]any)
		IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) bool
		IsNamespaceOwnerOrAdminOrOperator(ctx context.Context, namespace string, userId int64) bool
		IsSharedNamespace(namespace string) bool
		AddNewNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		RemoveNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		QueryNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		QueryOwnNamespaceReturnRes(ctx context.Context) (retMsg string)
		AddNamespaceAdminReturnRes(ctx context.Context, namespace string, userId int64) (retMsg string)
		RemoveNamespaceAdminReturnRes(ctx context.Context, namespace string, userId int64) (retMsg string)
		ResetNamespaceAdminReturnRes(ctx context.Context, namespace string) (retMsg string)
		ChangeNamespaceOwnerReturnRes(ctx context.Context, namespace, ownerId string) (retMsg string)
		SetNamespacePropertyPublicReturnRes(ctx context.Context, namespace string, value bool) (retMsg string)
		IsNamespacePropertyPublic(ctx context.Context, namespace string) bool
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
