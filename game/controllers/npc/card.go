package npc

import (
	"github.com/gin-gonic/gin"
	common2 "pokemon/common"
	"pokemon/game/ginapp"
	"strconv"
)

func CardSeries(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.JSONDATAOK("", gin.H{"series": gapp.OptSvc.NpcSrv.GetCardSeriesDatas()})
}

func UserCardSeries(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	gapp.JSONDATAOK("", gin.H{"cards": gapp.OptSvc.NpcSrv.GetCardSeriesData(gapp.Id(), id)})
}

func CardPrizes(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.JSONDATAOK("", gin.H{"prizes": gapp.OptSvc.NpcSrv.GetCardPrizeDatas(gapp.Id())})

}

func GetCardPrize(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.GetCardPrize(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func CardTitles(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.JSONDATAOK("", gin.H{"titles": gapp.OptSvc.NpcSrv.GetCardTitleDatas(gapp.Id())})

}

func UseCardTitle(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.UseCardTitle(gapp.Id(), id)
	if ok {
		common2.UpdateUserInfo2Chat(gapp.OptSvc.SysSrv.GetChatUserInfo(gapp.Id()))
	}
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func CancelCardTitle(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.CancelCardTitle(gapp.Id(), id)
	if ok {
		common2.UpdateUserInfo2Chat(gapp.OptSvc.SysSrv.GetChatUserInfo(gapp.Id()))
	}
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
