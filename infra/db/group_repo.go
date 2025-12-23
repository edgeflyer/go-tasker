package db

import (
	"context"
	"tasker/core/group"
	"tasker/pkg/apperror"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func groupToDomain(m *GroupModel) *group.Group {
	return &group.Group{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func groupToModel(t *group.Group) *GroupModel {
	return &GroupModel{
		ID:        t.ID,
		UserID:    t.UserID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (r *GroupRepository) GetByID(ctx context.Context, userID int64, ID int64) (*group.Group, error) {
	var m GroupModel
	tx := r.db.WithContext(ctx).Where("user_id = ? and id = ?", userID, ID).First(&m)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, apperror.New("GROUP_NOT_FOUND", "group not found")
		}
		return nil, apperror.New("DB_ERROR", "failed to get group")
	}
	return groupToDomain(&m), nil
}

func (r *GroupRepository) Create(ctx context.Context, group *group.Group) error {
	m := groupToModel(group)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return apperror.New("DB_ERROR", "failed to create group")
	}

	group.ID = m.ID
	return nil
}

func (r *GroupRepository) GetByUserIDAndName(ctx context.Context, userID int64, name string) (*group.Group, error) {
	var m GroupModel
	tx := r.db.WithContext(ctx).Where("user_id = ? and name = ?", userID, name).First(&m)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, apperror.New("DB_ERROR", "failed to get group")
	}
	return groupToDomain(&m), nil
}