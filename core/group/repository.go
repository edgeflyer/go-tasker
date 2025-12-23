package group

import "context"

type Repository interface {
	// 查询用户是否有这个组
	GetByID(ctx context.Context, userID int64, ID int64) (*Group, error)

	// 创建分组
	Create(ctx context.Context, group *Group) error

	GetByUserIDAndName(ctx context.Context, userID int64, name string) (*Group, error)

	Delete(ctx context.Context, userID, ID uint64) error

	Update(ctx context.Context, group *Group) (*Group, error)

	GetListByName(ctx context.Context, userID int64, name string) (*[]Group, error)

	GetListByUserID(ctx context.Context, userID int64) (*[]Group, error)
}
