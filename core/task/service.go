package task

import(
	"context"
	"sync"
	"time"
	"tasker/pkg/apperror"
)

type Status string

const(
	StatusPending Status = "pending"
	StatusCompleted Status = "completed"
)

type Task struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 创建任务时用的入参
type CreateTaskInput struct {
	Title string `json:"title"`
	Description string `json:"description"`
}

// 更新任务时用的入参（目前设置的比填)
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type ListTaskerFilter struct {
	// 后面可以加status / 分页等
}

// Service把task相关业务抽象出来
type Service interface {
	CreateTask(ctx context.Context, in CreateTaskInput) (*Task, error)
	GetTask(ctx context.Context, id int64) (*Task, error)
	ListTasks(ctx context.Context, filter ListTaskerFilter) ([]*Task, error)
	UpdateTask(ctx context.Context, id int64, in UpdateTaskInput) (*Task, error)
	DeleteTask(ctx context.Context, id int64) error
}

// 内存实现版(用map存任务，用锁保证并发安全)
type memoryService struct {
	mu sync.RWMutex
	seq int64
	tasks map[int64]*Task
}

func NewMemoryService() Service {
	return &memoryService{
		tasks: make(map[int64]*Task),
	}
}

// 工具函数
func validateStatus(status Status) bool {
	return status == StatusPending || status == StatusCompleted
}

// Service实现
func (s *memoryService) CreateTask(ctx context.Context, in CreateTaskInput) (*Task, error) {
	if in.Title == "" {
		return nil, apperror.New("INVALID_TITLE", "title is required")
	}

	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	task := &Task{
		ID: s.seq,
		Title: in.Title,
		Description: in.Description,
		Status: StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.tasks[task.ID] = task
	return task, nil
}

func (s *memoryService) GetTask(ctx context.Context,id int64) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, apperror.New("TASK_NOT_FOUND", "task not found")
	}
	return t, nil
}

func (s *memoryService) ListTasks(ctx context.Context, filter ListTaskerFilter) ([]*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (s *memoryService) UpdateTask(ctx context.Context, id int64, in UpdateTaskInput) (*Task, error) {
	if in.Title == "" {
		return nil, apperror.New("INVALID_TITLE", "title is reuired")
	}
	if !validateStatus(in.Status) {
		return nil, apperror.New("INVALID_STATUS", "status must be 'pending' or 'completed'")
	}


	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, apperror.New("TASK_NOT_FOUND", "task not found")
	}

	t.Title = in.Title
	t.Description = in.Description
	t.Status = in.Status
	t.UpdatedAt = time.Now()

	return t, nil
}

func (s *memoryService) DeleteTask(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return apperror.New("TASK_NOT_FOUD", "task not found")
	}
	delete(s.tasks, id)
	return nil
}