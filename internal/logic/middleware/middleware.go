package middleware

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"net/http"
	"qq-bot-backend/internal/model"
	"qq-bot-backend/internal/service"
	"time"
)

var (
	middlewareAccessIntervalMilliseconds *gvar.Var
)

type sMiddleware struct{}

func init() {
	service.RegisterMiddleware(New())
}

func New() *sMiddleware {
	return &sMiddleware{}
}

// Common 共用后置中间件，序列化 json 响应
func (s *sMiddleware) Common(r *ghttp.Request) {
	r.Middleware.Next()
	// 后置中间件
	// 设置 cdn 不缓存
	r.Response.Header().Set("Cache-Control", "no-cache")
	// 错误处理
	err := r.GetError()
	if err != nil {
		code := gerror.Code(err)
		msg := ""
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
			r.Response.WriteHeader(http.StatusInternalServerError)
		}
		msg = err.Error()
		r.Response.WriteJson(model.CommonResPrefix{
			Code:    code.Code(),
			Message: msg,
		})
	} else {
		// 获取响应内容
		res := r.GetHandlerResponse()
		if res != nil {
			r.Response.WriteJson(res)
			return
		}
	}
}

// AccessIntervalControl 访问间隔控制
func (s *sMiddleware) AccessIntervalControl(r *ghttp.Request) {
	// 前置中间件
	if v, err := gcache.Get(r.Context(), r.Session.MustId()); v != nil {
		if err != nil {
			return
		}
		r.Response.WriteHeader(http.StatusTooManyRequests)
		return
	}
	err := gcache.Set(r.Context(), r.Session.MustId(), true,
		service.Cfg().GetMiddlewareAccessIntervalMilliseconds(r.Context())*time.Millisecond)
	if err != nil {
		return
	}
	// 下一步
	r.Middleware.Next()
}
