// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IList interface {
		GetListData(ctx context.Context, listName string) (listMap map[string]any)
		AppendListData(ctx context.Context, listName string, newMap map[string]any) (n int, err error)
		AddListReturnRes(ctx context.Context, listName, namespace string) (retMsg string)
		RemoveListReturnRes(ctx context.Context, listName string) (retMsg string)
		RecoverListReturnRes(ctx context.Context, listName string) (retMsg string)
		ExportListReturnRes(ctx context.Context, listName string) (retMsg string)
		QueryListLenReturnRes(ctx context.Context, listName string) (retMsg string)
		QueryListReturnRes(ctx context.Context, listName string, keys ...string) (retMsg string)
		AddListDataReturnRes(ctx context.Context, listName, key string, value ...string) (retMsg string)
		RemoveListDataReturnRes(ctx context.Context, listName, key string) (retMsg string)
		ResetListDataReturnRes(ctx context.Context, listName string) (retMsg string)
		SetListDataReturnRes(ctx context.Context, listName, newListStr string) (retMsg string)
		AppendListDataReturnRes(ctx context.Context, listName, newListStr string) (retMsg string)
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
