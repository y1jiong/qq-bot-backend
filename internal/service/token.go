// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IToken interface {
		IsCorrectToken(ctx context.Context, token string) (correct bool, name string, ownerId, botId int64)
		UpdateLoginTime(ctx context.Context, token string)
		AddNewTokenReturnRes(ctx context.Context, name, token string) (retMsg string)
		RemoveTokenReturnRes(ctx context.Context, name string) (retMsg string)
		QueryTokenReturnRes(ctx context.Context, name string) (retMsg string)
		QueryOwnTokenReturnRes(ctx context.Context) (retMsg string)
		ChangeTokenOwnerReturnRes(ctx context.Context, name, ownerId string) (retMsg string)
		BindTokenBotId(ctx context.Context, name, botId string) (retMsg string)
		UnbindTokenBotId(ctx context.Context, name string) (retMsg string)
	}
)

var (
	localToken IToken
)

func Token() IToken {
	if localToken == nil {
		panic("implement not found for interface IToken, forgot register?")
	}
	return localToken
}

func RegisterToken(i IToken) {
	localToken = i
}
