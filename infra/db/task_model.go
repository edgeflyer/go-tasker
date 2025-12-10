package db
// Task 的DB模型

import (
	"time"
)

// TaskModel是存到postgres中的结构
type TaskModel struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`
	UserID int64 `gorm:"not null;index"`
	Title string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status string `gorm:"type:varchar(20);not null;index"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName可以自定义表名
func (TaskModel) TableName () string {
	return "tasks"
}