package v1

import "github.com/gogf/gf/v2/frame/g"

type CommandReq struct {
	g.Meta    `path:"/command" method:"post" tags:"api" summary:"命令"`
	Token     string `json:"token" v:"required" description:"必填"`
	Command   string `json:"command" v:"required" description:"必填"`
	GroupId   int64  `json:"group_id"`
	Timestamp int64  `json:"timestamp" v:"required" description:"单位：秒；超过 5 秒的请求会被拒绝"`
	Signature string `json:"signature" v:"required" description:"tokenName+token+timestamp+command 的 sha1 值的 base64 值"`
}
type CommandRes struct {
	Message string `json:"message"`
}
