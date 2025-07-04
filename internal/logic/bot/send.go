package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"qq-bot-backend/utility/segment"
	"sync"
)

func (s *sBot) SendMessage(ctx context.Context,
	userId, groupId int64,
	msg string,
	plain bool,
) (messageId int64, err error) {
	// 参数校验
	if userId == 0 && groupId == 0 {
		return 0, errors.New("userId 和 groupId 不能同时为 0")
	}
	if msg == "" {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "bot.SendMessage")
	defer span.End()
	span.SetAttributes(attribute.String("send_message.message", msg))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	if groupId != 0 {
		userId = 0
		span.SetAttributes(attribute.Int64("send_message.group_id", groupId))
		go func() {
			if err := s.MarkGroupMsgAsRead(ctx, groupId); err != nil {
				g.Log().Warning(ctx, err)
			}
		}()
	} else {
		span.SetAttributes(attribute.Int64("send_message.user_id", userId))
		go func() {
			if err := s.MarkPrivateMsgAsRead(ctx, userId); err != nil {
				g.Log().Warning(ctx, err)
			}
		}()
	}

	// 消息长度限制
	if len(msg) > messageLengthLimit {
		var botId int64
		var nickname string
		botId, nickname, err = s.GetLoginInfo(ctx)
		if err != nil {
			g.Log().Warning(ctx, err)
			return
		}
		if messageId, err = s.SendForwardMessage(ctx,
			userId,
			groupId,
			s.MessageToNodes(botId, nickname, msg),
		); err == nil {
			return
		}
		g.Log().Warning(ctx, err)
	}

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			MessageType string `json:"message_type,omitempty"`
			Message     any    `json:"message"`
			AutoEscape  bool   `json:"auto_escape,omitempty"`
			UserId      int64  `json:"user_id,omitempty"`
			GroupId     int64  `json:"group_id,omitempty"`
		} `json:"params"`
	}{
		Action: "send_msg",
		Echo:   echoSign,
		Params: struct {
			MessageType string `json:"message_type,omitempty"`
			Message     any    `json:"message"`
			AutoEscape  bool   `json:"auto_escape,omitempty"`
			UserId      int64  `json:"user_id,omitempty"`
			GroupId     int64  `json:"group_id,omitempty"`
		}{
			Message:    msg,
			AutoEscape: plain,
			UserId:     userId,
			GroupId:    groupId,
		},
	}
	// message segment
	if s.isMessageSegment(ctx) {
		if plain {
			req.Params.Message = segment.NewTextSegments(msg)
			req.Params.AutoEscape = false
		} else {
			req.Params.Message = segment.ParseMessage(msg)
		}
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
		messageId = s.getMessageIdFromData(asyncCtx)
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

// SendMsg 适用于**不需要**级联撤回的场景
func (s *sBot) SendMsg(ctx context.Context, msg string, plain ...bool) {
	p := false
	if len(plain) > 0 && plain[0] {
		p = true
	}
	if _, err := s.SendMessage(ctx, s.GetUserId(ctx), s.GetGroupId(ctx), msg, p); err != nil {
		g.Log().Warning(ctx, err)
	}
}

// SendMsgIfNotApiReq 适用于**非API请求**且**需要**级联撤回的场景
func (s *sBot) SendMsgIfNotApiReq(ctx context.Context, msg string, plain ...bool) {
	if s.isApiReq(ctx) {
		return
	}
	s.SendMsgCacheContext(ctx, msg, plain...)
}

// SendMsgCacheContext 适用于**需要**级联撤回的场景
func (s *sBot) SendMsgCacheContext(ctx context.Context, msg string, plain ...bool) {
	p := false
	if len(plain) > 0 && plain[0] {
		p = true
	}
	sentMsgId, err := s.SendMessage(ctx, s.GetUserId(ctx), s.GetGroupId(ctx), msg, p)
	if err != nil {
		return
	}
	if err = s.CacheMessageContext(ctx, sentMsgId); err != nil {
		g.Log().Error(ctx, err)
	}
}

func (s *sBot) SendForwardMessage(ctx context.Context,
	userId, groupId int64,
	nodes []map[string]any,
) (messageId int64, err error) {
	// 参数校验
	if userId == 0 && groupId == 0 {
		return 0, errors.New("userId 和 groupId 不能同时为 0")
	}
	if len(nodes) == 0 {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "bot.SendForwardMessage")
	defer span.End()
	{
		messagesJSON, _ := sonic.MarshalString(nodes)
		span.SetAttributes(attribute.String("send_forward_message.messages", messagesJSON))
	}
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	if groupId != 0 {
		userId = 0
		span.SetAttributes(attribute.Int64("send_forward_message.group_id", groupId))
		go func() {
			if err := s.MarkGroupMsgAsRead(ctx, groupId); err != nil {
				g.Log().Warning(ctx, err)
			}
		}()
	} else {
		span.SetAttributes(attribute.Int64("send_forward_message.user_id", userId))
		go func() {
			if err := s.MarkPrivateMsgAsRead(ctx, userId); err != nil {
				g.Log().Warning(ctx, err)
			}
		}()
	}

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			MessageType string           `json:"message_type,omitempty"`
			UserId      int64            `json:"user_id,omitempty"`
			GroupId     int64            `json:"group_id,omitempty"`
			Messages    []map[string]any `json:"messages"`
		} `json:"params"`
	}{
		Action: "send_forward_msg",
		Echo:   echoSign,
		Params: struct {
			MessageType string           `json:"message_type,omitempty"`
			UserId      int64            `json:"user_id,omitempty"`
			GroupId     int64            `json:"group_id,omitempty"`
			Messages    []map[string]any `json:"messages"`
		}{
			UserId:   userId,
			GroupId:  groupId,
			Messages: nodes,
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
		messageId = s.getMessageIdFromData(asyncCtx)
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

func (s *sBot) SendForwardMsg(ctx context.Context, msg string) {
	botId, nickname, err := s.GetLoginInfo(ctx)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	if _, err = s.SendForwardMessage(ctx,
		s.GetUserId(ctx),
		s.GetGroupId(ctx),
		s.MessageToNodes(botId, nickname, msg),
	); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendForwardMsgCacheContext(ctx context.Context, msg string) {
	botId, nickname, err := s.GetLoginInfo(ctx)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	sentMsgId, err := s.SendForwardMessage(ctx,
		s.GetUserId(ctx),
		s.GetGroupId(ctx),
		s.MessageToNodes(botId, nickname, msg),
	)
	if err != nil {
		return
	}
	if err = s.CacheMessageContext(ctx, sentMsgId); err != nil {
		g.Log().Error(ctx, err)
	}
}

func (s *sBot) SendFileToGroup(ctx context.Context, groupId int64, filePath, name, folder string) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SendFileToGroup")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("send_file_to_group.group_id", groupId),
		attribute.String("send_file_to_group.file_path", filePath),
		attribute.String("send_file_to_group.name", name),
		attribute.String("send_file_to_group.folder", folder),
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
			GroupId int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder,omitempty"`
		} `json:"params"`
	}{
		Action: "upload_group_file",
		Echo:   echoSign,
		Params: struct {
			GroupId int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder,omitempty"`
		}{
			GroupId: groupId,
			File:    filePath,
			Name:    name,
			Folder:  folder,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) SendFileToUser(ctx context.Context, userId int64, filePath, name string) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SendFileToUser")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("send_file_to_user.user_id", userId),
		attribute.String("send_file_to_user.file_path", filePath),
		attribute.String("send_file_to_user.name", name),
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
			UserId int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		} `json:"params"`
	}{
		Action: "upload_private_file",
		Echo:   echoSign,
		Params: struct {
			UserId int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		}{
			UserId: userId,
			File:   filePath,
			Name:   name,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) SendFile(ctx context.Context, filePath, name string) (err error) {
	if groupId := s.GetGroupId(ctx); groupId != 0 {
		return s.SendFileToGroup(ctx, groupId, filePath, name, "")
	}
	return s.SendFileToUser(ctx, s.GetUserId(ctx), filePath, name)
}

func (s *sBot) UploadFile(ctx context.Context, url string) (filePath string, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.UploadFile")
	defer span.End()
	span.SetAttributes(attribute.String("upload_file.url", url))
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
			URL string `json:"url"`
		} `json:"params"`
	}{
		Action: "download_file",
		Echo:   echoSign,
		Params: struct {
			URL string `json:"url"`
		}{
			URL: url,
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
		filePath = s.getFileFromData(asyncCtx)
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) ApproveJoinGroup(ctx context.Context, flag, subType string, approve bool, reason string,
) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.ApproveJoinGroup")
	defer span.End()
	span.SetAttributes(
		attribute.String("approve_join_group.flag", flag),
		attribute.String("approve_join_group.sub_type", subType),
		attribute.Bool("approve_join_group.approve", approve),
		attribute.String("approve_join_group.reason", reason),
	)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// 参数校验
	if approve {
		reason = ""
	}
	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			Flag    string `json:"flag"`
			SubType string `json:"sub_type"`
			Approve bool   `json:"approve"`
			Reason  string `json:"reason,omitempty"`
		} `json:"params"`
	}{
		Action: "set_group_add_request",
		Echo:   echoSign,
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

func (s *sBot) SetModel(ctx context.Context, model string) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SetModel")
	defer span.End()
	span.SetAttributes(attribute.String("set_model.model", model))
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

func (s *sBot) RecallMessage(ctx context.Context, messageId int64) (err error) {
	if messageId == 0 {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "bot.RecallMessage")
	defer span.End()
	span.SetAttributes(attribute.Int64("recall_message.message_id", messageId))
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
		Action: "delete_msg",
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) Mute(ctx context.Context, groupId, userId int64, seconds int) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.Mute")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("mute.group_id", groupId),
		attribute.Int64("mute.user_id", userId),
		attribute.Int("mute.seconds", seconds),
	)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// 参数校验
	if seconds > 2591940 {
		// 不大于 29 天 23 小时 59 分钟
		// (30*24*60-1)*60=2591940 秒
		seconds = 2591940
	}
	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId  int64 `json:"group_id"`
			UserId   int64 `json:"user_id"`
			Duration int   `json:"duration"`
		} `json:"params"`
	}{
		Action: "set_group_ban",
		Echo:   echoSign,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) SetGroupCard(ctx context.Context, groupId, userId int64, card string) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SetGroupCard")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("set_group_card.group_id", groupId),
		attribute.Int64("set_group_card.user_id", userId),
		attribute.String("set_group_card.card", card),
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
			GroupId int64  `json:"group_id"`
			UserId  int64  `json:"user_id"`
			Card    string `json:"card"`
		} `json:"params"`
	}{
		Action: "set_group_card",
		Echo:   echoSign,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) Kick(ctx context.Context, groupId, userId int64, reject ...bool) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.Kick")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("kick.group_id", groupId),
		attribute.Int64("kick.user_id", userId),
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
			GroupId          int64 `json:"group_id"`
			UserId           int64 `json:"user_id"`
			RejectAddRequest bool  `json:"reject_add_request,omitempty"`
		} `json:"params"`
	}{
		Action: "set_group_kick",
		Echo:   echoSign,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) ProfileLike(ctx context.Context, userId int64, times int) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.ProfileLike")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("profile_like.user_id", userId),
		attribute.Int("profile_like.times", times),
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
			UserId int64 `json:"user_id"`
			Times  int   `json:"times"`
		} `json:"params"`
	}{
		Action: "send_like",
		Echo:   echoSign,
		Params: struct {
			UserId int64 `json:"user_id"`
			Times  int   `json:"times"`
		}{
			UserId: userId,
			Times:  times,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) EmojiLike(ctx context.Context, messageId int64, emojiId string) (err error) {
	if messageId == 0 {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "bot.EmojiLike")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("emoji_like.message_id", messageId),
		attribute.String("emoji_like.emoji_id", emojiId),
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
			MessageId int64  `json:"message_id"`
			EmojiId   string `json:"emoji_id"`
		} `json:"params"`
	}{
		Action: "set_msg_emoji_like",
		Echo:   echoSign,
		Params: struct {
			MessageId int64  `json:"message_id"`
			EmojiId   string `json:"emoji_id"`
		}{
			MessageId: messageId,
			EmojiId:   emojiId,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) Poke(ctx context.Context, groupId, userId int64) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.Poke")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("poke.group_id", groupId),
		attribute.Int64("poke.user_id", userId),
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
			GroupId int64 `json:"group_id,omitempty"`
			UserId  int64 `json:"user_id"`
		} `json:"params"`
	}{
		Action: "send_poke",
		Echo:   echoSign,
		Params: struct {
			GroupId int64 `json:"group_id,omitempty"`
			UserId  int64 `json:"user_id"`
		}{
			GroupId: groupId,
			UserId:  userId,
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
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJSON); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

func (s *sBot) Okay(ctx context.Context) (err error) {
	if groupId := s.GetGroupId(ctx); groupId != 0 {
		return s.EmojiLike(ctx, s.GetMsgId(ctx), "124") // 124: OK
	} else {
		if s.GetMsgId(ctx) == 0 {
			return
		}
		return s.Poke(ctx, groupId, s.GetUserId(ctx))
	}
}

func (s *sBot) MarkAllAsRead(ctx context.Context) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.MarkAllAsRead")
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
		Action: "_mark_all_as_read",
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

func (s *sBot) MarkPrivateMsgAsRead(ctx context.Context, userId int64) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.MarkPrivateMsgAsRead")
	defer span.End()
	span.SetAttributes(attribute.Int64("mark_private_msg_as_read.user_id", userId))
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
			UserId int64 `json:"user_id"`
		} `json:"params"`
	}{
		Action: "mark_private_msg_as_read",
		Echo:   echoSign,
		Params: struct {
			UserId int64 `json:"user_id"`
		}{
			UserId: userId,
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

func (s *sBot) MarkGroupMsgAsRead(ctx context.Context, groupId int64) (err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.MarkGroupMsgAsRead")
	defer span.End()
	span.SetAttributes(attribute.Int64("mark_group_msg_as_read.group_id", groupId))
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
		} `json:"params"`
	}{
		Action: "mark_group_msg_as_read",
		Echo:   echoSign,
		Params: struct {
			GroupId int64 `json:"group_id"`
		}{
			GroupId: groupId,
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
