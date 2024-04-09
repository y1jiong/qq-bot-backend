// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// List is the golang structure for table list.
type List struct {
	ListName  string      `json:"list_name"  orm:"list_name"  ` //
	Namespace string      `json:"namespace"  orm:"namespace"  ` //
	ListJson  string      `json:"list_json"  orm:"list_json"  ` //
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" ` //
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" ` //
}
