package npc

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/game/ginapp"
	"pokemon/game/services/common"
	"strconv"
	"strings"
)

// 宠物神殿
func PetSdPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	//idStr := c.Query("id")

	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
	carryPets := gapp.OptSvc.PetSrv.GetCarryPets(user.ID)
	//
	//mbid := 0
	//if idStr != "" {
	//	id := com.StrTo(idStr).MustInt()
	//	if id != 0 {
	//		mbid = id
	//	}
	//}
	//mainPet := &models.UPet{}
	sdData := gapp.OptSvc.NpcSrv.GetPetSdPropInfo(user.ID)

	carryPetData := []gin.H{}
	for _, pet := range carryPets {
		pet.GetM()
		propIds := strings.Split(pet.MModel.ReMakePid, ",")
		petIds := strings.Split(pet.MModel.ReMakeId, ",")
		levels := strings.Split(pet.MModel.ReMakeLevel, ",")

		reMakeA := gin.H{}
		reMakeB := gin.H{}
		if pet.MModel.Wx == 7 {
			ssjhRule := common.GetSSJhRule(pet.MModel.ID)
			if ssjhRule == nil {
				reMakeA["prop"] = "无"
				reMakeA["topet"] = "无"
				reMakeA["level"] = 0
				reMakeA["jb"] = 0
			} else {
				levelItems := strings.Split(ssjhRule.NeedLevels, ",")
				propItems := strings.Split(ssjhRule.NeedProps, ",")
				if len(levelItems) <= pet.ReMakeTimes || len(propItems) <= pet.ReMakeTimes {
					reMakeA["prop"] = "无"
					reMakeA["topet"] = "无"
					reMakeA["level"] = 0
					reMakeA["jb"] = 0
				} else {
					items := strings.Split(propItems[pet.ReMakeTimes], ":")
					prop := common.GetMProp(com.StrTo(items[0]).MustInt())
					reMakeA["prop"] = prop.Name
					reMakeA["topet"] = pet.MModel.Name
					reMakeA["level"] = levelItems[pet.ReMakeTimes]
					reMakeA["jb"] = (ssjhRule.ZsProgress + pet.ReMakeTimes) * 10000
				}
			}
			reMakeB["prop"] = "无"
			reMakeB["topet"] = "无"
			reMakeB["level"] = 0
			reMakeB["jb"] = 0

		} else {
			apropIds := strings.Split(propIds[0], "|")
			aprop := common.GetMProp(com.StrTo(apropIds[0]).MustInt())
			if aprop == nil {
				reMakeA["prop"] = "无"
			} else {
				reMakeA["prop"] = aprop.Name
			}
			apet := common.GetMpet(com.StrTo(petIds[0]).MustInt())
			if apet == nil {
				reMakeA["topet"] = "无"
			} else {
				reMakeA["topet"] = apet.Name
			}
			reMakeA["level"] = com.StrTo(levels[0]).MustInt()
			reMakeA["jb"] = 1000

			//fmt.Printf("进化所需道具：%s\n", propIds)
			if len(propIds) > 1 {
				bpropIds := strings.Split(propIds[1], "|")
				bprop := common.GetMProp(com.StrTo(bpropIds[0]).MustInt())
				if bprop == nil {
					reMakeB["prop"] = "无"
				} else {
					reMakeB["prop"] = bprop.Name
				}
			} else {
				reMakeB["prop"] = "无"
			}
			if len(petIds) > 1 {
				bpet := common.GetMpet(com.StrTo(petIds[1]).MustInt())
				if bpet == nil {
					reMakeB["topet"] = "无"
				} else {
					reMakeB["topet"] = bpet.Name
				}
			} else {
				reMakeB["topet"] = "无"
			}
			if len(levels) > 1 {
				reMakeB["level"] = com.StrTo(levels[1]).MustInt()
			} else {
				reMakeB["level"] = 0
			}
			reMakeB["jb"] = 1000
		}

		cqjinbi := 0
		if pet.MModel.Wx < 7 {
			if pet.CC > 30 {
				cqjinbi = int(pet.CC * 1000)
				if cqjinbi > 6000000 {
					cqjinbi = 6000000
				}
			}
		}
		carryPetData = append(carryPetData, gin.H{
			"id":          pet.ID,
			"name":        pet.MModel.Name,
			"img":         pet.MModel.ImgCard,
			"level":       pet.Level,
			"A":           reMakeA,
			"B":           reMakeB,
			"cqjinbi":     cqjinbi,
			"remaketimes": pet.ReMakeTimes,
			"czl":         pet.Czl,
			"is_ss":       pet.MModel.Wx > 6,
		})
	}

	sdData["carrypets"] = carryPetData
	sdData["ssczl"] = userInfo.CzlSS
	sdData["hc_false_num"] = userInfo.HechengNums
	gapp.JSONDATAOK("", sdData)
}

// 宠物神殿-非神圣宠进化
func PetEvolution(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	pathStr := c.Query("path")
	apath := pathStr != "1"
	fzidStr := c.Query("fzid")
	fzid, err := strconv.Atoi(fzidStr)
	if err != nil {
		fzid = 0
	}
	ok, msg := gapp.OptSvc.PetSrv.Evolution(gapp.Id(), id, apath, fzid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 宠物神殿-宠物合成
func PetMerge(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	aidStr := c.Query("aid")
	aid, err := strconv.Atoi(aidStr)
	if aid < 1 || err != nil {
		gapp.JSONDATAOK("请选择主宠！", gin.H{"result": false})
		return
	}

	bidStr := c.Query("bid")
	bid, err := strconv.Atoi(bidStr)
	if bid < 1 || err != nil {
		gapp.JSONDATAOK("请选择副宠！", gin.H{"result": false})
		return
	}

	apidStr := c.Query("apid")
	apid, err := strconv.Atoi(apidStr)
	if err != nil {
		apid = 0
	}

	bpidStr := c.Query("bpid")
	bpid, err := strconv.Atoi(bpidStr)
	if err != nil {
		bpid = 0
	}
	ok, msg := gapp.OptSvc.PetSrv.Emerge(gapp.Id(), aid, bid, apid, bpid, c.Query("zbcheck") == "true", c.Query("protectcheck") == "true")
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 宠物神殿-宠物转生
func PetZhuansheng(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	aidStr := c.Query("aid")
	aid, err := strconv.Atoi(aidStr)
	if aid < 1 || err != nil {
		gapp.JSONDATAOK("请选择主宠！", gin.H{"result": false})
		return
	}

	bidStr := c.Query("bid")
	bid, err := strconv.Atoi(bidStr)
	if bid < 1 || err != nil {
		gapp.JSONDATAOK("请选择副宠", gin.H{"result": false})
		return
	}

	cidStr := c.Query("cid")
	cid, err := strconv.Atoi(cidStr)
	if bid < 1 || err != nil {
		gapp.JSONDATAOK("请选择涅槃兽！", gin.H{"result": false})
		return
	}

	apidStr := c.Query("apid")
	apid, err := strconv.Atoi(apidStr)
	if err != nil {
		apid = 0
	}

	bpidStr := c.Query("bpid")
	bpid, err := strconv.Atoi(bpidStr)
	if err != nil {
		bpid = 0
	}
	ok, msg := gapp.OptSvc.PetSrv.ZhuanSheng(gapp.Id(), aid, bid, cid, apid, bpid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 宠物神殿-宠物抽取成长
func PetCqCzl(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	apidStr := c.Query("apid")
	apid, err := strconv.Atoi(apidStr)
	if err != nil {
		apid = 0
	}

	bpidStr := c.Query("bpid")
	bpid, err := strconv.Atoi(bpidStr)
	if err != nil {
		bpid = 0
	}
	ok, msg := gapp.OptSvc.PetSrv.Chouqu(gapp.Id(), id, apid, bpid)
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok, "ss_czl": userInfo.CzlSS})
}

// 宠物神殿-神圣宠物转化成长
func PetZhCzl(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	czlStr := c.Query("czl")
	czl, err := strconv.Atoi(czlStr)
	if czl < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.PetSrv.Zhuanhua(gapp.Id(), id, czl)
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok, "ss_czl": userInfo.CzlSS})
}

// 宠物神殿-神圣宠进化
func PetSSEvolution(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	fzidStr := c.Query("fzid")
	fzid, err := strconv.Atoi(fzidStr)
	if err != nil {
		fzid = 0
	}
	ok, msg := gapp.OptSvc.PetSrv.SSEvolution(gapp.Id(), id, fzid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 宠物神殿-神圣宠转生信息
func PetSSZhuanShengInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	gapp.JSONDATAOK("", gin.H{"path": gapp.OptSvc.PetSrv.SSZhuanshengInfo(gapp.Id(), id)})

}

// 宠物神殿-神圣宠转生
func PetSSZhuanSheng(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	toidStr := c.Query("toid")
	toid, err := strconv.Atoi(toidStr)
	if toid < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}

	apidStr := c.Query("apid")
	apid, err := strconv.Atoi(apidStr)
	if err != nil {
		apid = 0
	}

	bpidStr := c.Query("bpid")
	bpid, err := strconv.Atoi(bpidStr)
	if err != nil {
		bpid = 0
	}
	ok, msg := gapp.OptSvc.PetSrv.SSZhuanSheng(gapp.Id(), id, toid, apid, bpid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
