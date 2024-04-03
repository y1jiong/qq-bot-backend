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
	IModule interface {
		TryApproveAddGroup(ctx context.Context) (catch bool)
		TryLockCard(ctx context.Context) (catch bool)
		TryAutoSetCard(ctx context.Context) (catch bool)
		TryKeywordRecall(ctx context.Context) (catch bool)
		TryGroupKeywordReply(ctx context.Context) (catch bool)
		TryLogLeave(ctx context.Context) (catch bool)
		TryLogApproval(ctx context.Context) (catch bool)
		TryUndoMessageRecall(ctx context.Context) (catch bool)
		TryKeywordReply(ctx context.Context) (catch bool)
		AutoLimit(ctx context.Context, kind, key string, limitTimes int, duration time.Duration) (limited bool, times int)
		MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string)
		AutoMute(ctx context.Context, kind string, groupId, userId int64, limitTimes, baseMinutes, limitMinutes int, duration time.Duration)
		WebhookGetHeadConnectOptionsTrace(ctx context.Context, method, url string) (body string, err error)
		WebhookPostPutPatchDelete(ctx context.Context, method, url string, payload any) (body string, err error)
	}
)

var (
	localModule IModule
)

func Module() IModule {
	if localModule == nil {
		panic("implement not found for interface IModule, forgot register?")
	}
	return localModule
}

func RegisterModule(i IModule) {
	localModule = i
}
