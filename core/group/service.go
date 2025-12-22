package group

import (
	"context"
	"time"
)

type Group struct {
	ID int64
	UserID int64
	Name string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type service struct {
	repo Repository
}

type Service interface{
	GetByID(ctx context.Context, userID int64)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}