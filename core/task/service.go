package task

import (
	"context"
	// "sync"
	"tasker/pkg/apperror"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"` // 在写完auth之后新增：任务属于哪个用户
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`

	DueDate *time.Time `json:"due_date"`
	Priority string `json:"priority"`
	GroupID *int64 `json:"group_id"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Group struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 创建任务时用的入参
type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate *time.Time `json:"due_date"`
	Priority string `json:"priority"`
	GroupID *int64 `json:"group_id"`
}

// 更新任务时用的入参（目前设置的必填)
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type ListTaskerFilter struct {
	Status   Status `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Query    string `json:"query"`
	Sort     string `json:"sort"` // created_desc/created_asc/status
}

type ListResult struct {
	Items    []*Task `json:"items"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

// Service把task相关业务抽象出来
type Service interface {
	CreateTask(ctx context.Context, userID int64, in CreateTaskInput) (*Task, error)
	GetTask(ctx context.Context, userID int64, id int64) (*Task, error)
	ListTasks(ctx context.Context, userID int64, filter ListTaskerFilter) (*ListResult, error)
	UpdateTask(ctx context.Context, userID int64, id int64, in UpdateTaskInput) (*Task, error)
	DeleteTask(ctx context.Context, userID int64, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// 实现Service方法
func (s *service) CreateTask(ctx context.Context, userID int64, in CreateTaskInput) (*Task, error) {
	if in.Title == "" {
		return nil, apperror.New("INVALID_TITLE", "title is required")
	}

	if in.Priority == "" {
		in.Priority = "low"
	}

	// 确认分组属于用户
	exist, err := s.repo.GetGroupByID(ctx, userID)

	now := time.Now()
	t := &Task{
		UserID:      userID,
		Title:       in.Title,
		Description: in.Description,
		Status:      StatusPending,
		DueDate: in.DueDate,
		Priority: in.Priority,
		GroupID: in.GroupID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *service) GetTask(ctx context.Context, userID int64, id int64) (*Task, error) {
	return s.repo.GetByID(ctx, userID, id)
}

func (s *service) ListTasks(ctx context.Context, userID int64, filter ListTaskerFilter) (*ListResult, error) {
	// normalize
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if filter.Status != "" && filter.Status != StatusPending && filter.Status != StatusCompleted {
		return nil, apperror.New("INVALID_STATUS", "status must be 'pending' or 'completed'")
	}

	// sort allowlist
	switch filter.Sort {
	case "", "created_desc", "created_asc", "status":
	default:
		return nil, apperror.New("INVALID_SORT", "sort must be created_desc/created_asc/status")
	}

	filter.Page = page
	filter.PageSize = pageSize

	return s.repo.List(ctx, userID, filter)
}

func (s *service) UpdateTask(ctx context.Context, userID int64, id int64, in UpdateTaskInput) (*Task, error) {
	if in.Title == "" {
		return nil, apperror.New("INVALID_TITLE", "title is required")
	}
	if in.Status != StatusPending && in.Status != StatusCompleted {
		return nil, apperror.New("INVALID_STATUS", "status must be 'pending' or 'completed'")
	}

	t, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	t.Title = in.Title
	t.Description = in.Description
	t.Status = in.Status
	t.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *service) DeleteTask(ctx context.Context, userID int64, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}
