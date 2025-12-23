package main

import (
	"net/http"
	"tasker/api/handler"
	"tasker/core/group"
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

	allowed := map[string]bool{
		"http://localhost:5173":         true,
		"http://8.213.208.227:15173":    true,
		"http://nullix.top":             true,
		"http://www.nullix.top":         true,
		"https://nullix.top":            true,
		"https://www.nullix.top":        true,
	}

	if allowed[origin] {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

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

	// 初始化group service
	groupRepo := db.NewGroupRepository(gormDB)
	groupSvc := group.NewService(groupRepo)

	// User相关
	userRepo := db.NewUserRepository(gormDB)
	userSvc := user.NewService(userRepo)
	userHandler := handler.NewAuthHandler(userSvc)
	userHandler.RegisterRoutes(r)

	// 初始化 Repository Service Handler
	taskRepo := db.NewTaskRepository(gormDB)
	taskSvc := task.NewService(taskRepo, groupSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	taskHandler.RegisterRoutes(r)

	// // task模块挂载（内存版）
	// taskSvc := task.NewMemoryService()
	// taskHandler := handler.NewTaskHandler(taskSvc)
	// taskHandler.RegisterRoutes(r)
	
	r.Run(":8080")
}

