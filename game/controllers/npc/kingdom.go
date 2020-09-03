package npc

import (
	"github.com/gin-gonic/gin"
	"pokemon/game/ginapp"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"strconv"
	"strings"
	"time"
)

// 皇宫-界面信息
func KingPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	kingData := gin.H{}
	awardData := gapp.OptSvc.NpcSrv.KingAwards()
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
	now := time.Now()
	getAwardItems := strings.Split(userInfo.PrizeItems, "|")
	kingData["day_award_status"] = false
	if getAwardItems[0] != "" {
		if lastPrizeDay, err := utils.YmdStrParseTime(getAwardItems[0]); err == nil {
			if now.Year() == lastPrizeDay.Year() && now.Month() == lastPrizeDay.Month() && now.Day() == lastPrizeDay.Day() {
				kingData["day_award_status"] = true
			}
		}
	}

	kingData["week_award_status"] = false
	if len(getAwardItems) > 1 && getAwardItems[1] != "" {
		if lastPrizeDay, err := utils.YmdStrParseTime(getAwardItems[1]); err == nil {
			year, week := lastPrizeDay.ISOWeek()
			nyear, nweek := now.ISOWeek()
			if year == nyear && week == nweek {
				kingData["week_award_status"] = true
			}
		}
	}
	kingData["is_weeken"] = false
	if now.Weekday() == 0 || now.Weekday() == 6 {
		kingData["is_weeken"] = true
	}

	kingData["holiday_award_status"] = false
	if len(getAwardItems) > 2 && getAwardItems[2] != "" {
		if lastPrizeDay, err := utils.YmdStrParseTime(getAwardItems[2]); err == nil {
			if now.Year() == lastPrizeDay.Year() && now.Month() == lastPrizeDay.Month() && now.Day() == lastPrizeDay.Day() {
				kingData["holiday_award_status"] = true
			}
		}
	}

	kingData["day_awards"] = awardData["day"]
	kingData["week_awards"] = awardData["week"]
	kingData["holiday_awards"] = awardData["holiday"]
	kingData["prestige"] = user.Prestige
	kingData["jprestige"] = user.Jprestige
	dataSetting := common.GetWelcome("dati")
	if dataSetting != nil {
		kingData["dati_content"] = dataSetting.Content
	} else {
		kingData["dati_content"] = "活动内容，见官方网站通知。"
	}
	datiPlayer := gapp.OptSvc.UserSrv.GetDatiPlayer(gapp.Id())
	if datiPlayer != nil {
		kingData["dati_right_time"] = datiPlayer.OkSum
	} else {
		kingData["dati_right_time"] = 0
	}
	kingData["danquan"] = gapp.OptSvc.NpcSrv.DanQuanCnt(gapp.Id())
	gapp.JSONDATAOK("", kingData)
}

// 皇宫-领取皇宫日常、周末、假期奖励
func GetDayPrize(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	prizeType := c.Query("type")
	ok, msg := gapp.OptSvc.NpcSrv.GetKingAwards(gapp.Id(), prizeType)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 皇宫-领取皇宫日常、周末、假期奖励
func GivePrestige(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	numidStr := c.Query("num")
	num, err := strconv.Atoi(numidStr)
	if num < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok := gapp.OptSvc.UserSrv.GivePrestige(gapp.Id(), num)
	var msg string
	if ok {
		msg = "捐赠贵族威望成功！"
	} else {
		msg = "您没有那么多威望可以捐赠！"
	}
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 皇宫-砸蛋
func Zadan(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	prizeType := c.Query("type")

	pStr := c.Query("position")
	position, err := strconv.Atoi(pStr)
	if position < 0 || position > 4 || err != nil {
		gapp.JSONDATAOK("参数出错！", gin.H{"result": false})
		return
	}
	ok, msg, leftSum, awards := gapp.OptSvc.NpcSrv.Zadan(gapp.Id(), position, prizeType)
	gapp.JSONDATAOK(msg, gin.H{
		"result": ok,
		"sum":    leftSum,
		"awards": awards,
	})
}
