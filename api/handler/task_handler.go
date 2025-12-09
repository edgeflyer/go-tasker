package handler

import(
	"context"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"tasker/core/task"
	"tasker/pkg/apperror"
	"tasker/pkg/response"
)

type TaskHandler struct {
	svc task.Service
}

func NewTaskHandler(svc task.Service) * TaskHandler {
	return &TaskHandler{svc: svc}
}

// 路由注册
func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/tasks")
	{
		g.POST("", h.CreateTask)
		g.GET("", h.ListTasks)
		g.GET("/:id", h.GetTask)
		g.PUT("/:id", h.UpdateTask)
		g.DELETE("/:id", h.DeleteTask)
	}
}

// 具体handler
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var in task.CreateTaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
		return
	}

	t, err := h.svc.CreateTask(context.Background(), in)
	if err != nil {
		if appErr, ok := apperror.IsAppError(err); ok {
			// 业务错误，一般是400
			response.Error(c, http.StatusBadRequest, appErr.Code, appErr.Message)
			return
		}
		// 未知错误
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, t)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.svc.ListTasks(context.Background(), task.ListTaskerFilter{})
	if err !=nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.Success(c, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	t, err := h.svc.GetTask(context.Background(), id)
	if err != nil {
		if appErr, ok := apperror.IsAppError(err); ok && appErr.Code == "TASK_NOT_FOUND" {
			response.Error(c, http.StatusNotFound, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.Success(c, t)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var in task.UpdateTaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
		return
	}

	t, err := h.svc.UpdateTask(context.Background(), id, in)
	if err != nil {
		if appErr, ok := apperror.IsAppError(err); ok {
			switch appErr.Code {
			case "TASK_NOT_FOUND":
				response.Error(c, http.StatusNotFound, appErr.Code, appErr.Message)
			case "INVALID_TITLE", "INVALID_STATUS":
				response.Error(c, http.StatusBadRequest, appErr.Code, appErr.Message)
			default:
				response.Error(c, http.StatusBadRequest, appErr.Code, appErr.Message)
			}
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.Success(c, t)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := h.svc.DeleteTask(context.Background(), id); err != nil {
		if appErr, ok := apperror.IsAppError(err); ok && appErr.Code == "TASL_NOT_FOUND" {
			response.Error(c, http.StatusNotFound, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	// 这里可以返回204，无内容，为了统一返回一条message
	response.Success(c, gin.H{"message": "task deleted"})
}

//工具函数：解析路径参数id
func parseIDParam(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id ,err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return 0, false
	}
	return id ,true
}