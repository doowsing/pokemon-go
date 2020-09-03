package npc

import (
	"github.com/gin-gonic/gin"
	"pokemon/common/rcache"
	"pokemon/game/ginapp"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"strconv"
	"strings"
	"time"
)

// 进入游戏的公告、区服名
func GameInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	user := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
	set := common.GetWelcome("ifrc")
	gameName := "遗迹服"
	gapp.JSONDATAOK("", gin.H{
		"guide":     user.NewGuideStep,
		"announce":  set.Content,
		"game_name": gameName,
	})
}

// 中心城镇
func CityPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	rcache.DelPetStatus(user.Mbid)
	rcache.SetInMap(gapp.Id(), 0)
	gapp.JSONDATAOK("", gin.H{"result": true})
}

// 宠物界面信息
func PetsPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSvc.UserSrv.GetUserById(uid)
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(uid)
	petCnt := gapp.OptSvc.PetSrv.GetPetCnt(uid)

	mainPet := &models.UPet{}
	mainPetId := user.Mbid
	// 宠物资料
	carryPets := gapp.OptSvc.PetSrv.GetCarryPets(uid)
	carryPetData := []gin.H{}
	for _, pet := range carryPets {
		pet.GetM()
		if pet.ID == mainPetId {
			mainPet = pet
		}
		carryPetData = append(carryPetData, gin.H{
			"id":    pet.ID,
			"name":  pet.MModel.Name,
			"img":   pet.MModel.ImgCard,
			"level": pet.Level,
			"wx":    pet.MModel.WxName,
		})
	}

	if mainPet.ID == 0 {
		mainPet = carryPets[0]
	}
	kxs := gapp.OptSvc.PetSrv.GetPetKx(mainPet.Kx)

	// 宠物装备
	zbs := gapp.OptSvc.PropSrv.GetPZbs(mainPet.ID)
	zbData := map[string]gin.H{}

	for i := 1; i <= 12; i++ {
		zbData[utils.GetZbPositionName(i)] = nil
	}

	for _, v := range zbs {
		v.GetM()
		if v.MModel != nil {
			position := v.MModel.Position
			if position == 0 {
				position = 11
			} else if position == 11 {
				position = 12
			}
			zbData[utils.GetZbPositionName(position)] = gin.H{
				"id":       v.ID,
				"name":     v.MModel.Name,
				"img":      v.MModel.Img,
				"position": position,
			}
		}
	}
	gapp.OptSvc.FightSrv.GetZbAttr(mainPet, zbs)

	// 宠物技能
	skills := gapp.OptSvc.PetSrv.GetPetSkill(mainPet.ID)
	skillData := []gin.H{}
	for _, skill := range skills {
		skill.GetM()
		skillData = append(skillData, gin.H{
			"id":        skill.ID,
			"name":      skill.MModel.Name,
			"level":     skill.Level,
			"enbale_up": skill.Level < 10 && mainPet.Level <= skill.Level,
		})
	}

	// 可学技能书
	skillBooks := gapp.OptSvc.PropSrv.GetCarryPropsByVaryName(uid, false, 5)
	studySkills := []gin.H{}
	for _, book := range skillBooks {
		s := gapp.OptSvc.PetSrv.GetMskillByPid(book.Pid)
		if s != nil && (s.Wx == 0 || s.Wx == mainPet.Wx) {
			studySkills = append(studySkills, gin.H{
				"id":   book.ID,
				"name": s.Name,
			})
		}
	}

	userData := gin.H{
		"id":           user.ID,
		"nickname":     user.Nickname,
		"sex":          user.Sex,
		"img":          "2" + user.Headimg + ".gif",
		"jinbi":        user.Money,
		"shuijing":     userInfo.Sj,
		"yuanbao":      user.Yb,
		"vip":          user.Vip,
		"showpettimes": userInfo.Bbshow,
		"pk_core":      user.ChallengeRecord,
	}
	mainPetData := gin.H{
		"id":          mainPet.ID,
		"level":       mainPet.Level,
		"name":        mainPet.MModel.Name,
		"wx":          mainPet.MModel.WxName,
		"img":         mainPet.MModel.ImgStand,
		"hp":          mainPet.ZbAttr.Hp,
		"mp":          mainPet.ZbAttr.Mp,
		"ac":          mainPet.ZbAttr.Ac,
		"mc":          mainPet.ZbAttr.Mc,
		"hits":        mainPet.ZbAttr.Hits,
		"miss":        mainPet.ZbAttr.Miss,
		"speed":       mainPet.ZbAttr.Speed,
		"czl":         mainPet.Czl,
		"nowexp":      mainPet.NowExp,
		"maxexp":      mainPet.LExp,
		"kx":          kxs,
		"remarktimes": mainPet.ReMakeTimes,
		"zbs":         zbData,
		"skills":      skillData,
		"skillbooks":  studySkills,
	}
	gapp.JSONDATAOK("", gin.H{
		"user":      userData,
		"pet_count": petCnt,
		"pets":      carryPetData,
		"main_pet":  mainPetData,
	})

}

// 宠物界面-下装备
func PetOffZb(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}

	ok, msg := gapp.OptSvc.PetSrv.PetOffzb(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 用户界面信息
func UserPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSvc.UserSrv.GetUserById(uid)
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(uid)
	petCnt := gapp.OptSvc.PetSrv.GetPetCnt(uid)
	mergename := "未婚"
	if userInfo.Merge > 0 {
		mergeUser := gapp.OptSvc.UserSrv.GetUserById(userInfo.Merge)
		mergename = mergeUser.Nickname
	}
	dlTimes := 1.0
	switch dlTimes {
	case 2:
		dlTimes = 1.5
		break
	case 3:
		dlTimes = 2.0
		break
	case 4:
		dlTimes = 2.5
		break
	case 5:
		dlTimes = 3.0
		break
	default:
		dlTimes = 1.0
	}
	leftDlTime := 0
	if leftDlTime = user.Maxdblexptime + user.AutofitStartTime - utils.NowUnix(); leftDlTime < 0 {
		leftDlTime = 0
	}
	if user.ChallengeRecord != "" {
		user.ChallengeRecord = "胜:0,败:0"
	} else {
		user.ChallengeRecord = "胜:" + strings.ReplaceAll(user.ChallengeRecord, ":", ", 败：")
	}
	tiaozhan := "允许"
	if userInfo.TiaoZhan != 1 {
		tiaozhan = "不允许"
	}
	userData := gin.H{
		"id":              user.ID,
		"nickname":        user.Nickname,
		"sex":             user.Sex,
		"img":             "3" + user.Headimg + ".gif",
		"jinbi":           user.Money,
		"shuijing":        userInfo.Sj,
		"yuanbao":         user.Yb,
		"vip":             user.Vip,
		"last_vip":        user.VipLast,
		"score":           user.Score,
		"pk_core":         user.ChallengeRecord,
		"prestige":        user.Prestige,
		"jpprestige":      user.Jprestige,
		"autotimes_jb":    user.AutoFightTimeM,
		"autotimes_yb":    user.AutoFightTimeYb,
		"autotimes_team":  userInfo.TeamAutoTimes,
		"double_exp":      dlTimes,
		"double_exp_time": leftDlTime,
		"enable_fight":    tiaozhan,
		"merry":           mergename,
	}
	gapp.JSONDATAOK("", gin.H{
		"user":        userData,
		"pet_count":   petCnt,
		"friend_list": []string{},
		"black_list":  []string{},
	})
}

// 公告牌界面信息
func PublicPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	forceUpdateStr := c.Query("update")
	forceUpdate := false
	if forceUpdateStr == "true" {
		forceUpdate = true
	}
	rankList := gapp.OptSvc.SysSrv.GetPublicRankLists(forceUpdate)
	// 消费排行榜
	openFlag, userList, timeSet, userCon := gapp.OptSvc.SysSrv.GetConsumptionInfo(gapp.Id())
	gapp.JSONDATAOK("", gin.H{
		"level_list":       rankList["level"],
		"sscz_list":        rankList["sscz"],
		"cz_list":          rankList["cz"],
		"czzc_list":        rankList["czzc"],
		"consumption_list": userList,
		"consumption_info": gin.H{
			"is_open":          openFlag,
			"consumption_user": userCon,
			"consumption_time": timeSet,
		},
		"public_announce": gapp.OptSvc.SysSrv.GetPublicContent(),
	})
}

// 神秘商店界面信息
func SmShopPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	updateStr := c.Query("update")
	smShopGoods := gapp.OptSvc.PropSrv.GetSmShopGood(updateStr == "true")
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
	smShopGoods["sj"] = userInfo.Sj
	smShopGoods["yb"] = user.Yb
	smShopGoods["vip"] = user.Vip
	gapp.JSONDATAOK("", smShopGoods)
}

// 神秘商店抢购列表
func SmShopQgList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.JSONDATAOK("", gapp.OptSvc.PropSrv.GetSmShopQgList())
}

// 道具商店界面信息
func DjShopPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	updateStr := c.Query("update")
	shopData := gapp.OptSvc.PropSrv.GetDjShopGood(updateStr == "true")
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	shopData["jb"] = user.Money
	shopData["prestige"] = user.Prestige
	gapp.JSONDATAOK("", shopData)
}

// 牧场界面信息
func MuchangPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSvc.UserSrv.GetUserById(uid)
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(uid)

	// 检查是否有密码
	need_pass := gapp.OptSvc.UserSrv.CheckNeedPwd(c.Query("passwd"), user.McPwd)

	allPets := []*models.UPet{}
	if !need_pass {
		allPets = gapp.OptSvc.PetSrv.GetAllPets(uid)
	} else {
		allPets = gapp.OptSvc.PetSrv.GetCarryPets(uid)
	}

	mcPets, crPets, ableTgPets, inTgPets := []gin.H{}, []gin.H{}, []gin.H{}, []gin.H{}
	for _, p := range allPets {
		p.GetM()
		//fmt.Printf("pet id:%d\n", p.ID)
		if p.Muchang == 0 && p.TgFlag != 1 {
			crPets = append(crPets, gin.H{
				"id":    p.ID,
				"name":  p.MModel.Name,
				"img":   p.MModel.ImgCard,
				"level": p.Level,
				"wx":    p.MModel.WxName,
			})
		} else if p.Muchang == 1 {
			if p.TgFlag == 0 {
				mcPets = append(mcPets, gin.H{
					"id":    p.ID,
					"name":  p.MModel.Name,
					"img":   p.MModel.ImgCard,
					"level": p.Level,
					"wx":    p.MModel.WxName,
				})
				if p.Level >= 10 {
					ableTgPets = append(ableTgPets, gin.H{
						"id":   p.ID,
						"name": p.MModel.Name,
					})
				}
			} else {
				inTgPets = append(inTgPets, gin.H{
					"id":      p.ID,
					"name":    p.MModel.Name,
					"tg_flag": p.TgFlag,
				})
			}

		}
	}
	gapp.JSONDATAOK("", gin.H{
		"mc_pets":      mcPets,
		"carry_pets":   crPets,
		"in_tg_pets":   inTgPets,
		"able_tg_pets": ableTgPets,
		"mc_max":       user.McPlace,
		"main_pet":     user.Mbid,
		"sj":           userInfo.Sj,
		"tg_time":      user.TgTime,
		"tg_max":       user.TgPlace,
		"need_pass":    need_pass,
	})

}

// 仓库界面信息
func CangkuPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSvc.UserSrv.GetUserById(uid)
	ckData := []gin.H{}

	// 检查是否有密码
	need_pass := gapp.OptSvc.UserSrv.CheckNeedPwd(c.Query("passwd"), user.CkPwd)
	if !need_pass {
		ckData = gapp.OptSvc.PropSrv.GetCkPropData(gapp.Id())
	}
	gapp.JSONDATAOK("", gin.H{
		"ckprops":   ckData,
		"yb":        user.Yb,
		"jb":        user.Money,
		"ck_max":    user.BasePlace,
		"need_pass": need_pass,
	})

}

// 拍卖所界面信息
func PaiMaiPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := gapp.Id()
	user := gapp.OptSvc.UserSrv.GetUserById(id)
	userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(id)

	jbProps, sjProps, ybProps := gapp.OptSvc.PropSrv.GetPmProps()
	jbData, sjData, ybData, myData := []gin.H{}, []gin.H{}, []gin.H{}, []gin.H{}
	selfPmProps := gapp.OptSvc.PropSrv.GetSelfPmProps(id)
	for _, prop := range jbProps {
		prop.GetM()
		jbData = append(jbData, gin.H{
			"id":      prop.ID,
			"name":    prop.MModel.Name,
			"price":   prop.Psell,
			"vary_id": prop.MModel.VaryName,
			"num":     prop.Psum,
		})
	}
	for _, prop := range sjProps {
		prop.GetM()
		sjData = append(sjData, gin.H{
			"id":      prop.ID,
			"name":    prop.MModel.Name,
			"price":   prop.Psj,
			"vary_id": prop.MModel.VaryName,
			"num":     prop.Psum,
		})
	}
	for _, prop := range ybProps {
		prop.GetM()
		ybData = append(ybData, gin.H{
			"id":      prop.ID,
			"name":    prop.MModel.Name,
			"price":   prop.Pyb,
			"vary_id": prop.MModel.VaryName,
			"num":     prop.Psum,
		})
	}
	now := time.Now()
	for _, prop := range selfPmProps {
		leftTimeStr := "已过期"
		if leftTime := time.Unix(int64(prop.Petime), 0).Sub(now); leftTime > 0 {
			leftTimeStr = utils.DurationFormatHms(leftTime)
		}
		myData = append(myData, gin.H{
			"id":        prop.ID,
			"name":      prop.MModel.Name,
			"price":     prop.PmMoneyStr,
			"vary_id":   prop.MModel.VaryName,
			"left_time": leftTimeStr,
			"num":       prop.Psum,
		})
	}
	shopData := gin.H{
		"jb_list": jbData,
		"sj_list": sjData,
		"yb_list": ybData,
		"my_list": myData,
		"pm_money": gin.H{
			"jb": user.PaiMoney,
			"sj": userInfo.Paisj,
			"yb": userInfo.Paiyb,
		},
		"jb": user.Money,
		"sj": userInfo.Sj,
		"yb": user.Yb,
	}
	gapp.JSONDATAOK("", shopData)

}

// 欢迎界面
func WelcomeContent(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	welcome := common.GetWelcome("welcome")
	welimg := common.GetWelcome("welimg")
	//href := services.GetWelcome("href")
	welcontent := common.GetWelcome("welcontent")
	gapp.JSONDATAOK("", gin.H{
		"pet_produce":      welcontent.Content,
		"pet_img":          strings.ReplaceAll(welimg.Content, "/images/welcome/", ""),
		"announce_content": welcome.Content,
	})
}
