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
		AutoLimit(ctx context.Context, kind, key string, limitTimes int, duration time.Duration) (limited bool, times int)
		ReverseSortedArrayFromMapKey(m map[string]any) (arr []string)
		AutoMute(ctx context.Context, kind string, groupId, userId int64, limitTimes, baseMinutes, limitMinutes int, duration time.Duration)
		MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string)
		IsOnKeywordLists(ctx context.Context, msg string, lists map[string]any) (in bool, hit, value string)
		WebhookGetHeadConnectOptionsTrace(ctx context.Context, header, method, url string) (statusCode int, contentType string, body []byte, err error)
		WebhookPostPutPatchDelete(ctx context.Context, header, method, url string, payload any) (statusCode int, contentType string, body []byte, err error)
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
