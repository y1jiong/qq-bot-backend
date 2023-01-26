// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IList interface {
		AddList(ctx context.Context, listName, namespace string)
		RemoveList(ctx context.Context, listName string)
		QueryList(ctx context.Context, listName string)
		GetList(ctx context.Context, listName string) (list map[string]any)
		AddListData(ctx context.Context, listName, key string, value ...string)
		RemoveListData(ctx context.Context, listName, key string)
		ResetListData(ctx context.Context, listName string)
	}
)

var (
	localList IList
)

func List() IList {
	if localList == nil {
		panic("implement not found for interface IList, forgot register?")
	}
	return localList
}

func RegisterList(i IList) {
	localList = i
}
