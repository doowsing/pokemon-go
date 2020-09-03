package chat

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	flag.Parse()
	log.SetFlags(0)
	app := gin.Default()
	initRouter(app)

	// 启动
	log.Fatal(app.Run(":2020"))
}
