// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ICfg interface {
		IsDebugEnabled(ctx context.Context) bool
		GetDebugToken(ctx context.Context) string
		GetOTelConfig(ctx context.Context) (endpoint string, traceToken string)
	}
)

var (
	localCfg ICfg
)

func Cfg() ICfg {
	if localCfg == nil {
		panic("implement not found for interface ICfg, forgot register?")
	}
	return localCfg
}

func RegisterCfg(i ICfg) {
	localCfg = i
}
