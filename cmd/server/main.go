package main

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"tasker/pkg/response"
	"tasker/core/task"
	"tasker/api/handler"
	
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

	// task模块挂载
	taskSvc := task.NewMemoryService()
	taskHandler := handler.NewTaskHandler(taskSvc)
	taskHandler.RegisterRoutes(r)
	
	r.Run(":8080")
}

