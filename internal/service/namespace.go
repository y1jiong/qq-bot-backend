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
		IsNamespaceOwnerOrAdmin(ctx context.Context, namespace string, userId int64) (yes bool)
		AddNamespaceList(ctx context.Context, namespace, listName string)
		RemoveNamespaceList(ctx context.Context, namespace, listName string)
		GetNamespaceList(ctx context.Context, namespace string) (lists map[string]any)
		GetNamespaceListIncludingPublic(ctx context.Context, namespace string) (lists map[string]any)
		IsPublicNamespace(namespace string) (yes bool)
		AddNewNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		RemoveNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		QueryNamespaceReturnRes(ctx context.Context, namespace string) (retMsg string)
		QueryOwnNamespaceReturnRes(ctx context.Context) (retMsg string)
		AddNamespaceAdminReturnRes(ctx context.Context, namespace string, userId int64) (retMsg string)
		RemoveNamespaceAdminReturnRes(ctx context.Context, namespace string, userId int64) (retMsg string)
		ResetNamespaceAdminReturnRes(ctx context.Context, namespace string) (retMsg string)
		ChangeNamespaceOwnerReturnRes(ctx context.Context, namespace, ownerId string) (retMsg string)
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
