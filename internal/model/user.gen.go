// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID          int32      `gorm:"column:id;primaryKey" json:"id"`
	Email       *string    `gorm:"column:email" json:"email"`
	Password    *string    `gorm:"column:password" json:"password"`
	GoogleID    *string    `gorm:"column:google_id" json:"google_id"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsActivated bool       `gorm:"column:is_activated;not null;default:FALSE" json:"is_activated"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
