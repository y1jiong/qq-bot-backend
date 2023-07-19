// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// List is the golang structure of table list for DAO operations like Where/Data.
type List struct {
	g.Meta    `orm:"table:list, do:true"`
	ListName  interface{} //
	Namespace interface{} //
	ListJson  interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
