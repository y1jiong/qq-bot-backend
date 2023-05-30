package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"qq-bot-backend/internal/service"
)

var (
	File = cFile{}
)

type cFile struct{}

func (c *cFile) GetCachedFileFromId(r *ghttp.Request) {
	ctx := r.Context()
	id := r.Get("id").String()
	content, err := service.File().GetCachedFileFromId(ctx, id)
	if err != nil {
		r.SetError(err)
		return
	}
	r.Response.Write(content)
}
