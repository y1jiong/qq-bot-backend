// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package command

import (
	"context"

	"qq-bot-backend/api/command/v1"
)

type ICommandV1 interface {
	Command(ctx context.Context, req *v1.CommandReq) (res *v1.CommandRes, err error)
}
