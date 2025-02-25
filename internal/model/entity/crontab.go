// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Crontab is the golang structure for table crontab.
type Crontab struct {
	Name       string      `json:"name"       orm:"name"       ` //
	Expression string      `json:"expression" orm:"expression" ` //
	BotId      int64       `json:"bot_id"     orm:"bot_id"     ` //
	Request    string      `json:"request"    orm:"request"    ` //
	CreatedAt  *gtime.Time `json:"created_at" orm:"created_at" ` //
}
