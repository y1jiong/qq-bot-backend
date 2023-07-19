// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Token is the golang structure for table token.
type Token struct {
	Name        string      `json:"name"          ` //
	Token       string      `json:"token"         ` //
	OwnerId     int64       `json:"owner_id"      ` //
	CreatedAt   *gtime.Time `json:"created_at"    ` //
	UpdatedAt   *gtime.Time `json:"updated_at"    ` //
	DeletedAt   *gtime.Time `json:"deleted_at"    ` //
	LastLoginAt *gtime.Time `json:"last_login_at" ` //
}
