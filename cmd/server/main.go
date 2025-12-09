package main

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"tasker/pkg/response"
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

	r.Run(":8080")
}

