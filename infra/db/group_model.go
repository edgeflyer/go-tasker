package db

import "time"

type GroupModel struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`
	// 联合唯一索引：确保同一个用户下，name不重复
	UserID int64 `gorm:"not null;index:idx_users_name,unique"`
	Name string `gorm:"type:varchar(50);not null;index:idx_users_name,unique"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (GroupModel) TableName() string {
	return "groups"
}