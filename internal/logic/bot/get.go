package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"qq-bot-backend/internal/service"
	"sync"
)

func (s *sBot) isApiReq(ctx context.Context) bool {
	return s.reqNodeFromCtx(ctx).Get("api_req").Exists()
}

func (s *sBot) isMessageSegment(ctx context.Context) bool {
	return s.reqNodeFromCtx(ctx).Get("_is_message_segment").Exists()
}

func (s *sBot) getEcho(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("echo").StrictString()
	return v
}

func (s *sBot) getEchoStatus(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("status").StrictString()
	return v
}

func (s *sBot) getEchoFailedMsg(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("wording").StrictString()
	return v
}

func (s *sBot) GetPostType(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("post_type").StrictString()
	return v
}

func (s *sBot) GetMsgType(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("message_type").StrictString()
	return v
}

func (s *sBot) GuessMsgType(groupId int64) string {
	if groupId != 0 {
		return "group"
	}
	return "private"
}

func (s *sBot) GetRequestType(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("request_type").StrictString()
	return v
}

func (s *sBot) GetNoticeType(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("notice_type").StrictString()
	return v
}

func (s *sBot) GetSubType(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("sub_type").StrictString()
	return v
}

func (s *sBot) GetMsgId(ctx context.Context) int64 {
	v, _ := s.reqNodeFromCtx(ctx).Get("message_id").StrictInt64()
	return v
}

func (s *sBot) GetMessage(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("raw_message").StrictString()
	if v == "" {
		v, _ = s.reqNodeFromCtx(ctx).Get("message").StrictString()
	}
	return v
}

func (s *sBot) GetUserId(ctx context.Context) int64 {
	v, _ := s.reqNodeFromCtx(ctx).Get("user_id").StrictInt64()
	return v
}

func (s *sBot) GetGroupId(ctx context.Context) int64 {
	v, _ := s.reqNodeFromCtx(ctx).Get("group_id").StrictInt64()
	return v
}

func (s *sBot) GetComment(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("comment").StrictString()
	return v
}

func (s *sBot) GetFlag(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("flag").StrictString()
	return v
}

func (s *sBot) GetTimestamp(ctx context.Context) int64 {
	v, _ := s.reqNodeFromCtx(ctx).Get("time").StrictInt64()
	return v
}

func (s *sBot) GetOperatorId(ctx context.Context) int64 {
	v, _ := s.reqNodeFromCtx(ctx).Get("operator_id").StrictInt64()
	return v
}

func (s *sBot) GetSelfId(ctx context.Context) int64 {
	v, _ := s.reqNodeFromCtx(ctx).Get("self_id").StrictInt64()
	return v
}

func (s *sBot) GetNickname(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("sender").Get("nickname").StrictString()
	return v
}

func (s *sBot) GetCard(ctx context.Context) string {
	v, _ := s.reqNodeFromCtx(ctx).Get("sender").Get("card").StrictString()
	return v
}

func (s *sBot) GetCardOrNickname(ctx context.Context) string {
	if card := s.GetCard(ctx); card != "" {
		return card
	}
	return s.GetNickname(ctx)
}

func (s *sBot) GetCardOldNew(ctx context.Context) (oldCard, newCard string) {
	oldCard, _ = s.reqNodeFromCtx(ctx).Get("card_old").StrictString()
	newCard, _ = s.reqNodeFromCtx(ctx).Get("card_new").StrictString()
	return
}

func (s *sBot) GetGroupMemberInfo(ctx context.Context, groupId, userId int64, noCache ...bool) (member *ast.Node, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.GetGroupMemberInfo")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("get_group_member_info.group_id", groupId),
		attribute.Int64("get_group_member_info.user_id", userId),
	)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId int64 `json:"group_id"`
			UserId  int64 `json:"user_id"`
			NoCache bool  `json:"no_cache"`
		} `json:"params"`
	}{
		Action: "get_group_member_info",
		Echo:   echoSign,
		Params: struct {
			GroupId int64 `json:"group_id"`
			UserId  int64 `json:"user_id"`
			NoCache bool  `json:"no_cache"`
		}{
			GroupId: groupId,
			UserId:  userId,
			NoCache: false,
		},
	}
	if len(noCache) > 0 && noCache[0] {
		req.Params.NoCache = true
	}
	reqJSON, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			if !req.Params.NoCache {
				member, err = s.GetGroupMemberInfo(ctx, groupId, userId, true)
			}
			return
		}
		member = s.getData(asyncCtx)
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) GetGroupMemberList(ctx context.Context, groupId int64, noCache ...bool) (members []any, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.GetGroupMemberList")
	defer span.End()
	span.SetAttributes(attribute.Int64("get_group_member_list.group_id", groupId))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId int64 `json:"group_id"`
			NoCache bool  `json:"no_cache"`
		} `json:"params"`
	}{
		Action: "get_group_member_list",
		Echo:   echoSign,
		Params: struct {
			GroupId int64 `json:"group_id"`
			NoCache bool  `json:"no_cache"`
		}{
			GroupId: groupId,
			NoCache: false,
		},
	}
	if len(noCache) > 0 && noCache[0] {
		req.Params.NoCache = true
	}
	reqJSON, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			if !req.Params.NoCache {
				members, err = s.GetGroupMemberList(ctx, groupId, true)
			}
			return
		}
		received := s.getData(asyncCtx)
		members, _ = received.Array()
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) RequestMessageFromCache(ctx context.Context, messageId int64) (messageMap map[string]any, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.RequestMessageFromCache")
	defer span.End()
	span.SetAttributes(attribute.Int64("request_message_from_cache.message_id", messageId))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	v, err := gcache.Get(ctx, s.getMessageAstNodeCacheKey(ctx))
	if err != nil {
		return
	}

	if v == nil || v.IsNil() {
		messageMap = make(map[string]any)
		return
	}

	node, ok := v.Val().(*ast.Node)
	if !ok {
		messageMap = make(map[string]any)
		return
	}

	return node.Map()
}

func (s *sBot) RequestMessage(ctx context.Context, messageId int64) (messageMap map[string]any, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.RequestMessage")
	defer span.End()
	span.SetAttributes(attribute.Int64("request_message.message_id", messageId))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			MessageId int64 `json:"message_id"`
		} `json:"params"`
	}{
		Action: "get_msg",
		Echo:   echoSign,
		Params: struct {
			MessageId int64 `json:"message_id"`
		}{
			MessageId: messageId,
		},
	}
	reqJSON, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
		received := s.getData(asyncCtx)
		messageMap, _ = received.Map()
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) GetGroupInfo(ctx context.Context, groupId int64, noCache ...bool) (infoMap map[string]any, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.GetGroupInfo")
	defer span.End()
	span.SetAttributes(attribute.Int64("get_group_info.group_id", groupId))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId int64 `json:"group_id"`
			NoCache bool  `json:"no_cache"`
		} `json:"params"`
	}{
		Action: "get_group_info",
		Echo:   echoSign,
		Params: struct {
			GroupId int64 `json:"group_id"`
			NoCache bool  `json:"no_cache"`
		}{
			GroupId: groupId,
			NoCache: false,
		},
	}
	if len(noCache) > 0 && noCache[0] {
		req.Params.NoCache = true
	}
	reqJSON, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			if !req.Params.NoCache {
				infoMap, err = s.GetGroupInfo(ctx, groupId, true)
			}
			return
		}
		received := s.getData(asyncCtx)
		infoMap, _ = received.Map()
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) GetLoginInfo(ctx context.Context) (userId int64, nickname string) {
	ctx, span := gtrace.NewSpan(ctx, "bot.GetLoginInfo")
	defer span.End()
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
	}{
		Action: "get_login_info",
		Echo:   echoSign,
	}
	reqJSON, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
		received := s.getData(asyncCtx)
		userId, _ = received.Get("user_id").StrictInt64()
		nickname, _ = received.Get("nickname").StrictString()
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) IsGroupOwnerOrAdmin(ctx context.Context) bool {
	role, _ := s.reqNodeFromCtx(ctx).Get("sender").Get("role").StrictString()
	// lazy load user role
	if role == "" {
		member, err := s.GetGroupMemberInfo(ctx, s.GetGroupId(ctx), s.GetUserId(ctx))
		if err != nil {
			g.Log().Warning(ctx, err)
			return false
		}
		role, err = member.Get("role").StrictString()
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}
		_, _ = s.reqNodeFromCtx(ctx).Set("sender", *member)
	}
	return role == "owner" || role == "admin"
}

func (s *sBot) IsGroupOwnerOrAdminOrSysTrusted(ctx context.Context) bool {
	return s.IsGroupOwnerOrAdmin(ctx) || service.User().IsSystemTrustedUser(ctx, gconv.Int64(s.GetUserId(ctx)))
}

func (s *sBot) GetVersionInfo(ctx context.Context) (appName, appVersion, protocolVersion string, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.GetVersionInfo")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
	}{
		Action: "get_version_info",
		Echo:   echoSign,
	}
	reqJSON, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
		received := s.getData(asyncCtx)
		appName, _ = received.Get("app_name").StrictString()
		appVersion, _ = received.Get("app_version").StrictString()
		protocolVersion, _ = received.Get("protocol_version").StrictString()
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) GetLikes(ctx context.Context) []map[string]any {
	v, _ := s.reqNodeFromCtx(ctx).Get("likes").Array()
	if len(v) == 0 {
		return nil
	}
	likes := make([]map[string]any, 0, len(v))
	for _, item := range v {
		likes = append(likes, gconv.Map(item))
	}
	return likes
}
