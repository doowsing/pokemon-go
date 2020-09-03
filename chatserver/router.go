package chat

import "github.com/gin-gonic/gin"

func initRouter(r *gin.Engine) {
	r.GET("/", upgradeWs)
	r.GET("/check", checkExist)
	r.POST("/sysmsg", processSysMsg)
}

func upgradeWs(c *gin.Context) {
	serveWs(c.Writer, c.Request)
}

func processSysMsg(c *gin.Context) {
	msg := c.PostForm("data")
	ProcessSysMsg(msg)
}

func checkExist(c *gin.Context) {

}
