package errcode

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"net/http"
)

var (
	MessageExpired   = gcode.New(http.StatusBadRequest, "message expired", nil)
	SignatureError   = gcode.New(http.StatusBadRequest, "signature error", nil)
	CommandNotFound  = gcode.New(http.StatusBadRequest, "command not found", nil)
	GroupNotBinding  = gcode.New(http.StatusBadRequest, "group not binding", nil)
	PermissionDenied = gcode.New(http.StatusForbidden, "permission denied", nil)
	FileNotFound     = gcode.New(http.StatusNotFound, "file not found", nil)
	BotNotConnected  = gcode.New(http.StatusInternalServerError, "bot not connected", nil)
)
