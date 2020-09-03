package npc

import (
	"github.com/gin-gonic/gin"
	"pokemon/game/ginapp"
	"strconv"
)

// 扫雷信息
func SaoLeiInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	bestAwards := gapp.OptSvc.NpcSrv.GetUserSaoleiAward(gapp.Id())
	level, enableSaolei := gapp.OptSvc.UserSrv.GetSaoleiStatus(gapp.Id())
	cgkSum, fhkSum, sxkSum := gapp.OptSvc.NpcSrv.GetSaoleiPropNum(gapp.Id())
	gapp.JSONDATAOK("", gin.H{
		"enable_sl": enableSaolei,
		"awards":    bestAwards,
		"card_num": gin.H{
			"cgk": cgkSum,
			"sxk": sxkSum,
			"fhk": fhkSum,
		},
		"level": level,
	})
}

// 扫雷-刷新奖励
func UpdateSaoLeiAwards(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.UpdateSaoLeiAward(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 扫雷-开始扫雷
func StartSaoLei(c *gin.Context) {

	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	pStr := c.Query("position")
	position, err := strconv.Atoi(pStr)
	if position < 1 || position > 9 || err != nil {
		gapp.JSONDATAOK("参数出错！", gin.H{"result": false})
		return
	}
	msg, result := gapp.OptSvc.NpcSrv.StartSaoLei(gapp.Id(), position)
	gapp.JSONDATAOK(msg, result)
}

// 扫雷-开始闯关
func IntoSaoLei(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.UseSaoleiTicketInto(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 扫雷-复活
func EasterSaoLei(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.EasterSaoLei(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
