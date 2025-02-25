package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/utility/segment"
)

func (s *sBot) RewriteMessage(ctx context.Context, message string) {
	_, _ = s.reqNodeFromCtx(ctx).Set("raw_message", ast.NewString(message))
}

func (s *sBot) SetHistory(ctx context.Context, history string) error {
	const historyKey = "_history"
	node := s.reqNodeFromCtx(ctx)
	if !node.Get(historyKey).Valid() {
		_, _ = node.Set(historyKey, ast.NewNull())
	}
	if node.Get(historyKey).Get(history).Valid() {
		return errors.New("history already exists")
	}
	_, _ = node.Get(historyKey).Set(history, ast.NewNull())
	return nil
}

func (s *sBot) getMessageAstNodeCacheKey(ctx context.Context) string {
	return cacheKeyMsgIdPrefix + gconv.String(s.GetSelfId(ctx)) + "_" + gconv.String(s.GetMsgId(ctx))
}

func (s *sBot) CacheMessageAstNode(ctx context.Context) {
	_ = gcache.Set(ctx,
		s.getMessageAstNodeCacheKey(ctx),
		s.reqNodeFromCtx(ctx),
		messageContextTTL,
	)
}

func (s *sBot) tryMessageSegmentToString(ctx context.Context) {
	node := s.reqNodeFromCtx(ctx)

	messageNode := node.Get("message")

	if !messageNode.Exists() || messageNode.TypeSafe() != ast.V_ARRAY {
		return
	}
	_, _ = node.Set("_is_message_segment", ast.NewObject([]ast.Pair{}))

	if rawMsgNode := node.Get("raw_message"); rawMsgNode.Exists() && rawMsgNode.TypeSafe() == ast.V_STRING {
		_, _ = node.Set("message", *rawMsgNode)
		return
	}

	jsonBytes, err := messageNode.MarshalJSON()
	if err != nil {
		return
	}

	segments, err := segment.ParseJSON(jsonBytes)
	if err != nil {
		return
	}

	_, _ = node.Set("message", ast.NewString(segments.String()))
}
