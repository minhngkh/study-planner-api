// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTask = "task"

// Task mapped from table <task>
type Task struct {
	ID            int32     `gorm:"column:id;primaryKey" json:"id"`
	UserID        int32     `gorm:"column:user_id" json:"user_id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	Description   string    `gorm:"column:description" json:"description"`
	Priority      string    `gorm:"column:priority;not null" json:"priority"`
	EstimatedTime int32     `gorm:"column:estimated_time" json:"estimated_time"`
	Status        string    `gorm:"column:status;not null" json:"status"`
	StartTime     time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime       time.Time `gorm:"column:end_time" json:"end_time"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName Task's table name
func (*Task) TableName() string {
	return TableNameTask
}
