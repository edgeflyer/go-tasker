package db

// TaskRepository的GORM实现

import (
	"context"
	"gorm.io/gorm"
	"tasker/core/task"
	"tasker/pkg/apperror"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}


// model和domain转换
func toDomain(m *TaskModel) *task.Task {
	return &task.Task{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      task.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toModel(t *task.Task) *TaskModel {
	return &TaskModel{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// 实现Repository接口
func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	m := toModel(t)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return apperror.New("DB_ERROR", "failed to create task")
	}
	// 回填自增ID
	t.ID = m.ID
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*task.Task, error) {
	var m TaskModel
	tx := r.db.WithContext(ctx).First(&m, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, apperror.New("TASK_NOT_FOUND", "task not found")
		}
		return nil, apperror.New("DB_ERROR", "failed to get task")
	}
	return toDomain(&m), nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*task.Task, error) {
	var models []TaskModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, apperror.New("DB_ERROR", "failed to list tasks")
	}

	result := make([]*task.Task, 0, len(models))
	for _, m := range models {
		t := toDomain(&m)
		result = append(result, t)
	}
	return result, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	m := toModel(t)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return apperror.New("DB_ERROR", "failed to update task")
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx).Delete(&TaskModel{}, id)
	if tx.Error != nil {
		return apperror.New("DB_ERROR", "failed to delete task")
	}
	if tx.RowsAffected == 0 {
		return apperror.New("TASK_NOT_FOUND", "task not found")
	}
	return nil
}