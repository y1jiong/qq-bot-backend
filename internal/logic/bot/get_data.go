package bot

import (
	"context"
	"github.com/bytedance/sonic/ast"
)

func (s *sBot) getData(ctx context.Context) *ast.Node {
	return s.reqJsonFromCtx(ctx).Get("data")
}

func (s *sBot) getFileFromData(ctx context.Context) string {
	v, _ := s.getData(ctx).Get("file").StrictString()
	return v
}
