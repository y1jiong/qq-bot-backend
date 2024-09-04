package process

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
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

func (s *sProcess) IsBotProcessEnabled() bool {
	return atomic.LoadInt32(&botProcessState) == enabled
}

func (s *sProcess) PauseBotProcess() bool {
	return atomic.CompareAndSwapInt32(&botProcessState, botProcessState, disabled)
}

func (s *sProcess) ContinueBotProcess() bool {
	return atomic.CompareAndSwapInt32(&botProcessState, botProcessState, enabled)
}

func (s *sProcess) Process(ctx context.Context) {
	if service.Bot().GetPostType(ctx) == "meta_event" {
		// 跳过处理元事件 心跳包 生命周期
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "process.Process")
	defer span.End()

	// 优先处理命令
	if catch, retMsg := service.Command().TryCommand(ctx, service.Bot().GetMessage(ctx)); catch {
		// 处理成功放弃后续逻辑
		if retMsg != "" {
			service.Bot().SendMsgCacheContext(ctx, retMsg)
		}
		return
	}
	// 是否暂停处理
	if !s.IsBotProcessEnabled() {
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
