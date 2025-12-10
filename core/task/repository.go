package task

import (
	"context"
)

// Repository抽象了对task的 持久化操作
type Repository interface {
	Create(ctx context.Context, t *Task) error
	GetByID(ctx context.Context, userID, id int64) (*Task, error)
	List(ctx context.Context, userID int64) ([]*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, userID, id int64) error
}