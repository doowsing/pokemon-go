package main

import (
	"log"
	"pokemon/pkg/app"
)

func main() {
	app := app.NewApp()
	defer app.Destory()

	// 启动
	log.Fatal(app.Launch())
}
