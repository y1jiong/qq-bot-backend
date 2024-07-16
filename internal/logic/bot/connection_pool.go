package bot

import (
	"context"
	"sync"
)

var (
	connectionPool = sync.Map{}
)

func (s *sBot) JoinConnectionPool(ctx context.Context, key int64) {
	connectionPool.Store(key, ctx)
}

func (s *sBot) LeaveConnectionPool(key int64) {
	connectionPool.Delete(key)
}

func (s *sBot) LoadConnectionPool(key int64) context.Context {
	if v, ok := connectionPool.Load(key); ok {
		return v.(context.Context)
	}
	return nil
}
