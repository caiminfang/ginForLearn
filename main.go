package main

import (
	"hello/app"
	"hello/routes"
)

func init() {
	app.Init()
}
func run() {
	//defer 最后处理，主要用于资源处理
	defer app.Stop()
	app.StartHttpServer(routes.Register)
}
func main() {
	run()
}
