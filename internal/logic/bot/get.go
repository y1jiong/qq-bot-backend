package bot

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
	"sync"
)

func (s *sBot) isApiReq(ctx context.Context) bool {
	return s.reqJsonFromCtx(ctx).Get("api_req").Exists()
}

func (s *sBot) getEcho(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("echo").StrictString()
	return v
}

func (s *sBot) getEchoStatus(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("status").StrictString()
	return v
}

func (s *sBot) getEchoFailedMsg(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("wording").StrictString()
	return v
}

func (s *sBot) GetPostType(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("post_type").StrictString()
	return v
}

func (s *sBot) GetMsgType(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("message_type").StrictString()
	return v
}

func (s *sBot) GetRequestType(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("request_type").StrictString()
	return v
}

func (s *sBot) GetNoticeType(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("notice_type").StrictString()
	return v
}

func (s *sBot) GetSubType(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("sub_type").StrictString()
	return v
}

func (s *sBot) GetMsgId(ctx context.Context) int64 {
	v, _ := s.reqJsonFromCtx(ctx).Get("message_id").StrictInt64()
	return v
}

func (s *sBot) GetMessage(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("message").StrictString()
	return v
}

func (s *sBot) GetUserId(ctx context.Context) int64 {
	v, _ := s.reqJsonFromCtx(ctx).Get("user_id").StrictInt64()
	return v
}

func (s *sBot) GetGroupId(ctx context.Context) int64 {
	v, _ := s.reqJsonFromCtx(ctx).Get("group_id").StrictInt64()
	return v
}

func (s *sBot) GetComment(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("comment").StrictString()
	return v
}

func (s *sBot) GetFlag(ctx context.Context) string {
	v, _ := s.reqJsonFromCtx(ctx).Get("flag").StrictString()
	return v
}

func (s *sBot) GetTimestamp(ctx context.Context) int64 {
	v, _ := s.reqJsonFromCtx(ctx).Get("time").StrictInt64()
	return v
}

func (s *sBot) GetOperatorId(ctx context.Context) int64 {
	v, _ := s.reqJsonFromCtx(ctx).Get("operator_id").StrictInt64()
	return v
}

func (s *sBot) GetSelfId(ctx context.Context) int64 {
	v, _ := s.reqJsonFromCtx(ctx).Get("self_id").StrictInt64()
	return v
}

func (s *sBot) GetGroupMemberInfo(ctx context.Context, groupId, userId int64) (member ast.Node, err error) {
	// echo sign
	echoSign := guid.S()
	// 参数
	res := struct {
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
	resJson, err := sonic.ConfigStd.Marshal(res)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		defer wg.Done()
		if err = s.defaultEchoProcess(rsyncCtx); err != nil {
			return
		}
		member = *s.getData(rsyncCtx)
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, resJson)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) GetGroupMemberList(ctx context.Context, groupId int64, noCache ...bool) (members []any, err error) {
	// echo sign
	echoSign := guid.S()
	// 参数
	res := struct {
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
			NoCache: true,
		},
	}
	resJson, err := sonic.ConfigStd.Marshal(res)
	if err != nil {
		g.Log().Error(ctx, err)
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
		members, _ = received.Array()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, resJson)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) GetCardOldNew(ctx context.Context) (oldCard, newCard string) {
	oldCard, _ = s.reqJsonFromCtx(ctx).Get("card_old").StrictString()
	newCard, _ = s.reqJsonFromCtx(ctx).Get("card_new").StrictString()
	return
}

func (s *sBot) RequestMessage(ctx context.Context, messageId int64) (messageMap map[string]any, err error) {
	// echo sign
	echoSign := guid.S()
	// 参数
	res := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			MessageId string `json:"message_id"`
		} `json:"params"`
	}{
		Action: "get_msg",
		Echo:   echoSign,
		Params: struct {
			MessageId string `json:"message_id"`
		}{
			MessageId: gconv.String(messageId),
		},
	}
	resJson, err := sonic.ConfigStd.Marshal(res)
	if err != nil {
		g.Log().Error(ctx, err)
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
		messageMap, _ = received.Map()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, resJson)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) GetGroupInfo(ctx context.Context, groupId int64, noCache ...bool) (infoMap map[string]any, err error) {
	// echo sign
	echoSign := guid.S()
	// 参数
	res := struct {
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
		res.Params.NoCache = true
	}
	resJson, err := sonic.ConfigStd.Marshal(res)
	if err != nil {
		g.Log().Error(ctx, err)
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
		infoMap, _ = received.Map()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, resJson)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) GetLoginInfo(ctx context.Context) (userId int64, nickname string) {
	// echo sign
	echoSign := guid.S()
	// 参数
	res := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
	}{
		Action: "get_login_info",
		Echo:   echoSign,
	}
	resJson, err := sonic.ConfigStd.Marshal(res)
	if err != nil {
		g.Log().Error(ctx, err)
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
		userId, _ = received.Get("user_id").StrictInt64()
		nickname, _ = received.Get("nickname").StrictString()
	}
	wg.Add(1)
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, resJson)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	wg.Wait()
	return
}

func (s *sBot) IsGroupOwnerOrAdmin(ctx context.Context) (yes bool) {
	role, _ := s.reqJsonFromCtx(ctx).Get("sender").Get("role").StrictString()
	// lazy load user role
	if role == "" {
		member, err := s.GetGroupMemberInfo(ctx, s.GetGroupId(ctx), s.GetUserId(ctx))
		if err != nil {
			g.Log().Warning(ctx, err)
			return
		}
		role, err = member.Get("role").StrictString()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		params := []ast.Pair{
			{
				Key:   "role",
				Value: ast.NewString(role),
			},
		}
		_, err = s.reqJsonFromCtx(ctx).Set("sender", ast.NewObject(params))
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	return role == "owner" || role == "admin"
}
