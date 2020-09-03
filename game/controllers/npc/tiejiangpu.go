package npc

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/game/ginapp"
	"pokemon/game/services/common"
	"strconv"
	"strings"
)

// 铁匠铺界面信息
func TieJiangPuPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	updateStr := c.Query("update")
	shopData := gapp.OptSvc.NpcSrv.GetTjpShopGood(updateStr == "true")
	qhData, fjData, xqData := []gin.H{}, []gin.H{}, []gin.H{}
	fjSetting := common.GetWelcome("biodegradable_equipment")
	fjPositions := strings.Split(fjSetting.Content, ",")
	xqSetting := common.GetWelcome("allow_to_use_gam")
	xqPositions := strings.Split(xqSetting.Content, ",")

	props := gapp.OptSvc.PropSrv.GetCarryProps(gapp.Id(), false)
	for _, prop := range props {
		prop.GetM()
		propInfo := gin.H{
			"id":        prop.ID,
			"name":      prop.MModel.Name,
			"price":     prop.MModel.SellJb,
			"vary_id":   prop.MModel.VaryName,
			"vary_name": prop.MModel.VaryNameStr,
			"sum":       prop.Sums,
			"img":       prop.MModel.Img,
		}
		if prop.MModel.VaryName == 9 && prop.Zbing == 0 {
			if prop.MModel.PlusFlag > 0 {
				// 可强化
				qhData = append(qhData, propInfo)
			}
			if com.IsSliceContainsStr(fjPositions, strconv.Itoa(prop.MModel.Position)) {
				// 可分解
				fjData = append(fjData, propInfo)
			}
			if com.IsSliceContainsStr(xqPositions, strconv.Itoa(prop.MModel.Position)) {
				// 可镶嵌
				xqData = append(xqData, propInfo)
			}
		} else if prop.MModel.VaryName == 10 || prop.MModel.VaryName == 11 {
			// 强化消耗物、龙珠
			qhData = append(qhData, propInfo)
		} else if prop.MModel.VaryName == 25 || prop.MModel.VaryName == 26 || prop.MModel.VaryName == 27 {
			// 水晶、洗练石、保底石
			xqData = append(xqData, propInfo)
		}
	}

	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	shopData["qh_list"] = qhData
	shopData["fj_list"] = fjData
	shopData["xq_list"] = xqData
	shopData["jb"] = user.Money
	shopData["prestige"] = user.Prestige
	shopData["left_fj_times"] = gapp.OptSvc.NpcSrv.GetZbFJTimes(gapp.Id())

	gapp.JSONDATAOK("", shopData)
}

// 铁匠铺-分解装备
func FenJieEquip(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.FenjieZb(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 铁匠铺-强化装备
func QiangHuaEquip(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	fzidStr := c.Query("fzid")
	fzid, err := strconv.Atoi(fzidStr)
	if err != nil || fzid < 1 {
		fzid = 0
	}
	ok, msg := gapp.OptSvc.NpcSrv.QiangHuaEquip(gapp.Id(), id, fzid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 铁匠铺-强化装备所需的信息
func QiangHuaEquipInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	result, msg := gapp.OptSvc.NpcSrv.QiangHuaInfo(gapp.Id(), id)
	gapp.JSONDATAOK(msg, result)
}

// 铁匠铺-合成水晶、镶嵌装备
func MergeProps(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id1Str := c.Query("id1")
	id1, err := strconv.Atoi(id1Str)
	if err != nil || id1 < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	id2Str := c.Query("id2")
	id2, err := strconv.Atoi(id2Str)
	if err != nil || id2 < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	fzidStr := c.Query("fzid")
	fzid, err := strconv.Atoi(fzidStr)
	if err != nil || fzid < 1 {
		fzid = 0
	}
	ok, msg := gapp.OptSvc.NpcSrv.MergeProps(gapp.Id(), id1, id2, fzid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
