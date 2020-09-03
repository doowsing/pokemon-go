package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/game/ginapp"
)

// 给聊天服务器调用的http接口
func ChatLogin(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	info := gapp.OptSvc.SysSrv.GetChatUserInfo(gapp.Id())
	gapp.JSONDATAOK("", info)
}

func SetGroupUnReady(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := c.Query("id")
	uuid := c.Query("uuid")
	gapp.OptSvc.FightSrv.SetUserUnReady(com.StrTo(id).MustInt(), uuid)

}
