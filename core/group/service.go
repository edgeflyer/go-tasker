package group

import (
	"context"
	"time"
)

type Group struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type service struct {
	repo Repository
}

type Service interface {
	GetGroup(ctx context.Context, userID int64, ID int64) (*Group, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetGroup(ctx context.Context, userID int64, ID int64) (*Group, error) {
	return s.repo.GetByID(ctx, userID, ID)
}
