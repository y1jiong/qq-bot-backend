package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
	"sync"
)

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

func (s *sBot) GetTimestamp(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("time").MustInt64()
}

func (s *sBot) GetOperatorId(ctx context.Context) int64 {
	return s.reqJsonFromCtx(ctx).Get("operator_id").MustInt64()
}

func (s *sBot) GetGroupMemberList(ctx context.Context, groupId int64, noCache ...bool) (members []any, err error) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "get_group_member_list")
	// echo sign
	echoSign := guid.S()
	resJson.Set("echo", echoSign)
	// 参数
	params := make(map[string]any)
	params["group_id"] = groupId
	if len(noCache) > 0 && noCache[0] {
		params["no_cache"] = true
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		defer wg.Done()
		if err = s.defaultEchoProcess(rsyncCtx); err != nil {
			return
		}
		received := s.getData(rsyncCtx)
		members = received.MustArray()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) GetCardOldNew(ctx context.Context) (oldCard, newCard string) {
	oldCard = s.reqJsonFromCtx(ctx).Get("card_old").MustString()
	newCard = s.reqJsonFromCtx(ctx).Get("card_new").MustString()
	return
}

func (s *sBot) RequestMessage(ctx context.Context, messageId int64) (messageMap map[string]any, err error) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "get_msg")
	// echo sign
	echoSign := guid.S()
	resJson.Set("echo", echoSign)
	// 参数
	params := make(map[string]any)
	params["message_id"] = gconv.String(messageId)
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		defer wg.Done()
		if err = s.defaultEchoProcess(rsyncCtx); err != nil {
			return
		}
		received := s.getData(rsyncCtx)
		messageMap = received.MustMap()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) GetGroupInfo(ctx context.Context, groupId int64, noCache ...bool) (infoMap map[string]any, err error) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "get_group_info")
	// echo sign
	echoSign := guid.S()
	resJson.Set("echo", echoSign)
	// 参数
	params := make(map[string]any)
	params["group_id"] = groupId
	if len(noCache) > 0 && noCache[0] {
		params["no_cache"] = true
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		defer wg.Done()
		if err = s.defaultEchoProcess(rsyncCtx); err != nil {
			return
		}
		received := s.getData(rsyncCtx)
		infoMap = received.MustMap()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}
