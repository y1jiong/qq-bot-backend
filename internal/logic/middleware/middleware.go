package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"net/http"
	"qq-bot-backend/internal/service"
	"time"
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

func (s *sMiddleware) RateLimit(r *ghttp.Request) {
	// 前置中间件
	cacheKey := "RateLimit" + r.GetRemoteIp()
	limitTimes := 2
	intervalTime := time.Second
	// Rate Limit
	exist := false
	if v, err := gcache.Get(r.Context(), cacheKey); v != nil {
		if err != nil {
			r.SetError(err)
			return
		}
		if v.Int() >= limitTimes {
			r.Response.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, exist, err = gcache.Update(r.Context(), cacheKey, v.Int()+1)
		if err != nil {
			r.SetError(err)
			return
		}
	}
	if !exist {
		err := gcache.Set(r.Context(), cacheKey, 1, intervalTime)
		if err != nil {
			r.SetError(err)
			return
		}
	}
	// 下一步
	r.Middleware.Next()
}
