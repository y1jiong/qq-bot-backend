// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ICrontab interface {
		Run(ctx context.Context)
		GlanceReturnRes(ctx context.Context) (retMsg string)
		QueryReturnRes(ctx context.Context, name string) (retMsg string)
		AddReturnRes(ctx context.Context, name string, expr string, botId int64, reqJSON []byte) (retMsg string)
		RemoveReturnRes(ctx context.Context, name string) (retMsg string)
	}
)

var (
	localCrontab ICrontab
)

func Crontab() ICrontab {
	if localCrontab == nil {
		panic("implement not found for interface ICrontab, forgot register?")
	}
	return localCrontab
}

func RegisterCrontab(i ICrontab) {
	localCrontab = i
}
