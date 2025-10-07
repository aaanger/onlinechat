package response

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, message)
	return
}
