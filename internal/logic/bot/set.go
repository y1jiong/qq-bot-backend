package bot

import (
	"context"
	"github.com/bytedance/sonic/ast"
)

func (s *sBot) RewriteMessage(ctx context.Context, message string) error {
	_, err := s.reqJsonFromCtx(ctx).Set("message", ast.NewString(message))
	return err
}
