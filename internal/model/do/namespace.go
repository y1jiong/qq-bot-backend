// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Namespace is the golang structure of table namespace for DAO operations like Where/Data.
type Namespace struct {
	g.Meta      `orm:"table:namespace, do:true"`
	Namespace   any         //
	OwnerId     any         //
	SettingJson any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
