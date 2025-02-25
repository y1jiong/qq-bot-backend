package bot

import (
	"context"
	"sync"
)

var (
	connection = sync.Map{}
)

func (s *sBot) JoinConnection(ctx context.Context, key int64) {
	connection.Store(key, ctx)
}

func (s *sBot) LeaveConnection(key int64) {
	connection.Delete(key)
}

func (s *sBot) LoadConnection(key int64) context.Context {
	if v, ok := connection.Load(key); ok {
		if ctx, okay := v.(context.Context); okay {
			return ctx
		}
	}
	return nil
}
