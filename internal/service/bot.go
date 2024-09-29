// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/bytedance/sonic/ast"
	"github.com/gorilla/websocket"
)

type (
	IBot interface {
		CtxWithWebSocket(parent context.Context, conn *websocket.Conn) context.Context
		CtxNewWebSocketMutex(parent context.Context) context.Context
		CtxWithReqJson(ctx context.Context, reqJson *ast.Node) context.Context
		Process(ctx context.Context, rawJson []byte, nextProcess func(ctx context.Context))
		JoinConnectionPool(ctx context.Context, key int64)
		LeaveConnectionPool(key int64)
		LoadConnectionPool(key int64) context.Context
		CacheMessageContext(ctx context.Context, userId int64, lastMessageId int64, currentMessageId int64) error
		GetCachedMessageContext(ctx context.Context, userId int64, lastMessageId int64) (currentMessageIds []int64, exist bool, err error)
		Forward(ctx context.Context, url string, key string) (err error)
		GetPostType(ctx context.Context) string
		GetMsgType(ctx context.Context) string
		GuessMsgType(groupId int64) string
		GetRequestType(ctx context.Context) string
		GetNoticeType(ctx context.Context) string
		GetSubType(ctx context.Context) string
		GetMsgId(ctx context.Context) int64
		GetMessage(ctx context.Context) string
		GetUserId(ctx context.Context) int64
		GetGroupId(ctx context.Context) int64
		GetComment(ctx context.Context) string
		GetFlag(ctx context.Context) string
		GetTimestamp(ctx context.Context) int64
		GetOperatorId(ctx context.Context) int64
		GetSelfId(ctx context.Context) int64
		GetNickname(ctx context.Context) string
		GetCard(ctx context.Context) string
		GetCardOrNickname(ctx context.Context) string
		GetCardOldNew(ctx context.Context) (oldCard string, newCard string)
		GetGroupMemberInfo(ctx context.Context, groupId int64, userId int64, noCache ...bool) (member *ast.Node, err error)
		GetGroupMemberList(ctx context.Context, groupId int64, noCache ...bool) (members []any, err error)
		RequestMessage(ctx context.Context, messageId int64) (messageMap map[string]any, err error)
		GetGroupInfo(ctx context.Context, groupId int64, noCache ...bool) (infoMap map[string]any, err error)
		GetLoginInfo(ctx context.Context) (userId int64, nickname string)
		IsGroupOwnerOrAdmin(ctx context.Context) bool
		IsGroupOwnerOrAdminOrSysTrusted(ctx context.Context) bool
		SendMessage(ctx context.Context, messageType string, userId int64, groupId int64, msg string, plain bool) (messageId int64, err error)
		// SendPlainMsg 适用于*事件*触发的消息发送
		SendPlainMsg(ctx context.Context, msg string)
		// SendMsg 适用于*事件*触发的消息发送
		SendMsg(ctx context.Context, msg string)
		SendMsgIfNotApiReq(ctx context.Context, msg string, richText ...bool)
		// SendMsgCacheContext 适用于*非事件*触发的消息发送
		SendMsgCacheContext(ctx context.Context, msg string, richText ...bool)
		SendFileToGroup(ctx context.Context, groupId int64, filePath string, name string, folder string)
		SendFileToUser(ctx context.Context, userId int64, filePath string, name string)
		SendFile(ctx context.Context, filePath string, name string)
		UploadFile(ctx context.Context, url string) (filePath string, err error)
		ApproveJoinGroup(ctx context.Context, flag string, subType string, approve bool, reason string)
		SetModel(ctx context.Context, model string)
		RecallMessage(ctx context.Context, msgId int64)
		MutePrototype(ctx context.Context, groupId int64, userId int64, seconds int)
		Mute(ctx context.Context, seconds int)
		SetGroupCard(ctx context.Context, groupId int64, userId int64, card string)
		Kick(ctx context.Context, groupId int64, userId int64, reject ...bool)
		RewriteMessage(ctx context.Context, message string)
		SetHistory(ctx context.Context, history string) error
	}
)

var (
	localBot IBot
)

func Bot() IBot {
	if localBot == nil {
		panic("implement not found for interface IBot, forgot register?")
	}
	return localBot
}

func RegisterBot(i IBot) {
	localBot = i
}
