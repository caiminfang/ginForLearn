package routes

import (
	"github.com/gin-gonic/gin"
	"hello/api"
)

func Register(engine *gin.Engine) {
	engine.GET("/test", api.TestHandle)
}
