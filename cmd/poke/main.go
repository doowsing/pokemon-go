package main

import (
	"log"
	"pokemon/game/app"
)

func main() {
	newApp := app.NewApp()
	defer newApp.Destory()

	// 启动
	log.Fatal(newApp.Launch())
}
