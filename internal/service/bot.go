// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IBot interface {
		CtxWithWebSocket(parent context.Context, ws *ghttp.WebSocket) context.Context
		CtxNewWebSocketMutex(parent context.Context) context.Context
		Process(ctx context.Context, rawJson []byte, nextProcess func(ctx context.Context))
		IsGroupOwnerOrAdmin(ctx context.Context) (yes bool)
		GetPostType(ctx context.Context) string
		GetMsgType(ctx context.Context) string
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
		GetGroupMemberList(ctx context.Context, groupId int64, noCache ...bool) (members []any, err error)
		GetCardOldNew(ctx context.Context) (oldCard, newCard string)
		RequestMessage(ctx context.Context, messageId int64) (messageMap map[string]any, err error)
		GetGroupInfo(ctx context.Context, groupId int64, noCache ...bool) (infoMap map[string]any, err error)
		SendMessage(ctx context.Context, messageType string, uid, gid int64, msg string, plain bool)
		SendPlainMsg(ctx context.Context, msg string)
		SendMsg(ctx context.Context, msg string)
		SendFileToGroup(ctx context.Context, gid int64, filePath, name, folder string)
		SendFileToUser(ctx context.Context, uid int64, filePath, name string)
		SendFile(ctx context.Context, name, url string)
		ApproveJoinGroup(ctx context.Context, flag, subType string, approve bool, reason string)
		SetModel(ctx context.Context, model string)
		RecallMessage(ctx context.Context, msgId int64)
		MutePrototype(ctx context.Context, groupId, userId int64, seconds int)
		Mute(ctx context.Context, seconds int)
		SetGroupCard(ctx context.Context, groupId, userId int64, card string)
		Kick(ctx context.Context, groupId, userId int64, reject ...bool)
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
