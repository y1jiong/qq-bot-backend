// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	UserId      int64       `json:"user_id"      orm:"user_id"      ` //
	SettingJson string      `json:"setting_json" orm:"setting_json" ` //
	CreatedAt   *gtime.Time `json:"created_at"   orm:"created_at"   ` //
	UpdatedAt   *gtime.Time `json:"updated_at"   orm:"updated_at"   ` //
}
