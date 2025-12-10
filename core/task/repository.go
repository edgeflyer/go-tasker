package task

import (
	"context"
)

// Repository抽象了对task的 持久化操作
type Respository interface {
	Create(ctx context.Context, t *Task) error
	GetByID(ctx context.Context, id int64) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id int64) error
}