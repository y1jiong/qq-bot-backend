package errcode

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"net/http"
)

var (
	Forbidden       = gcode.New(http.StatusForbidden, http.StatusText(http.StatusForbidden), nil)
	Conflict        = gcode.New(http.StatusConflict, http.StatusText(http.StatusConflict), nil)
	TooEarly        = gcode.New(http.StatusTooEarly, http.StatusText(http.StatusTooEarly), nil)
	TooManyRequests = gcode.New(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests), nil)
	InternalError   = gcode.New(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
)
