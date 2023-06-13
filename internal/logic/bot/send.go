package bot

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
)

func (s *sBot) SendMessage(ctx context.Context, messageType string, uid, gid int64, msg string, plain bool) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "send_msg")
	// 参数
	params := make(map[string]any)
	params["message_type"] = messageType
	params["message"] = msg
	if plain {
		// 以纯文本方式发送
		params["auto_escape"] = true
	}
	if uid == 0 && gid == 0 {
		return
	}
	if gid != 0 {
		params["group_id"] = gid
	} else {
		params["user_id"] = uid
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendPlainMsg(ctx context.Context, msg string) {
	s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, true)
}

func (s *sBot) SendMsg(ctx context.Context, msg string) {
	s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, false)
}

func (s *sBot) SendFileToGroup(ctx context.Context, gid int64, filePath, name, folder string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "upload_group_file")
	// 参数
	params := make(map[string]any)
	params["group_id"] = gid
	params["file"] = filePath
	params["name"] = name
	if folder != "" {
		params["folder"] = folder
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendFileToUser(ctx context.Context, uid int64, filePath, name string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "upload_private_file")
	// 参数
	params := make(map[string]any)
	params["user_id"] = uid
	params["file"] = filePath
	params["name"] = name
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendFile(ctx context.Context, name, url string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "download_file")
	// echo sign
	echoSign := guid.S()
	resJson.Set("echo", echoSign)
	// 参数
	params := make(map[string]any)
	params["url"] = url
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if s.DefaultEchoProcess(ctx, rsyncCtx) {
			return
		}
		filePath := s.GetFileFromData(rsyncCtx)
		groupId := s.GetGroupId(ctx)
		if groupId != 0 {
			s.SendFileToGroup(ctx, groupId, filePath, name, "")
			return
		}
		userId := s.GetUserId(ctx)
		s.SendFileToUser(ctx, userId, filePath, name)
	}
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) ApproveJoinGroup(ctx context.Context, flag, subType string, approve bool, reason string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "set_group_add_request")
	// 参数
	params := make(map[string]any)
	params["flag"] = flag
	params["sub_type"] = subType
	params["approve"] = approve
	// 当不予通过时，给出理由
	if !approve {
		params["reason"] = reason
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// 发送响应
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SetModel(ctx context.Context, model string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "_set_model_show")
	// echo sign
	echoSign := guid.S()
	resJson.Set("echo", echoSign)
	// 参数
	params := make(map[string]any)
	params["model"] = model
	params["model_show"] = model
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	// callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if s.DefaultEchoProcess(ctx, rsyncCtx) {
			return
		}
		s.SendPlainMsg(ctx, "已更改机型为 '"+model+"'")
	}
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) RecallMessage(ctx context.Context, msgId int64) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "delete_msg")
	// 参数
	params := make(map[string]any)
	params["message_id"] = msgId
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) MutePrototype(ctx context.Context, groupId, userId int64, seconds int) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "set_group_ban")
	// 参数
	params := make(map[string]any)
	params["group_id"] = groupId
	params["user_id"] = userId
	if seconds > 2591940 {
		// 不大于 29 天 23 小时 59 分钟
		// (30*24*60-1)*60=2591940 秒
		seconds = 2591940
	}
	params["duration"] = seconds
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) Mute(ctx context.Context, seconds int) {
	s.MutePrototype(ctx, s.GetGroupId(ctx), s.GetUserId(ctx), seconds)
}

func (s *sBot) SetGroupCard(ctx context.Context, groupId, userId int64, card string) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "set_group_card")
	// 参数
	params := make(map[string]any)
	params["group_id"] = groupId
	params["user_id"] = userId
	params["card"] = card
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) Kick(ctx context.Context, groupId, userId int64, reject ...bool) {
	// 初始化响应
	resJson := sj.New()
	resJson.Set("action", "set_group_kick")
	// 参数
	params := make(map[string]any)
	params["group_id"] = groupId
	params["user_id"] = userId
	if len(reject) > 0 && reject[0] {
		params["reject_add_request"] = true
	}
	// 参数打包
	resJson.Set("params", params)
	res, err := resJson.Encode()
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	err = s.webSocketFromCtx(ctx).WriteMessage(websocket.TextMessage, res)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}
