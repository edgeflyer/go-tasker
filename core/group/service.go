package group

import (
	"context"
	"tasker/pkg/apperror"
	"time"
)

type Group struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}

type service struct {
	repo Repository
}

type Service interface {
	GetGroup(ctx context.Context, userID int64, ID int64) (*Group, error)
	CreateGroup(ctx context.Context, userID int64, name string) (*Group, error)
	FindGroupByName(ctx context.Context, userID int64, name string) (*Group, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetGroup(ctx context.Context, userID int64, ID int64) (*Group, error) {
	return s.repo.GetByID(ctx, userID, ID)
}

func (s *service) CreateGroup(ctx context.Context, userID int64, name string) (*Group, error) {
	// 校验name
	if name == "" {
		return nil, apperror.New("INVALID_GROUP_NAME", "invalid group name")
	}
	
	// 组装group
	now := time.Now()
	g := &Group{
		UserID: userID,
		Name: name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

func (s *service) FindGroupByName(ctx context.Context, userID int64, name string) (*Group, error) {
	return s.repo.GetByUserIDAndName(ctx, userID, name)
}