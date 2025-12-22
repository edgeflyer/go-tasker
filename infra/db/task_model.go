package db

// Task 的DB模型

import (
	"time"
)

// TaskModel是存到postgres中的结构
type TaskModel struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`
	// 关联用户（集联删除）
	UserID int64 `gorm:"not null;index"`
	Title string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status string `gorm:"type:varchar(20);not null;index"`

	DueData *time.Time `gorm:"not null"`
	Priority string `gorm:"type:varchar(20);default:'low';index"`
	GroupID *int64 `gorm:"index"`

	// OnDelete:SET NULL意思是如果这个组被删除了，这些人物的GroupID自动变成NULL
	Group GroupModel `gorm:"foreignKey:GroupID;constraint:GroupID;constraint:OnDelete:SET NULL"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName可以自定义表名
func (TaskModel) TableName () string {
	return "tasks"
}