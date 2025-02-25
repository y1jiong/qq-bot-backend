package crontab

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"qq-bot-backend/internal/service"
	"strings"
)

type sCrontab struct{}

func New() *sCrontab {
	return &sCrontab{}
}

func init() {
	service.RegisterCrontab(New())
}

func (s *sCrontab) Run(ctx context.Context) {

}

func (s *sCrontab) add(ctx context.Context, name, expr string, botId int64, reqJSON []byte) (err error) {
	if len(strings.Fields(expr)) == 5 {
		expr = "# " + expr
	}
	_, err = gcron.AddSingleton(ctx, expr, func(ctx context.Context) {
		botCtx := service.Bot().LoadConnection(botId)
		if botCtx == nil {
			return
		}
		service.Bot().Process(botCtx, reqJSON, service.Process().Process)
	}, name)
	return
}

func (s *sCrontab) remove(name string) {
	gcron.Remove(name)
}
