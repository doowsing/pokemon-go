package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"html/template"
	"pokemon/pkg/ginapp"
	"pokemon/pkg/models"
	"pokemon/pkg/services"
	"pokemon/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func LoginPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	gapp.HTML("page/login.jet.html", gin.H{})
}

func WelcomePage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	welcomeWord := gapp.OptSrv.SysSrv.GetWelcome("welcome").Content
	welimg := gapp.OptSrv.SysSrv.GetWelcome("welimg").Content
	welcontent := gapp.OptSrv.SysSrv.GetWelcome("welcontent").Content
	autosum := 0
	gapp.HTML("page/welcome.jet.html", gin.H{
		"welcomeWord": welcomeWord,
		"welimg":      welimg,
		"welcontent":  welcontent,
		"autosum":     autosum,
	})
}

func HomePage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.String("ok")
}

func IframePage(c *gin.Context) {

	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ifrc := gapp.OptSrv.SysSrv.GetWelcome("ifrc").Content
	gapp.HTML("page/iframe.jet.html", gin.H{
		"ifrc": template.HTML(ifrc),
	})
}

func IndexPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	gapp.HTML("page/index.jet.html", gin.H{})
}

func GamePage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := gapp.Id()
	iframe := gapp.OptSrv.SysSrv.GetWelcome("iframe").Content
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(id)
	//gapp.String("欢迎您，%s", user.Nickname)
	gameName := "曙光服"
	var (
		htmlName, chatServer string
	)
	if c.DefaultQuery("ish5", "1") == "1" {
		htmlName = "game_h5.jet.html"
		ip, port := "127.0.0.1", 1986
		sessionid := gapp.Session().SessionId()
		fmt.Printf("sessionId:%s\n", sessionid)
		str := utils.Md5(sessionid + user.Account + user.Password + "0" + user.Nickname)
		chatServer = fmt.Sprintf(`%s|%s|%s|30|%s|%s|%s|%s|%s|%s`, ip, sessionid, port, user.ID, user.Account, user.Password, 0, user.Nickname, str)
	} else {
		htmlName = "game_flash.jet.html"
		chatServer = "127.0.0.1:1986"
	}

	gapp.HTML("page/"+htmlName, gin.H{
		"user":       user,
		"userinfo":   userInfo,
		"title":      gameName,
		"iframeurl":  iframe,
		"chatServer": chatServer,
	})
}

func ExporePage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := gapp.Id()
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	gapp.HTML("page/expore_mod.jet.html", gin.H{
		"openmap": user.OpenMap,
	})
}

func ExporeNewPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := gapp.Id()
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	if n := c.DefaultQuery("n", "1"); n == "2" {
		gapp.HTML("page/expore_new2_mod.jet.html", gin.H{"openmap": user.OpenMap})
	} else {
		gapp.HTML("page/expore_new_mod.jet.html", gin.H{"openmap": user.OpenMap})
	}
}

// 宠物界面信息
func PetsPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSrv.UserSrv.GetUserById(uid)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(uid)
	petCnt := gapp.OptSrv.PetSrv.GetPetCnt(uid)

	mainPet := &models.UPet{}
	mainPetId := user.Mbid
	// 宠物资料
	carryPets := gapp.OptSrv.PetSrv.GetCarryPets(uid)
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
	kxs := gapp.OptSrv.PetSrv.GetPetKx(mainPet.Kx)

	// 宠物装备
	zbs := gapp.OptSrv.PropSrv.GetPZbs(mainPet.ID)
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
	gapp.OptSrv.FightSrv.GetZbAttr(mainPet, zbs)

	// 宠物技能
	skills := gapp.OptSrv.PetSrv.GetPetSkill(mainPet.ID)
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
	skillBooks := gapp.OptSrv.PropSrv.GetCarryPropsByVaryName(uid, false, 5)
	studySkills := []gin.H{}
	for _, book := range skillBooks {
		s := gapp.OptSrv.PetSrv.GetMskillByPid(book.Pid)
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
		"hp":          mainPet.Hp,
		"mp":          mainPet.Mp,
		"ac":          mainPet.Ac,
		"mc":          mainPet.Mc,
		"hits":        mainPet.Hits,
		"miss":        mainPet.Miss,
		"speed":       mainPet.Speed,
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

	ok, msg := gapp.OptSrv.PetSrv.PetOffzb(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 用户界面信息
func UserPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSrv.UserSrv.GetUserById(uid)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(uid)
	petCnt := gapp.OptSrv.PetSrv.GetPetCnt(uid)
	mergename := "未婚"
	if userInfo.Merge > 0 {
		mergeUser := gapp.OptSrv.UserSrv.GetUserById(userInfo.Merge)
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
	rankList := gapp.OptSrv.SysSrv.GetPublicRankLists(forceUpdate)
	// 消费排行榜
	openFlag, userList, timeSet, userCon := gapp.OptSrv.SysSrv.GetConsumptionInfo(gapp.Id())
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
		"public_announce": gapp.OptSrv.SysSrv.GetPublicContent(),
	})
}

// 神秘商店界面信息
func SmShopPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	updateStr := c.Query("update")
	smShopGoods := gapp.OptSrv.PropSrv.GetSmShopGood(updateStr == "true")
	user := gapp.OptSrv.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(gapp.Id())
	smShopGoods["sj"] = userInfo.Sj
	smShopGoods["yb"] = user.Yb
	smShopGoods["vip"] = user.Vip
	gapp.JSONDATAOK("", smShopGoods)
}

// 神秘商店抢购列表
func SmShopQgList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.JSONDATAOK("", gapp.OptSrv.PropSrv.GetSmShopQgList())
}

// 道具商店界面信息
func DjShopPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	updateStr := c.Query("update")
	shopData := gapp.OptSrv.PropSrv.GetDjShopGood(updateStr == "true")
	user := gapp.OptSrv.UserSrv.GetUserById(gapp.Id())
	shopData["jb"] = user.Money
	shopData["prestige"] = user.Prestige
	gapp.JSONDATAOK("", shopData)
}

// 牧场界面信息
func MuchangPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	user := gapp.OptSrv.UserSrv.GetUserById(uid)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(uid)

	// 检查是否有密码
	need_pass := gapp.OptSrv.UserSrv.CheckNeedPwd(c.Query("passwd"), user.McPwd)

	allPets := []*models.UPet{}
	if !need_pass {
		allPets = gapp.OptSrv.PetSrv.GetAllPets(uid)
	} else {
		allPets = gapp.OptSrv.PetSrv.GetCarryPets(uid)
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
	user := gapp.OptSrv.UserSrv.GetUserById(uid)
	ckData := []gin.H{}

	// 检查是否有密码
	need_pass := gapp.OptSrv.UserSrv.CheckNeedPwd(c.Query("passwd"), user.CkPwd)
	if !need_pass {
		ckData = gapp.OptSrv.PropSrv.GetCkPropData(gapp.Id())
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
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(id)

	jbProps, sjProps, ybProps := gapp.OptSrv.PropSrv.GetPmProps()
	jbData, sjData, ybData, myData := []gin.H{}, []gin.H{}, []gin.H{}, []gin.H{}
	selfPmProps := gapp.OptSrv.PropSrv.GetSelfPmProps(id)
	for _, prop := range jbProps {
		prop.GetM()
		jbData = append(jbData, gin.H{
			"id":      prop.ID,
			"name":    prop.MModel.Name,
			"price":   prop.Psell,
			"vary_id": prop.MModel.VaryName,
		})
	}
	for _, prop := range sjProps {
		prop.GetM()
		sjData = append(sjData, gin.H{
			"id":      prop.ID,
			"name":    prop.MModel.Name,
			"price":   prop.Psell,
			"vary_id": prop.MModel.VaryName,
		})
	}
	for _, prop := range ybProps {
		prop.GetM()
		ybData = append(ybData, gin.H{
			"id":      prop.ID,
			"name":    prop.MModel.Name,
			"price":   prop.Psell,
			"vary_id": prop.MModel.VaryName,
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

// 铁匠铺界面信息
func TieJiangPuPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	updateStr := c.Query("update")
	shopData := gapp.OptSrv.PropSrv.GetTjpShopGood(updateStr == "true")
	qhData, fjData, xqData := []gin.H{}, []gin.H{}, []gin.H{}
	fjSetting := services.GetWelcome("biodegradable_equipment")
	fjPositions := strings.Split(fjSetting.Content, ",")
	xqSetting := services.GetWelcome("allow_to_use_gam")
	xqPositions := strings.Split(xqSetting.Content, ",")

	props := gapp.OptSrv.PropSrv.GetCarryProps(gapp.Id(), false)
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

	user := gapp.OptSrv.UserSrv.GetUserById(gapp.Id())
	shopData["qh_list"] = qhData
	shopData["fj_list"] = fjData
	shopData["xq_list"] = xqData
	shopData["jb"] = user.Money
	shopData["prestige"] = user.Prestige
	shopData["left_fj_times"] = gapp.OptSrv.PropSrv.GetZbFJTimes(gapp.Id())

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
	ok, msg := gapp.OptSrv.PropSrv.FenjieZb(gapp.Id(), id)
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
	ok, msg := gapp.OptSrv.PropSrv.QiangHuaEquip(gapp.Id(), id, fzid)
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
	result, msg := gapp.OptSrv.PropSrv.QiangHuaInfo(gapp.Id(), id)
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
	ok, msg := gapp.OptSrv.PropSrv.MergeProps(gapp.Id(), id1, id2, fzid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 皇宫-界面信息
func KingPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	kingData := gin.H{}
	awardData := gapp.OptSrv.PropSrv.KingAwards()
	user := gapp.OptSrv.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(gapp.Id())
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
	kingData["jpprestige"] = user.Jprestige
	dataSetting := services.GetWelcome("dati")
	if dataSetting != nil {
		kingData["dati_content"] = dataSetting.Content
	} else {
		kingData["dati_content"] = "活动内容，见官方网站通知。"
	}
	datiPlayer := gapp.OptSrv.UserSrv.GetDatiPlayer(gapp.Id())
	if datiPlayer != nil {
		kingData["dati_right_time"] = datiPlayer.OkSum
	} else {
		kingData["dati_right_time"] = 0
	}
	kingData["danquan"] = gapp.OptSrv.PropSrv.DanQuanCnt(gapp.Id())
	gapp.JSONDATAOK("", kingData)
}

// 皇宫-领取皇宫日常、周末、假期奖励
func GetDayPrize(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	prizeType := c.Query("type")
	ok, msg := gapp.OptSrv.PropSrv.GetKingAwards(gapp.Id(), prizeType)
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
	ok := gapp.OptSrv.UserSrv.GivePrestige(gapp.Id(), num)
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
	ok, msg, leftSum, awards := gapp.OptSrv.PropSrv.Zadan(gapp.Id(), position, prizeType)
	gapp.JSONDATAOK(msg, gin.H{
		"result": ok,
		"sum":    leftSum,
		"awards": awards,
	})
}

// 扫雷信息
func SaoLeiInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	bestAwards := gapp.OptSrv.PropSrv.GetUserSaoleiAward(gapp.Id())
	level, enableSaolei := gapp.OptSrv.UserSrv.GetSaoleiStatus(gapp.Id())
	cgkSum, fhkSum, sxkSum := gapp.OptSrv.PropSrv.GetSaoleiPropNum(gapp.Id())
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
	ok, msg := gapp.OptSrv.PropSrv.UpdateSaoLeiAward(gapp.Id())
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
	msg, result := gapp.OptSrv.PropSrv.StartSaoLei(gapp.Id(), position)
	gapp.JSONDATAOK(msg, result)
}

// 扫雷-开始闯关
func IntoSaoLei(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSrv.PropSrv.UseSaoleiTicketInto(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 扫雷-复活
func EasterSaoLei(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSrv.PropSrv.EasterSaoLei(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 宠物神殿
func PetSdPage(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	//idStr := c.Query("id")

	user := gapp.OptSrv.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(gapp.Id())
	carryPets := gapp.OptSrv.PetSrv.GetCarryPets(user.ID)
	//
	//mbid := 0
	//if idStr != "" {
	//	id := com.StrTo(idStr).MustInt()
	//	if id != 0 {
	//		mbid = id
	//	}
	//}
	//mainPet := &models.UPet{}
	sdData := gapp.OptSrv.PropSrv.GetPetSdPropInfo(user.ID)
	carryPetData := []gin.H{}
	for _, pet := range carryPets {
		pet.GetM()
		propIds := strings.Split(pet.MModel.ReMakePid, ",")
		petIds := strings.Split(pet.MModel.ReMakeId, ",")
		levels := strings.Split(pet.MModel.ReMakeLevel, ",")

		reMakeA := gin.H{}
		reMakeB := gin.H{}
		if pet.MModel.Wx == 6 {
			ssjhRule := services.GetSSJhRule(pet.MModel.ID)
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
					prop := services.GetMProp(com.StrTo(items[0]).MustInt())
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
			aprop := services.GetMProp(com.StrTo(apropIds[0]).MustInt())
			if aprop == nil {
				reMakeA["prop"] = "无"
			} else {
				reMakeA["prop"] = aprop.Name
			}
			apet := services.GetMpet(com.StrTo(petIds[0]).MustInt())
			if apet == nil {
				reMakeA["topet"] = "无"
			} else {
				reMakeA["topet"] = apet.Name
			}
			reMakeA["level"] = com.StrTo(levels[0]).MustInt()
			reMakeA["jb"] = 1000

			fmt.Printf("进化所需道具：%s\n", propIds)
			if len(propIds) > 1 {
				bpropIds := strings.Split(propIds[1], "|")
				bprop := services.GetMProp(com.StrTo(bpropIds[0]).MustInt())
				if bprop == nil {
					reMakeB["prop"] = "无"
				} else {
					reMakeB["prop"] = bprop.Name
				}
			} else {
				reMakeB["prop"] = "无"
			}
			if len(petIds) > 1 {
				bpet := services.GetMpet(com.StrTo(petIds[0]).MustInt())
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
	ok, msg := gapp.OptSrv.PetSrv.Evolution(gapp.Id(), id, apath, fzid)
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
	ok, msg := gapp.OptSrv.PetSrv.Emerge(gapp.Id(), aid, bid, apid, bpid, c.Query("zbcheck") == "true", c.Query("protectcheck") == "true")
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
	ok, msg := gapp.OptSrv.PetSrv.ZhuanSheng(gapp.Id(), aid, bid, cid, apid, bpid)
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
	ok, msg := gapp.OptSrv.PetSrv.Chouqu(gapp.Id(), id, apid, bpid)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(gapp.Id())
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
	ok, msg := gapp.OptSrv.PetSrv.Zhuanhua(gapp.Id(), id, czl)
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(gapp.Id())
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
	ok, msg := gapp.OptSrv.PetSrv.SSEvolution(gapp.Id(), id, fzid)
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
	gapp.JSONDATAOK("", gin.H{"path": gapp.OptSrv.PetSrv.SSZhuanshengInfo(gapp.Id(), id)})

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
	ok, msg := gapp.OptSrv.PetSrv.SSZhuanSheng(gapp.Id(), id, toid, apid, bpid)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
