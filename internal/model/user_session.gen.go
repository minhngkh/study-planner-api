// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUserSession = "user_session"

// UserSession mapped from table <user_session>
type UserSession struct {
	ID           int32     `gorm:"column:id;primaryKey" json:"id"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	RefreshToken string    `gorm:"column:refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName UserSession's table name
func (*UserSession) TableName() string {
	return TableNameUserSession
}
