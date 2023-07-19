// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Namespace is the golang structure for table namespace.
type Namespace struct {
	Namespace   string      `json:"namespace"    ` //
	OwnerId     int64       `json:"owner_id"     ` //
	SettingJson string      `json:"setting_json" ` //
	CreatedAt   *gtime.Time `json:"created_at"   ` //
	UpdatedAt   *gtime.Time `json:"updated_at"   ` //
	DeletedAt   *gtime.Time `json:"deleted_at"   ` //
}
