package state

import "sync/atomic"

const (
	enabled  int32 = 1
	disabled int32 = 0
)

var (
	// bot 默认启用处理
	botProcessState = enabled
)

func (s *sState) IsBotProcess() bool {
	return atomic.LoadInt32(&botProcessState) == enabled
}

func (s *sState) PauseBotProcess() bool {
	return atomic.CompareAndSwapInt32(&botProcessState, botProcessState, disabled)
}

func (s *sState) ContinueBotProcess() bool {
	return atomic.CompareAndSwapInt32(&botProcessState, botProcessState, enabled)
}
