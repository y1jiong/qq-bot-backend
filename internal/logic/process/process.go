package process

import (
	"context"
	"qq-bot-backend/internal/service"
	"sync/atomic"
)

type sProcess struct{}

func New() *sProcess {
	return &sProcess{}
}

func init() {
	service.RegisterProcess(New())
}

const (
	enabled  int32 = 1
	disabled int32 = 0
)

var (
	// 默认启用处理
	botProcessState = enabled
)

func (s *sProcess) IsBotProcess() bool {
	return atomic.LoadInt32(&botProcessState) == enabled
}

func (s *sProcess) PauseBotProcess() bool {
	return atomic.CompareAndSwapInt32(&botProcessState, botProcessState, disabled)
}

func (s *sProcess) ContinueBotProcess() bool {
	return atomic.CompareAndSwapInt32(&botProcessState, botProcessState, enabled)
}

func (s *sProcess) Process(ctx context.Context) {
	// 优先处理命令
	if service.Command().TryCommand(ctx) {
		// 处理成功放弃后续逻辑
		return
	}
	// 是否暂停处理
	if !s.IsBotProcess() {
		return
	}
	// 捕捉 echo
	if service.Bot().CatchEcho(ctx) {
		return
	}
	// 处理分支
	switch service.Bot().GetPostType(ctx) {
	case "message":
		processMessage(ctx)
	case "request":
		processRequest(ctx)
	case "notice":
		processNotice(ctx)
	}
}
