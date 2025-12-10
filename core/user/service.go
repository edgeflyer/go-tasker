package user

import (
	"context"
	"time"
	"tasker/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 注册时用的输入
type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录时用的输入
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Service定义用户相关业务行为
type Service interface {
	Register(ctx context.Context, in RegisterInput) (*User, error)
	Login(ctx context.Context, in LoginInput) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

// Repository 抽象用户数据存取（后面用Postgres实现）
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, in RegisterInput) (*User, error) {
	if in.Username == "" {
		return nil, apperror.New("INVALID_USERNAME", "username is required")
	}
	if len(in.Password) < 6 {
		return nil, apperror.New("INVALID_PASSWORD", "password must be at least 6 characters")
	}

	// 检查用户名是否已存在
	existing, err := s.repo.GetByUsername(ctx, in.Username)
	if err == nil && existing != nil {
		return nil, apperror.New("USERNAME_EXISTS", "username already exists")
	}

	// bcrypt 哈希密码
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.New("INTERNAL_ERROR", "failed to hash password")
	}

	now := time.Now()
	u := &User{
		Username:  in.Username,
		Password:  string(hash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, u); err !=nil {
		return nil, err
	}

	// 返回钱把Password清掉，避免意外泄漏
	u.Password = ""
	return u, nil
}

func (s *service) Login(ctx context.Context, in LoginInput) (*User, error) {
	if in.Username == "" || in.Password == "" {
		return nil, apperror.New("INVALID_CREDENTIALS", "username and password are required")
	}

	u, err := s.repo.GetByUsername(ctx, in.Username)
	if err != nil {
		// 对外统一成账号密码错误，避免信息泄漏
		return nil, apperror.New("INVALID_CREDENTIALS", "invalid username or password")
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return nil, apperror.New("INVALID_CREDENTIALS", "invalid username or password")
	}

	// 登录成功，隐藏密码
	u.Password = ""
	return u, nil
}


func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}