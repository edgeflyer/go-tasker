package db

import (
	"context"
	"gorm.io/gorm"
	"tasker/core/user"
	"tasker/pkg/apperror"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func userToDomain(m *UserModel) *user.User {
	return &user.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password, // 这里仍然是哈希，Service 会决定是否清掉
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func userToModel(u *user.User) *UserModel {
	return &UserModel{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// 实现user.Repository接口

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	m :=userToModel(u)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		// 这里可以更细化处理唯一约束错误，先简单用一个统一的DB_ERROR
		return apperror.New("DB_ERROR", "failed to create user")
	}
	u.ID = m.ID
	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var m UserModel
	tx := r.db.WithContext(ctx).Where("username = ?", username).First(&m)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, apperror.New("USER_NOT_FOUND", "user not found")
		}
		return nil, apperror.New("DB_ERROR", "failed to get user")
	}
	return userToDomain(&m), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	var m UserModel
	tx := r.db.WithContext(ctx).First(&m, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, apperror.New("USER_NOT_FOUND", "user not found")
		}
		return nil, apperror.New("DB_ERROR", "failed to get user by id")
	}
	return userToDomain(&m), nil
}