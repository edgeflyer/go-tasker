package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"tasker/pkg/jwtutil"
	"tasker/pkg/response"
)

// AuthMiddleware 验证JWT， 成功的话把userID写进gin.Context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing Authorization header")
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid Authorization header format")
			c.Abort()
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))

		claims, err := jwtutil.ParseToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			c.Abort()
			return
		}

		// 把userID放进context，后面的handler可以取出来用
		c.Set("userID", claims.UserID)

		c.Next()
	}
}