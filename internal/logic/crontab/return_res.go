package crontab

import (
	"context"
)

func (s *sCrontab) GlanceReturnRes(ctx context.Context) (retMsg string) {
	panic("implement me")
}

func (s *sCrontab) QueryReturnRes(ctx context.Context, name string) (retMsg string) {
	panic("implement me")
}

func (s *sCrontab) AddReturnRes(ctx context.Context, name, expr string, selfId int64, reqJSON []byte) (retMsg string) {
	panic("implement me")
}

func (s *sCrontab) RemoveReturnRes(ctx context.Context, name string) (retMsg string) {
	panic("implement me")
}
