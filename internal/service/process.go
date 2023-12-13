// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IProcess interface {
		IsEnabledBotProcess() bool
		PauseBotProcess() bool
		ContinueBotProcess() bool
		Process(ctx context.Context)
	}
)

var (
	localProcess IProcess
)

func Process() IProcess {
	if localProcess == nil {
		panic("implement not found for interface IProcess, forgot register?")
	}
	return localProcess
}

func RegisterProcess(i IProcess) {
	localProcess = i
}
