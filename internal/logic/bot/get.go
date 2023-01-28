package bot

import "context"

func (s *sBot) getEcho(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("echo").MustString()
}

func (s *sBot) getEchoStatus(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("status").MustString()
}

func (s *sBot) getEchoFailedMsg(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("wording").MustString()
}

func (s *sBot) GetPostType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("post_type").MustString()
}

func (s *sBot) GetMsgType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("message_type").MustString()
}

func (s *sBot) GetRequestType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("request_type").MustString()
}

func (s *sBot) GetNoticeType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("notice_type").MustString()
}

func (s *sBot) GetSubType(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("sub_type").MustString()
}

func (s *sBot) GetMsgId(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("message_id").MustInt64()
}

func (s *sBot) GetMessage(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("message").MustString()
}

func (s *sBot) GetUserId(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("user_id").MustInt64()
}

func (s *sBot) GetGroupId(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("group_id").MustInt64()
}

func (s *sBot) GetComment(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("comment").MustString()
}

func (s *sBot) GetFlag(ctx context.Context) string {
	return s.reqJsonFromCtx(ctx).Get("flag").MustString()
}
