package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/game/ginapp"
	"pokemon/game/services/common"
	"strconv"
)

// 神圣战场玩家列表以及战场现况
func SSBattleUserList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data := gapp.OptSvc.FightSrv.GetSSBattleUserList()
	introduce := common.GetWelcome("battle")
	if introduce != nil {
		data["introduce"] = introduce.Content
	} else {
		data["introduce"] = "未更新"
	}
	gapp.JSONDATAOK("", data)
}

// 进入神圣战场
func SSBattleEnter(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	msg, data := gapp.OptSvc.FightSrv.SSBattleEnter(gapp.Id(), c.Query("faction") == "1")
	gapp.JSONDATAOK(msg, data)
}

// 进入神圣战场战斗
func SSBattleStartFight(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.SSBattleStartFight(gapp.Id(), c.Query("level"))
	gapp.JSONDATAOK(msg, data)
}

// 进入神圣战场攻击
func SSBattleAttack(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.SSBattleAttack(gapp.Id())
	gapp.JSONDATAOK(msg, data)
}

// 神圣战场使用道具
func SSBattleUseProp(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := com.StrTo(c.Query("id")).MustInt()

	data, msg := gapp.OptSvc.FightSrv.SSBattleUseProp(gapp.Id(), id)
	gapp.JSONDATAOK(msg, data)
}

// 神圣战场用户信息
func SSBattleUserInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	battleUser := gapp.OptSvc.FightSrv.GetSSBattleUser(gapp.Id())
	if battleUser != nil {
		rank := "未上榜"
		if battleUser.Tops > 0 {
			rank = strconv.Itoa(battleUser.Tops)
		}
		gapp.JSONDATAOK("", gin.H{"now_num": battleUser.JgValue, "last_num": battleUser.CurJgValue, "rank": rank})
	} else {

		rank := "未上榜"
		gapp.JSONDATAOK("", gin.H{"now_num": 0, "last_num": 0, "rank": rank})
	}

}

// 神圣战场商店
func SSBattleStore(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data := gapp.OptSvc.FightSrv.SSBattleStoreData()
	gapp.JSONDATAOK("", gin.H{"goods": data})
}

// 神圣战场领取排行奖励
func SSBattleGetAward(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	t := com.StrTo(c.Query("type")).MustInt()

	ok, msg := gapp.OptSvc.FightSrv.SSBattleGetAward(gapp.Id(), t)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 神圣战场兑换经验
func SSBattleConvertExp(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	num := com.StrTo(c.Query("num")).MustInt()

	ok, msg := gapp.OptSvc.FightSrv.SSBattleConvertExp(gapp.Id(), num)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 神圣战场兑换经验
func SSBattleConvertProp(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	pid := com.StrTo(c.Query("id")).MustInt()
	num := com.StrTo(c.Query("num")).MustInt()

	ok, msg := gapp.OptSvc.FightSrv.SSBattleConvertProp(gapp.Id(), pid, num)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
