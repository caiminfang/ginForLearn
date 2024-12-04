package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// successRet 成功响应，响应结果码200，有返回数据
func successRet(ctx *gin.Context, data interface{}) {
	ret(ctx, http.StatusOK, 200, data, "成功")
}

func ret(ctx *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}
