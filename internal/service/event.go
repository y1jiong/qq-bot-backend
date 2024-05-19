// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IEvent interface {
		TryApproveAddGroup(ctx context.Context) (catch bool)
		TryLockCard(ctx context.Context) (catch bool)
		TryAutoSetCard(ctx context.Context) (catch bool)
		TryKeywordRecall(ctx context.Context) (catch bool)
		TryGroupKeywordReply(ctx context.Context) (catch bool)
		TryLogLeave(ctx context.Context) (catch bool)
		TryLogApproval(ctx context.Context) (catch bool)
		TryUndoMessageRecall(ctx context.Context) (catch bool)
		TryKeywordReply(ctx context.Context) (catch bool)
	}
)

var (
	localEvent IEvent
)

func Event() IEvent {
	if localEvent == nil {
		panic("implement not found for interface IEvent, forgot register?")
	}
	return localEvent
}

func RegisterEvent(i IEvent) {
	localEvent = i
}
