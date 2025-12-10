package main

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tasker/pkg/response"
	"tasker/core/task"
	"tasker/api/handler"
	"tasker/infra/db"
	"tasker/core/user"
	
)

func main() {
	r := gin.Default()

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

