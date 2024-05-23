// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package message

import (
	"context"

	"qq-bot-backend/api/message/v1"
)

type IMessageV1 interface {
	Message(ctx context.Context, req *v1.MessageReq) (res *v1.MessageRes, err error)
}
