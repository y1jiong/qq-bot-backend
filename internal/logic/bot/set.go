package bot

import (
	"context"
	"github.com/bytedance/sonic/ast"
)

func (s *sBot) RewriteMessage(ctx context.Context, message string) {
	_, _ = s.reqJsonFromCtx(ctx).Set("message", ast.NewString(message))
}
