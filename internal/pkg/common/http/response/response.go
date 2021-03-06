package response

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, gin.H{
		"success":    false,
		"statusCode": httpCode,
		"message":    err.Error(),
	})
}

func ErrorWithMessage(c *gin.Context, httpCode int, msg string, err error) {
	c.JSON(httpCode, gin.H{
		"success":    false,
		"statusCode": httpCode,
		"message":    msg,
		"error":      err.Error(),
	})
}

func Success(c *gin.Context, httpCode int, data interface{}) {
	c.JSON(httpCode, gin.H{
		"success":    true,
		"statusCode": httpCode,
		"message":    "success",
		"data":       data,
	})
}

func SuccessWithMessage(c *gin.Context, httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, gin.H{
		"success":    true,
		"statusCode": httpCode,
		"message":    msg,
		"data":       data,
	})
}
