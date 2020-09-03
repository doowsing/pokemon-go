package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/pkg/ginapp"
	"pokemon/pkg/services"
	"pokemon/pkg/utils"
	"strconv"
	"strings"
)

var FightCtl = &FightController{}

type FightController struct {
}

func CheckOpenMap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	//id := gapp.Id()
	mapId := com.StrTo(c.DefaultQuery("n", "0")).MustInt()
	if mapId < 1 {
		return
	}
	optType := c.Query("type")
	dMap := services.GetMMap(mapId)
	if dMap == nil {
		gapp.String("1")
		return
	}
	switch optType {
	case "1":
		//user := gapp.OptSrv.UserSrv.GetUserById(id)
		//openMap := strings.Split(user.OpenMap, ",")
		//for _, v := range openMap {
		//	if v == strconv.Itoa(mapId) {
		//		gapp.String("10")
		//		return
		//	}
		//}
		gapp.String("12")
		return
	}
}

func OpenMap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := gapp.Id()
	mapId := com.StrTo(c.DefaultQuery("open", "0")).MustInt()
	if mapId < 1 {
		return
	}
	dMap := services.GetMMap(mapId)
	if dMap == nil {
		gapp.String("没有该地图！")
		return
	}

	user := gapp.OptSrv.UserSrv.GetUserById(id)
	openMap := strings.Split(user.OpenMap, ",")
	for _, v := range openMap {
		if v == strconv.Itoa(mapId) {
			gapp.String("该地图已经打开了!")
			return
		}
	}
	if gapp.OptSrv.FightSrv.OpenMap(id, mapId) {
		gapp.String("地图打开成功!")
	} else {
		gapp.String("您的包裹中没有打开该地图的钥匙！")
	}
}

func TeamModPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := gapp.Id()
	mapId := com.StrTo(c.DefaultQuery("n", "0")).MustInt()
	if mapId < 1 {
		gapp.String("非法参数！")
		return
	}

	mmap := services.GetMMap(mapId)
	if mmap == nil {
		gapp.String("非法参数！")
		return
	}
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	pets := gapp.OptSrv.PetSrv.GetCarryPets(id)
	hasDiff := false
	ifrteamh := 210
	if mmap.ID == 16 || mmap.ID >= 100 {
		hasDiff = true
		ifrteamh -= 80
	}

	for _, pet := range pets {
		pet.GetM()
	}
	fmt.Printf("1 pets:%s\n", pets[0].MModel)
	czlLimit := "无限制"
	if czlItems := strings.Split(mmap.CzlProp, "|"); czlItems[0] != "" {
		czlLimit = czlItems[0]
	}
	inTeam := false
	isLeader := false

	gapp.HTML("page/team_mod.jet.html", gin.H{
		"user":     user,
		"pets":     pets,
		"mmap":     mmap,
		"czlLimit": czlLimit,
		"ifrteamh": ifrteamh,
		"hasDiff":  hasDiff,
		"inTeam":   inTeam,
		"isLeader": isLeader,
	})
}

func TeamInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	inTeam := false
	isLeader := false
	teamList := ""
	now := utils.NowUnix()
	gapp.HTML("page/team.jet.html", gin.H{
		"id":       gapp.Id(),
		"inTeam":   inTeam,
		"isLeader": isLeader,
		"now":      now,
		"teamList": teamList,
	})
}

func FbModPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	mapId := com.StrTo(c.DefaultQuery("mapid", "0")).MustInt()
	if mapId < 1 {
		gapp.Redirect("/function/Expore_Mod.php")
		return
	}
	mmap := services.GetMMap(mapId)
	if mmap == nil {
		gapp.Redirect("/function/Expore_Mod.php")
		return
	}
	map2img := map[int]string{
		151: "02",
		12:  "10",
		13:  "11",
	}
	ImgId, ok := map2img[mapId]

	if !ok {
		ImgId = strconv.Itoa(mapId)
	}
	fbSetting := services.GetFbSetting(mapId)
	if fbSetting == nil {
		gapp.Redirect("/function/Expore_Mod.php")
		return
	}
	id := gapp.Id()
	num := 1
	leftTime := "已开启"
	fbRecord := gapp.OptSrv.FightSrv.GetFbRecord(id, mapId)
	if fbRecord != nil {
		for i, v := range fbSetting.GwIds {
			if v == fbRecord.GpcId {
				num = i + 1
				break
			}
		}
		if t := fbRecord.SrcTime - (utils.NowUnix() - fbRecord.LeftTime); fbRecord.LeftTime > 0 && t > 0 {
			leftTime = strconv.Itoa(t) + "秒"
		}
	}
	gpc := services.GetGpc(fbSetting.GwIds[num-1])
	if gpc == nil {
		gapp.String("地图怪物不存在！")
		return
	}

	user := gapp.OptSrv.UserSrv.GetUserById(id)
	pets := gapp.OptSrv.PetSrv.GetCarryPets(id)
	cpets := []map[string]string{}
	for i, pet := range pets {
		pet.GetM()
		sel := 50
		if pet.ID == user.Mbid {
			sel = 100
		}
		devWidth := 111
		if i == 2 {
			devWidth = 156
		}
		cpets = append(cpets, map[string]string{
			"id":       strconv.Itoa(pet.ID),
			"name":     pet.MModel.Name,
			"img":      pet.MModel.ImgCard,
			"sel":      strconv.Itoa(sel),
			"devWidth": strconv.Itoa(devWidth),
		})
	}
	lpet := len(cpets)
	for i := lpet - 1; i < 3; i++ {
		devWidth := 111
		if i == 2 {
			devWidth = 156
		}
		cpets = append(cpets, map[string]string{
			"id":       strconv.Itoa(0),
			"devWidth": strconv.Itoa(devWidth),
		})
	}
	fmt.Printf("cpets:%s\n", cpets)
	gapp.HTML("page/fb_mod.jet.html", gin.H{
		"nickname":  user.Nickname,
		"mbid":      user.Mbid,
		"mapid":     mapId,
		"introduce": mmap.Description,
		"mapImgId":  ImgId,
		"mapLevel":  fbSetting.Level,
		"leftTime":  leftTime,
		"gwCnt":     len(fbSetting.GwIds),
		"num":       num,
		"gwName":    gpc.Name,
		"cpets":     cpets,
	})

}
