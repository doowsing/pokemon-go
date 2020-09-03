package controllers

import (
	"github.com/gin-gonic/gin"
	"pokemon/common/rcache"
	"pokemon/game/ginapp"
)

func CheckUnExpireProp(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.OptSvc.PropSrv.CheckPropExpire()
}

func DelZeroProp(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.OptSvc.PropSrv.CheckPropValid()
}

func EndSSBattle(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

}

func ClearSaoLei(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	rcache.ClearZbfjTimes()

	rcache.ClearSaoleiAward()
	rcache.ClearSaoleiDieUserLevel()
	rcache.ClearSaoleiTodayUser()
	rcache.ClearSaoleiTicketUser()

}
