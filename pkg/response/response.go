package response

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

// 响应成功统一结构: {"data":...}
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// 自定义状态码的成功响应
func SuccessWithStatus(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"data": data,
	})
}

/*错误响应统一结构:
{
	"error": 
		{
			"code: "...",
			"message": "..."
		}
}
*/
type ErrorBody struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

func Error(c *gin.Context, status int, code, msg string) {
	c.JSON(status, gin.H{
		"error": ErrorBody {
			Code: code,
			Message: msg,
		},
	})
}