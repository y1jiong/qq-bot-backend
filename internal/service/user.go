// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IUser interface {
		QueryUserReturnRes(ctx context.Context, userId int64) (retMsg string)
		SystemTrustUserReturnRes(ctx context.Context, userId int64) (retMsg string)
		SystemDistrustUserReturnRes(ctx context.Context, userId int64) (retMsg string)
		GrantOpTokenReturnRes(ctx context.Context, userId int64) (retMsg string)
		RevokeOpTokenReturnRes(ctx context.Context, userId int64) (retMsg string)
		GrantOpNamespaceReturnRes(ctx context.Context, userId int64) (retMsg string)
		RevokeOpNamespaceReturnRes(ctx context.Context, userId int64) (retMsg string)
		GrantGetRawMsgReturnRes(ctx context.Context, userId int64) (retMsg string)
		RevokeGetRawMsgReturnRes(ctx context.Context, userId int64) (retMsg string)
		GrantRecallReturnRes(ctx context.Context, userId int64) (retMsg string)
		RevokeRecallReturnRes(ctx context.Context, userId int64) (retMsg string)
		IsSystemTrustedUser(ctx context.Context, userId int64) bool
		CanOpToken(ctx context.Context, userId int64) bool
		CanOpNamespace(ctx context.Context, userId int64) bool
		CanGetRawMessage(ctx context.Context, userId int64) bool
		CanRecallMessage(ctx context.Context, userId int64) bool
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
