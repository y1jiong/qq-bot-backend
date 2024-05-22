package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic/ast"
)

func (s *sBot) RewriteMessage(ctx context.Context, message string) {
	_, _ = s.reqJsonFromCtx(ctx).Set("raw_message", ast.NewString(message))
}

func (s *sBot) SetHistory(ctx context.Context, history string) error {
	const historyKey = "_history"
	node := s.reqJsonFromCtx(ctx)
	if !node.Get(historyKey).Valid() {
		_, _ = node.Set(historyKey, ast.NewNull())
	}
	if node.Get(historyKey).Get(history).Valid() {
		return errors.New("history already exists")
	}
	_, _ = node.Get(historyKey).Set(history, ast.NewNull())
	return nil
}
