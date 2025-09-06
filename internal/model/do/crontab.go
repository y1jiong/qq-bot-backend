// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Crontab is the golang structure of table crontab for DAO operations like Where/Data.
type Crontab struct {
	g.Meta     `orm:"table:crontab, do:true"`
	Name       any         //
	Expression any         //
	CreatorId  any         //
	BotId      any         //
	Request    any         //
	CreatedAt  *gtime.Time //
}
