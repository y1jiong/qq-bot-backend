// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"time"
)

type (
	IUtil interface {
		AutoMute(ctx context.Context, kind string, groupId int64, userId int64, limitTimes int, baseMinutes int, limitMinutes int, duration time.Duration)
		IsOnKeywordLists(ctx context.Context, msg string, lists map[string]any) (in bool, hit string, value string)
		MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string)
	}
)

var (
	localUtil IUtil
)

func Util() IUtil {
	if localUtil == nil {
		panic("implement not found for interface IUtil, forgot register?")
	}
	return localUtil
}

func RegisterUtil(i IUtil) {
	localUtil = i
}
