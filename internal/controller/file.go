package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

var (
	File = cFile{}
)

type cFile struct{}

func (c *cFile) GetCachedFileById(r *ghttp.Request) {
	ctx := r.Context()
	ctx, span := gtrace.NewSpan(ctx, "controller.file.GetCachedFileById")
	defer span.End()
	defer func() {
		if err := r.GetError(); err != nil {
			span.RecordError(err)
		}
	}()

	id := r.Get("id").String()
	content, err := service.File().GetCachedFileById(ctx, id)
	if err != nil {
		r.SetError(err)
		return
	}
	r.Response.Write(content)
}
