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
	r.Middleware.Next()

	if err := r.GetError(); err != nil {
		if code := gerror.Code(err); code != gcode.CodeNil && code.Code() >= 100 && code.Code() < 600 {
			r.Response.WriteHeader(code.Code())
		}
	}
}

func (s *sMiddleware) RateLimit(r *ghttp.Request) {
	cacheKey := "RateLimit_" + r.GetRemoteIp()
	const limitTimes = 2
	// Rate Limit
	if v, err := gcache.GetOrSet(r.Context(), cacheKey, 0, time.Second); err == nil {
		times := v.Int()
		if times >= limitTimes {
			r.Response.WriteHeader(http.StatusTooManyRequests)
			return
		}
		if _, _, err = gcache.Update(r.Context(), cacheKey, times+1); err != nil {
			r.SetError(err)
			return
		}
	} else {
		r.SetError(err)
		return
	}

	r.Middleware.Next()
}
