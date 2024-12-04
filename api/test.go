package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//func TestHandle() (ctx *gin.Context) {
//	successRet(ctx, map[string]interface{}{
//		"mode": "test",
//	})
//	return
//}

func TestHandle(c *gin.Context) {
	// 处理请求
	c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
}
