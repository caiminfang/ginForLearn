package app

import (
	"fmt"
	"go.uber.org/zap"
	"hello/configs"
)

var (
	Conf *configs.Config
	// Zap日志组件
	zapLogger *zap.Logger
	// 常规使用的日志组件
	Logger *zap.SugaredLogger
)

const version = "2.0.0"

func Init() {
	// 初始化配置信息
	Conf = configs.Init()

	fmt.Println("app 配置init")
}

func Stop() {
	fmt.Println("app close")
}
