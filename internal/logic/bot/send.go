package bot

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
)

func (s *sBot) SendMessage(ctx context.Context, messageType string, uid, gid int64, msg string, plain bool) {
	// 参数校验
	if uid == 0 && gid == 0 {
		return
	}
	if gid != 0 {
		uid = 0
	}
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			MessageType string `json:"message_type,omitempty"`
			Message     string `json:"message"`
			AutoEscape  bool   `json:"auto_escape,omitempty"`
			UserId      int64  `json:"user_id,omitempty"`
			GroupId     int64  `json:"group_id,omitempty"`
		} `json:"params"`
	}{
		Action: "send_msg",
		Params: struct {
			MessageType string `json:"message_type,omitempty"`
			Message     string `json:"message"`
			AutoEscape  bool   `json:"auto_escape,omitempty"`
			UserId      int64  `json:"user_id,omitempty"`
			GroupId     int64  `json:"group_id,omitempty"`
		}{
			MessageType: messageType,
			Message:     msg,
			AutoEscape:  plain,
			UserId:      uid,
			GroupId:     gid,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
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

func (s *sBot) SendPlainMsgIfNotApiReq(ctx context.Context, msg string) {
	if s.isApiReq(ctx) {
		return
	}
	s.SendPlainMsg(ctx, msg)
}

func (s *sBot) SendFileToGroup(ctx context.Context, gid int64, filePath, name, folder string) {
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			GroupId int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder,omitempty"`
		} `json:"params"`
	}{
		Action: "upload_group_file",
		Params: struct {
			GroupId int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder,omitempty"`
		}{
			GroupId: gid,
			File:    filePath,
			Name:    name,
			Folder:  folder,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendFileToUser(ctx context.Context, uid int64, filePath, name string) {
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			UserId int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		} `json:"params"`
	}{
		Action: "upload_private_file",
		Params: struct {
			UserId int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		}{
			UserId: uid,
			File:   filePath,
			Name:   name,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendFile(ctx context.Context, name, url string) {
	// echo sign
	echoSign := guid.S()
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			Url string `json:"url"`
		} `json:"params"`
	}{
		Action: "download_file",
		Echo:   echoSign,
		Params: struct {
			Url string `json:"url"`
		}{
			Url: url,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if err = s.defaultEchoProcess(rsyncCtx); err != nil {
			s.SendPlainMsgIfNotApiReq(ctx, err.Error())
			return
		}
		filePath := s.getFileFromData(rsyncCtx)
		groupId := s.GetGroupId(ctx)
		if groupId != 0 {
			s.SendFileToGroup(ctx, groupId, filePath, name, "")
			return
		}
		userId := s.GetUserId(ctx)
		s.SendFileToUser(ctx, userId, filePath, name)
	}
	timeout := func(ctx context.Context) {
		s.SendPlainMsgIfNotApiReq(ctx, "上传文件超时")
	}
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback, timeout)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) ApproveJoinGroup(ctx context.Context, flag, subType string, approve bool, reason string) {
	// 参数校验
	if approve {
		reason = ""
	}
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			Flag    string `json:"flag"`
			SubType string `json:"sub_type"`
			Approve bool   `json:"approve"`
			Reason  string `json:"reason,omitempty"`
		} `json:"params"`
	}{
		Action: "set_group_add_request",
		Params: struct {
			Flag    string `json:"flag"`
			SubType string `json:"sub_type"`
			Approve bool   `json:"approve"`
			Reason  string `json:"reason,omitempty"`
		}{
			Flag:    flag,
			SubType: subType,
			Approve: approve,
			Reason:  reason,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SetModel(ctx context.Context, model string) {
	// echo sign
	echoSign := guid.S()
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			Model     string `json:"model"`
			ModelShow string `json:"model_show"`
		} `json:"params"`
	}{
		Action: "_set_model_show",
		Echo:   echoSign,
		Params: struct {
			Model     string `json:"model"`
			ModelShow string `json:"model_show"`
		}{
			Model:     model,
			ModelShow: model,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	callback := func(ctx context.Context, rsyncCtx context.Context) {
		if err = s.defaultEchoProcess(rsyncCtx); err != nil {
			s.SendPlainMsgIfNotApiReq(ctx, err.Error())
			return
		}
		s.SendPlainMsgIfNotApiReq(ctx, "已更改机型为 '"+model+"'")
	}
	timeout := func(ctx context.Context) {
		s.SendPlainMsgIfNotApiReq(ctx, "更改机型超时")
	}
	// echo
	err = s.pushEchoCache(ctx, echoSign, callback, timeout)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) RecallMessage(ctx context.Context, msgId int64) {
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			MessageId int64 `json:"message_id"`
		} `json:"params"`
	}{
		Action: "delete_msg",
		Params: struct {
			MessageId int64 `json:"message_id"`
		}{
			MessageId: msgId,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) MutePrototype(ctx context.Context, groupId, userId int64, seconds int) {
	// 参数校验
	if seconds > 2591940 {
		// 不大于 29 天 23 小时 59 分钟
		// (30*24*60-1)*60=2591940 秒
		seconds = 2591940
	}
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			GroupId  int64 `json:"group_id"`
			UserId   int64 `json:"user_id"`
			Duration int   `json:"duration"`
		} `json:"params"`
	}{
		Action: "set_group_ban",
		Params: struct {
			GroupId  int64 `json:"group_id"`
			UserId   int64 `json:"user_id"`
			Duration int   `json:"duration"`
		}{
			GroupId:  groupId,
			UserId:   userId,
			Duration: seconds,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) Mute(ctx context.Context, seconds int) {
	s.MutePrototype(ctx, s.GetGroupId(ctx), s.GetUserId(ctx), seconds)
}

func (s *sBot) SetGroupCard(ctx context.Context, groupId, userId int64, card string) {
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			GroupId int64  `json:"group_id"`
			UserId  int64  `json:"user_id"`
			Card    string `json:"card"`
		} `json:"params"`
	}{
		Action: "set_group_card",
		Params: struct {
			GroupId int64  `json:"group_id"`
			UserId  int64  `json:"user_id"`
			Card    string `json:"card"`
		}{
			GroupId: groupId,
			UserId:  userId,
			Card:    card,
		},
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) Kick(ctx context.Context, groupId, userId int64, reject ...bool) {
	// 参数
	req := struct {
		Action string `json:"action"`
		Params struct {
			GroupId          int64 `json:"group_id"`
			UserId           int64 `json:"user_id"`
			RejectAddRequest bool  `json:"reject_add_request,omitempty"`
		} `json:"params"`
	}{
		Action: "set_group_kick",
		Params: struct {
			GroupId          int64 `json:"group_id"`
			UserId           int64 `json:"user_id"`
			RejectAddRequest bool  `json:"reject_add_request,omitempty"`
		}{
			GroupId:          groupId,
			UserId:           userId,
			RejectAddRequest: false,
		},
	}
	if len(reject) > 0 && reject[0] {
		req.Params.RejectAddRequest = true
	}
	reqJson, err := sonic.ConfigStd.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	err = s.writeMessage(ctx, websocket.TextMessage, reqJson)
	if err != nil {
		g.Log().Warning(ctx, err)
	}
}
