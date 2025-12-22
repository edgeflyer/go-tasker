package group

import "context"


type Repository interface {
	// 查询用户是否有这个组
	GetByID(ctx context.Context, userID, ID int64) (*Group, error)
}