package main

import (
	"net/http"
	"tasker/api/handler"
	"tasker/core/task"
	"tasker/core/user"
	"tasker/infra/db"
	"tasker/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		// 开发阶段：只允许你的前端地址
		if origin == "http://localhost:5173" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 关键：预检请求直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})


	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{
			"message": "pong",
		})
	})

	r.GET("/demo-error", func(c *gin.Context) {
		response.Error(c, http.StatusBadRequest, "DEMO_ERROR", "this is a demo error")
	})

	// 初始化数据库
	var gormDB *gorm.DB = db.NewPostgresDB()

	// User相关
	userRepo := db.NewUserRepository(gormDB)
	userSvc := user.NewService(userRepo)
	userHandler := handler.NewAuthHandler(userSvc)
	userHandler.RegisterRoutes(r)

	// 初始化 Repository Service Handler
	taskRepo := db.NewTaskRepository(gormDB)
	taskSvc := task.NewService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)
	taskHandler.RegisterRoutes(r)

	// // task模块挂载（内存版）
	// taskSvc := task.NewMemoryService()
	// taskHandler := handler.NewTaskHandler(taskSvc)
	// taskHandler.RegisterRoutes(r)
	
	r.Run(":8080")
}

