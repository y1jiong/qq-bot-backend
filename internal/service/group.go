// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IGroup interface {
		BindNamespace(ctx context.Context, groupId int64, namespace string)
		Unbind(ctx context.Context, groupId int64)
		BindQuery(ctx context.Context, groupId int64)
		IsGroupBindNamespaceOwnerOrAdmin(ctx context.Context, groupId, userId int64) (yes bool)
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
