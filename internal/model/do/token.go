// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Token is the golang structure of table token for DAO operations like Where/Data.
type Token struct {
	g.Meta       `orm:"table:token, do:true"`
	Name         interface{} //
	Token        interface{} //
	OwnerId      interface{} //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	DeletedAt    *gtime.Time //
	LastLoginAt  *gtime.Time //
	BindingBotId interface{} //
}
