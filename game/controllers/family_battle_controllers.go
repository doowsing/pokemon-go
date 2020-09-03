package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/game/ginapp"
)

func FamilyBattleInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data := gapp.OptSvc.FightSrv.FamilyBattleInfo(gapp.Id())
	gapp.JSONDATAOK("", data)
}

func FamilyBattleInvite(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	id := com.StrTo(c.Query("id")).MustInt()

	ok, msg := gapp.OptSvc.FightSrv.FamilyBattleInvite(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func FamilyBattleAccept(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	id := com.StrTo(c.Query("id")).MustInt()

	ok, msg := gapp.OptSvc.FightSrv.FamilyBattleAccept(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func FamilyBattleStartFight(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.FamilyBattleStartFight(gapp.Id())
	gapp.JSONDATAOK(msg, data)
}

func FamilyBattleAttack(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.FamilyBattleAttack(gapp.Id())
	gapp.JSONDATAOK(msg, data)
}
