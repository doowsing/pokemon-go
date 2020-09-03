package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/common/persistence"
	"pokemon/game/ginapp"
	"pokemon/game/services/common"
	"strings"
)

func ShowEggSetting(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	typeStr := c.Query("type")
	typeName := ""
	if typeStr == "1" {
		typeName = "golden_eggs"
	} else if typeStr == "2" {
		typeName = "silver_eggs"
	} else if typeStr == "3" {
		typeName = "copper_eggs"
	}
	if typeName == "" {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	eggSetting := common.GetWelcome(typeName)
	if eggSetting == nil {
		gapp.JSONDATAOK("参数错误！", nil)
	}
	awardData := []gin.H{}
	for _, s := range strings.Split(eggSetting.Content, ",") {
		items := strings.Split(s, ":")
		pid := com.StrTo(items[0]).MustInt()
		rankItems := strings.Split(items[4], "-")
		prop := common.GetMProp(pid)
		awardData = append(awardData, gin.H{
			"id":   prop.ID,
			"name": prop.Name,
			"rate": com.StrTo(rankItems[1]).MustInt() - com.StrTo(rankItems[0]).MustInt(),
		})
	}
	gapp.JSONDATAOK("", gin.H{"type": typeName, "awards": awardData})
}

func ShowRedisClusterStatus(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.JSONDATAOK("", persistence.GetRedisCluster().Stats())
}
