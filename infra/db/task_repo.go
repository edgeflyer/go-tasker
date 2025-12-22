package db

// TaskRepository的GORM实现

import (
	"context"
	"tasker/core/task"
	"tasker/pkg/apperror"

	"gorm.io/gorm"
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
		UserID:      m.UserID,
		Title:       m.Title,
		Description: m.Description,
		Status:      task.Status(m.Status),


		DueDate: m.DueData,
		Priority: m.Priority,
		GroupID: m.GroupID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toModel(t *task.Task) *TaskModel {
	return &TaskModel{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		DueData: t.DueDate,
		Priority: t.Priority,
		GroupID: t.GroupID,
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

func (r *TaskRepository) GetByID(ctx context.Context, userID, id int64) (*task.Task, error) {
	var m TaskModel
	tx := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&m)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, apperror.New("TASK_NOT_FOUND", "task not found")
		}
		return nil, apperror.New("DB_ERROR", "failed to get task")
	}
	return toDomain(&m), nil
}

func (r *TaskRepository) List(ctx context.Context, userID int64, filter task.ListTaskerFilter) (*task.ListResult, error) {
	db := r.db.WithContext(ctx).Model(&TaskModel{}).Where("user_id = ?", userID)
	if filter.Status != "" {
		db = db.Where("status = ?", string(filter.Status))
	}
	if filter.Query != "" {
		q := "%" + filter.Query + "%"
		db = db.Where("title ILIKE ? OR description ILIKE ?", q, q)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, apperror.New("DB_ERROR", "failed to count tasks")
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	var models []TaskModel
	order := "created_at DESC"
	switch filter.Sort {
	case "created_asc":
		order = "created_at ASC"
	case "status":
		order = "status ASC, created_at DESC"
	}

	if err := db.Order(order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&models).Error; err != nil {
		return nil, apperror.New("DB_ERROR", "failed to list tasks")
	}

	items := make([]*task.Task, 0, len(models))
	for _, m := range models {
		items = append(items, toDomain(&m))
	}

	return &task.ListResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	m := toModel(t)
	tx := r.db.WithContext(ctx).Model(&TaskModel{}).Where("id = ? AND user_id = ?", t.ID, t.UserID).Updates(map[string]any{
		"title":       m.Title,
		"description": m.Description,
		"status":      m.Status,
		"updated_at":  m.UpdatedAt,
	})
	if tx.Error != nil {
		return apperror.New("DB_ERROR", "failed to update task")
	}
	if tx.RowsAffected == 0 {
		return apperror.New("TASK_NOT_FOUND", "task not found")
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, userID, id int64) error {
	tx := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&TaskModel{})
	if tx.Error != nil {
		return apperror.New("DB_ERROR", "failed to delete task")
	}
	if tx.RowsAffected == 0 {
		return apperror.New("TASK_NOT_FOUND", "task not found")
	}
	return nil
}
