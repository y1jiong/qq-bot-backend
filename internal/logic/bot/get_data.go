package bot

import (
	"context"
	"github.com/bytedance/sonic/ast"
	"qq-bot-backend/utility/segment"
)

func (s *sBot) getData(ctx context.Context) *ast.Node {
	return s.reqNodeFromCtx(ctx).Get("data")
}

func (s *sBot) getFileFromData(ctx context.Context) string {
	v, _ := s.getData(ctx).Get("file").StrictString()
	return v
}

func (s *sBot) getMessageIdFromData(ctx context.Context) int64 {
	v, _ := s.getData(ctx).Get("message_id").Int64()
	return v
}

func (s *sBot) MessageToFakeNode(userId int64, nickname, message string) map[string]any {
	return map[string]any{
		"type": "node",
		"data": map[string]any{
			"user_id":  userId,
			"nickname": nickname,
			"content":  segment.ParseMessage(message),
		},
	}
}
