package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"qq-bot-backend/internal/service"
)

type sMiddleware struct{}

func New() *sMiddleware {
	return &sMiddleware{}
}

func init() {
	service.RegisterMiddleware(New())
}

func (s *sMiddleware) ErrCodeToHttpStatus(r *ghttp.Request) {
	// 下一步
	r.Middleware.Next()
	// 后置中间件
	err := r.GetError()
	code := gerror.Code(err)
	if err != nil && code != gcode.CodeNil && code.Code() >= 100 && code.Code() < 600 {
		r.Response.WriteHeader(code.Code())
	}
}
