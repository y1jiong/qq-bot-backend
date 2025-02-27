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
		TryForward(ctx context.Context) (caught bool)
		TryApproveAddGroup(ctx context.Context) (caught bool)
		TryLockCard(ctx context.Context) (caught bool)
		TryAutoSetCard(ctx context.Context) (caught bool)
		TryKeywordRecall(ctx context.Context) (caught bool)
		TryGroupKeywordReply(ctx context.Context) (caught bool)
		TryLogLeave(ctx context.Context) (caught bool)
		TryLogApproval(ctx context.Context) (caught bool)
		TryUndoMessageRecall(ctx context.Context) (caught bool)
		TryKeywordReply(ctx context.Context) (caught bool)
		TryCascadingRecall(ctx context.Context) (caught bool)
		TryEmojiRecall(ctx context.Context) (caught bool)
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
