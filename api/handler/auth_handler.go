package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"tasker/core/user"
	"tasker/pkg/apperror"
	"tasker/pkg/jwtutil"
	"tasker/pkg/response"
)

type AuthHandler struct {
	userSvc user.Service
}

func NewAuthHandler(userSvc user.Service) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// 注册路由
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/auth")
	{
		g.POST("/register", h.Register)
		g.POST("/login", h.Login)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var in user.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
		return
	}

	u ,err := h.userSvc.Register(context.Background(), in)
	if err != nil {
		if appErr, ok := apperror.IsAppError(err); ok {
			switch appErr.Code {
			case "INVALID_USERNAME", "INVALID_PASSWORD":
				response.Error(c, http.StatusConflict, appErr.Code, appErr.Message)
			case "USERNAME_EXISTS":
				response.Error(c, http.StatusConflict, appErr.Code, appErr.Message)
			default:
				response.Error(c, http.StatusBadRequest, appErr.Code, appErr.Message)
			}
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	response.SuccessWithStatus(c, http.StatusCreated, gin.H{
		"user": u,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in user.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
		return
	}

	u, err := h.userSvc.Login(context.Background(), in)
	if err != nil {
		if appErr, ok := apperror.IsAppError(err); ok {
			// 登录失败统一认为是401
			if appErr.Code == "INVALID_CREDENTIALS" {
				response.Error(c, http.StatusUnauthorized, appErr.Code, appErr.Message)
				return
			}
			response.Error(c, http.StatusBadRequest, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	// 登录成功
	token, err := jwtutil.GenerateToken(u.ID, 2*time.Hour)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TOKEN_ERROR", "failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": u,
	})
}