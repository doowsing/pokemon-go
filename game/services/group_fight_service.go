package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/psampaz/slice"
	"github.com/unknwon/com"
	"math/rand"
	common2 "pokemon/common"
	"pokemon/common/rcache"
	"pokemon/common/rpc-client/rpc-group"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/services/group-helper"
	"pokemon/game/utils"
	"strings"
	"time"
)

func (fs *FightService) StartGroupFight(userId, multiple int, useSj bool) (fightInfo gin.H, msg string) {
	fightInfo = gin.H{
		"result":               false,
		"waittime":             10,
		"is_fb":                false,
		"fb_need_sj":           0,
		"fb_to_card":           false,
		"fb_to_boss_card":      false,
		"auto_start":           false,
		"users":                nil,
		"leader_id":            0,
		"interval_auto_attack": 5,
		"mapid":                0,
		"multiple":             0,
		"gpcs":                 nil,
		"pet":                  nil,
		"gpc":                  nil,
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		msg = "未组队不可进入！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "未组队不可进入！"
		return
	}
	if userId != group.Leader {
		msg = "不是队长不可进入战斗"
		return
	}
	if len(group.Member) < 2 {
		msg = "队伍至少二人才可以参与战斗！"
		return
	}
	if !group.AllReady() {
		msg = "有队员未准备！"
		return
	}
	mapId := rcache.GetInMap(userId)
	if mapId == 0 {
		msg = "地图不存在！"
		return
	}
	var mmap *models.Map
	if !utils.IsNormalMap(mapId) && !((utils.IsMultipleMap(mapId) || utils.IsYiWangMap(mapId)) && multiple < 3) {
		msg = "该地图不可组队战斗！"
		return
	}
	mmap = common.GetMMap(mapId + multiple)
	if mmap == nil {
		msg = "地图不存在！"
		return
	}
	if group.MapId != mapId {
		msg = "所在地图与组队地图不一致！"
		return
	}
	defer group.Save()
	defer func() {
		_ = msg
		common2.GroupStartFight(group.UUId, gin.H{"code": 200, "msg": msg, "data": fightInfo})
	}()

	autoFlag := rcache.GetAutoFightFlag(group.Leader)
	if autoFlag > 0 {
		fightInfo["auto_start"] = true
	}
	fightInfo["mapid"] = mapId
	now := utils.NowUnix()
	if t := rcache.GetFightCoolTime(group.Leader, now); t > 0 {
		fightInfo["waittime"] = t
		msg = "战斗等待中！"
		return
	}
	if autoFlag == rcache.AutoFightJb {
		if !fs.OptSvc.UserSrv.DecreaseAutoJb(group.Leader) {
			rcache.DelAutoFightFlag(group.Leader)
		}
	} else if autoFlag == rcache.AutoFightYb {
		if !fs.OptSvc.UserSrv.DecreaseAutoYb(group.Leader) {
			rcache.DelAutoFightFlag(group.Leader)
		}
	}

	if autoFlag > 0 {
		rpc_group.SetFightUser(group, group.Leader)
	}

	if utils.IsYiWangMap(mapId) {
		// 副本处理
		// 检测是否可以打
		fightInfo["is_fb"] = true
		if group.FbNeverStart() {
			// 检测是否需要水晶进入
			if group.FbNeedSj() {
				if !useSj {
					msg = "您的队伍有人今日已参与过遗忘副本，队长需支付500水晶，是否支付？"
					fightInfo["need_sj"] = 500
					return
				} else {
					if !fs.OptSvc.UserSrv.DecreaseSj(userId, 500) {
						msg = "您的水晶不足！"
						return
					}
				}
			}
			rpc_group.SetMultiple(group, multiple)
		}
		if group.Level == 2 {
			if len(group.Gpc) <= group.GpcIndex {
				fightInfo["fb_to_boss_card"] = true
				return
			}
			group.GpcIndex++
		} else {
			if group.Level > 0 && group.Process == 0 {
				// 打完了一层怪物
				if group.StartCardTime == 0 || now-group.StartCardTime < 30 {
					fightInfo["fb_to_card"] = true
					return
				}
			}

			if group.Process == 0 {
				rpc_group.SetFbStartCardTime(group, 0)
			}

			if len(group.Gpc) == 0 {
				gpcGroup := &models.GpcGroup{}
				fs.GetDb().Where("map_id=? and step_id=? and group_id=", mapId+multiple, group.Level+1, group.Process+1).First(gpcGroup)
				if len(gpcGroup.GpcList) == 0 {
					msg = "地图配置不存在！请联系管理员"
					return
				}
				group.Gpc = gpcGroup.GpcList
			}
		}

		gpcData := []gin.H{}
		for i, id := range group.Gpc {
			gpc := common.GetGpc(id)
			if gpc == nil {
				msg = "怪物配置不存在！请联系管理员"
				return
			}
			gpcData = append(gpcData, gin.H{
				"id":    gpc.ID,
				"name":  gpc.Name,
				"level": gpc.Level,
			})
			if i == 0 {
				fightInfo["gpc"] = gin.H{
					"id":          gpc.ID,
					"name":        gpc.Name,
					"level":       gpc.Level,
					"hp":          gpc.Hp,
					"stand_img":   gpc.ImgStand,
					"attack_img1": gpc.ImgAck,
					"attack_img2": gpc.ImgAck,
				}
			}
		}
		fightInfo["gpcs"] = gpcData
	} else {
		gpcNum := utils.RandInt(1, int(float64(len(group.Member))*1.5))
		gpcIds := slice.ShuffleInt(slice.CopyInt(mmap.GpcIds))[:gpcNum]
		if len(gpcIds) == 0 {
			msg = "地图配置不存在！请联系管理员"
			return
		}
		gpcData := []gin.H{}
		for i, id := range gpcIds {
			gpc := common.GetGpc(id)
			if gpc == nil {
				msg = "怪物配置不存在！请联系管理员"
				return
			}
			gpcData = append(gpcData, gin.H{
				"id":    gpc.ID,
				"name":  gpc.Name,
				"level": gpc.Level,
			})
			if i == 0 {
				fightInfo["gpc"] = gin.H{
					"id":          gpc.ID,
					"name":        gpc.Name,
					"level":       gpc.Level,
					"hp":          gpc.Hp,
					"stand_img":   gpc.ImgStand,
					"attack_img1": gpc.ImgAck,
					"attack_img2": gpc.ImgAck,
				}
			}
		}
		fightInfo["gpcs"] = gpcData
	}
	userDatas := []gin.H{}
	for _, m := range group.Member {
		userDatas = append(userDatas, gin.H{"id": m.Id, "nickname": m.Nickname})
	}
	if group.FightUserId == 0 {
		rpc_group.SetFightUser(group, group.Leader)
	}

	rcache.SetFightTime(group.Leader, now)
	user := fs.OptSvc.UserSrv.GetUserById(group.FightUserId)
	mainPet := fs.OptSvc.PetSrv.GetPet(group.FightUserId, user.Mbid)
	mainPet.GetM()
	if utils.IsSsMap(mapId) {
		if mainPet.MModel.Wx != 7 {
			msg = "该地图只有神圣宠物可进入！"
			return
		}
	}
	if mmap.CzlProp != "" && mmap.CzlProp != "0" {
		czlItems := strings.Split(mmap.CzlProp, "|")
		if int(mainPet.CC) < com.StrTo(czlItems[0]).MustInt() {
			msg = "准入成长不足！"
			return
		}
	}
	fs.GetZbAttr(mainPet, nil)
	skills := fs.OptSvc.PetSrv.GetPetSkill(mainPet.ID)
	skillData := []gin.H{}
	for _, skill := range skills {
		skill.GetM()
		if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
			continue
		}
		skillData = append(skillData, gin.H{
			"id":   skill.ID,
			"name": skill.MModel.Name,
		})
	}

	petStatus := rcache.GetPetStatus(mainPet.ID)
	if petStatus == nil {
		petStatus = &rcache.PetStatus{}
		rcache.SetPetStatus(mainPet.ID, 0, 0)
	} else {
		if autoFlag == rcache.AutoFightJb {
			petStatus.DeHp = 0
			petStatus.DeMp = 0
			rcache.SetPetStatus(mainPet.ID, 0, 0)
		} else if autoFlag == rcache.AutoFightYb {
			petStatus.DeHp = 0
			petStatus.DeMp = 0
			rcache.SetPetStatus(mainPet.ID, 0, 0)
		}
	}
	petHp := mainPet.ZbAttr.Hp - petStatus.DeHp
	petMp := mainPet.ZbAttr.Mp - petStatus.DeMp
	if petHp < 1 || petMp < 1 {
		if petHp < 1 {
			petHp = 1
		}
		if petMp < 1 {
			petMp = 0
		}
		rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-petHp, mainPet.ZbAttr.Mp-petMp)
	}

	// 重置本次战斗信息
	rpc_group.ResetFightInfo(group)

	fightInfo["interval_auto_attack"] = 10 - int(mainPet.ZbAttr.Special["time"])
	fightInfo["pet"] = gin.H{
		"id":          mainPet.ID,
		"name":        mainPet.MModel.Name,
		"level":       mainPet.Level,
		"hp":          petHp,
		"max_hp":      mainPet.ZbAttr.Hp,
		"mp":          petMp,
		"max_mp":      mainPet.ZbAttr.Mp,
		"exp":         mainPet.NowExp,
		"max_exp":     mainPet.LExp,
		"header_img":  mainPet.MModel.ImgHead,
		"stand_img":   mainPet.MModel.ImgStand,
		"attack_img1": mainPet.MModel.ImgAck,
		"attack_img2": mainPet.MModel.ImgEffect,
		"skills":      skillData,
	}
	fightInfo["users"] = userDatas
	fightInfo["leader_id"] = group.Leader
	fightInfo["mapid"] = mapId
	fightInfo["multiple"] = multiple
	fightInfo["result"] = true
	return
}

func (fs *FightService) GroupAttack(userId, skillId int) (fightInfo gin.H, msg string) {
	msg = ""
	result := gin.H{
		"finish_msg": "",
		"award":      []gin.H{},
	}
	fightInfo = gin.H{
		"success":  false, // 攻击是否成功：缺蓝或技能不对，则攻击失效
		"finish":   false, // 战斗是否结束
		"result":   result,
		"now_user": 0,
		"pet":      nil,
		"new_pet":  nil,
		"gpc":      nil,
		"new_gpc":  nil,
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		fightInfo["finish"] = true
		result["finish_msg"] = "未组队不可战斗！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		fightInfo["finish"] = true
		result["finish_msg"] = "未组队不可战斗！"
		return
	}
	if userId != group.FightUserId {
		fightInfo["success"] = false
		result["finish_msg"] = "未轮到你战斗！"
	}
	if len(group.Member) < 2 {
		fightInfo["finish"] = true
		result["finish_msg"] = "队伍至少二人才可以参与战斗！"
		return
	}
	defer group.Save()
	defer common2.GroupAttack(group.UUId, gin.H{"code": 200, "msg": msg, "data": fightInfo})
	mapId := rcache.GetInMap(userId)
	fightInfo["mapid"] = mapId
	fightInfo["multiple"] = group.Multiple

	now := utils.NowUnix()
	if rcache.GetAttackCoolTime(userId, now) > 0 {
		fightInfo["finish"] = true
		return
	}
	if now-rcache.GetFightTime(userId) > 20 {
		fightInfo["finish"] = true
		return
	}

	user := fs.OptSvc.UserSrv.GetUserById(userId)
	mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()
	fs.GetZbAttr(mainPet, nil)

	petStatus := rcache.GetPetStatus(mainPet.ID)
	if petStatus == nil {
		petStatus = &rcache.PetStatus{}
	}
	petHp := mainPet.ZbAttr.Hp - petStatus.DeHp
	petMp := mainPet.ZbAttr.Mp - petStatus.DeMp

	gpc := group.GetGpc()
	if gpc == nil {
		fightInfo["finish"] = true
		result["finish_msg"] = "怪物信息不存在，战斗已结束！"
		return
	}

	gpcHp := gpc.Hp - group.GpcDeHp
	if petHp < 1 || gpcHp < 1 {
		fightInfo["finish"] = true
		result["finish_msg"] = "怪物死亡，战斗已结束！"
		return
	}
	var skill *models.Uskill
	if skillId == 0 {
		// 可能是自动攻击，查看自动攻击技能
		skillId = rcache.GetAutoFightSkill(userId)
	}
	if skillId == 0 {
		// 默认取普通攻击
		skill = fs.OptSvc.PetSrv.GetSkillBySid(mainPet.ID, 1)
	} else {
		skill = fs.OptSvc.PetSrv.GetSkill(skillId)
	}

	if skill == nil || skill.Bid != mainPet.ID {
		fightInfo["success"] = false
		msg = "技能无效！"
		return
	}
	skill.GetM()
	if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
		fightInfo["success"] = false
		msg = "技能无效！"
		return
	}
	if skill.Ump > petMp {
		fightInfo["success"] = false
		msg = "蓝量不足以发动技能！"
		return
	}

	gpcSkill := fs.GetGpcRandSkill(gpc)
	if gpcSkill == nil {
		fightInfo["finish"] = true
		result["finish_msg"] = "怪物攻击失效，战斗错误！"
		return
	}

	// 开始计算伤害
	success, apResult, bpResult := fs.GetFightResult(mainPet, skill, petStatus, gpc, gpcSkill, group.GpcDeHp)
	if !success {
		fightInfo["success"] = false
		msg = "蓝量不足以发动技能！"
		return
	}
	if apResult == nil || bpResult == nil {
		fightInfo["success"] = false
		msg = "技能类型出错！"
		return
	}
	fightInfo["success"] = true
	fightInfo["pet"] = apResult
	fightInfo["gpc"] = bpResult
	mmap := common.GetMMap(mapId)
	if bpResult.Die {
		// 战斗胜利

		// 累积本层战斗奖励
		rpc_group.AddMoney(group, gpc.Money)
		rpc_group.AddExp(group, gpc.Exp)
		rpc_group.AddAwardProp(group.UUId, gpc.Drops)

		// 看看是否有下一只怪物
		group.GpcIndex += 1
		if newGpc := group.GetGpc(); newGpc != nil {
			// 取怪物信息
			fightInfo["new_gpc"] = gin.H{
				"id":          newGpc.ID,
				"name":        newGpc.Name,
				"level":       newGpc.Level,
				"hp":          newGpc.Hp,
				"stand_img":   newGpc.ImgStand,
				"attack_img1": newGpc.ImgAck,
				"attack_img2": newGpc.ImgAck,
			}
		} else {
			// 结算奖励
			result["finish"] = true
			getResults := fs.processGroupFightVictory(group)
			// 战斗胜利处理
			result["award"] = getResults
			if utils.IsYiWangMap(mapId) {
				if group.Level < 2 {
					rpc_group.SetFbProcess(group, group.Process+1)
					if group.Process == 5 {
						result["fb_process"] = fmt.Sprintf("第%d关通关", group.Level+1)
						rpc_group.SetFbLevel(group, group.Level+1)
						rpc_group.SetFbProcess(group, 0)
						rpc_group.SetGpcs(group, []int{})
						rpc_group.SetGpcIndex(group, 0)
						rpc_group.ResetCardAwards(group)
					} else {

						result["fb_process"] = fmt.Sprintf("%d关%d组", group.Level+1, group.Process+1)
					}
				}
			}

		}
	} else if apResult.Die {
		// 战斗失败，切换战斗者
		if rpc_group.SetNextFightUserId(group) {
			// 存在下一个队员
			fightInfo["new_user"] = group.FightUserId
			user := fs.OptSvc.UserSrv.GetUserById(group.FightUserId)
			new_mainPet := fs.OptSvc.PetSrv.GetPet(group.FightUserId, user.Mbid)
			new_mainPet.GetM()
			if utils.IsSsMap(mapId) {
				if new_mainPet.MModel.Wx != 7 {
					fightInfo["finish"] = true
					result["finish_msg"] = "队员主宠为非神圣宠物，不可在此战斗！"
					return
				}
			}
			if mmap.CzlProp != "" && mmap.CzlProp != "0" {
				czlItems := strings.Split(mmap.CzlProp, "|")
				if int(new_mainPet.CC) < com.StrTo(czlItems[0]).MustInt() {
					fightInfo["finish"] = true
					result["finish_msg"] = "队员主宠成长不足，不可在此战斗！"
					return
				}
			}
			fs.GetZbAttr(new_mainPet, nil)
			skills := fs.OptSvc.PetSrv.GetPetSkill(new_mainPet.ID)
			skillData := []gin.H{}
			for _, skill := range skills {
				skill.GetM()
				if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
					continue
				}
				skillData = append(skillData, gin.H{
					"id":   skill.ID,
					"name": skill.MModel.Name,
				})
			}

			petStatus := rcache.GetPetStatus(new_mainPet.ID)
			if petStatus == nil {
				petStatus = &rcache.PetStatus{}
				rcache.SetPetStatus(new_mainPet.ID, 0, 0)
			} else {
				autoFlag := rcache.GetAutoFightFlag(group.Leader)
				if autoFlag == rcache.AutoFightJb {
					petStatus.DeHp = 0
					petStatus.DeMp = 0
					rcache.SetPetStatus(new_mainPet.ID, 0, 0)
				} else if autoFlag == rcache.AutoFightYb {
					petStatus.DeHp = 0
					petStatus.DeMp = 0
					rcache.SetPetStatus(new_mainPet.ID, 0, 0)
				}
			}
			petHp := new_mainPet.ZbAttr.Hp - petStatus.DeHp
			petMp := new_mainPet.ZbAttr.Mp - petStatus.DeMp
			if petHp < 1 || petMp < 1 {
				if petHp < 1 {
					petHp = 1
				}
				if petMp < 1 {
					petMp = 0
				}
				rcache.SetPetStatus(new_mainPet.ID, new_mainPet.ZbAttr.Hp-petHp, new_mainPet.ZbAttr.Mp-petMp)
			}

			fightInfo["interval_auto_attack"] = 10 - int(new_mainPet.ZbAttr.Special["time"])
			fightInfo["new_pet"] = gin.H{
				"id":          new_mainPet.ID,
				"name":        new_mainPet.MModel.Name,
				"level":       new_mainPet.Level,
				"hp":          petHp,
				"max_hp":      new_mainPet.ZbAttr.Hp,
				"mp":          petMp,
				"max_mp":      new_mainPet.ZbAttr.Mp,
				"exp":         new_mainPet.NowExp,
				"max_exp":     new_mainPet.LExp,
				"header_img":  new_mainPet.MModel.ImgHead,
				"stand_img":   new_mainPet.MModel.ImgStand,
				"attack_img1": new_mainPet.MModel.ImgAck,
				"attack_img2": new_mainPet.MModel.ImgEffect,
				"skills":      skillData,
			}
		} else {
			// 均已战败，本轮战斗失败
			result["finish"] = true
			result["finish_msg"] = "战斗失败！"

			// 战斗失败处理
			hasNext, finishMsg := fs.processFightDefeat(user, mmap)
			result["finish_msg"] = finishMsg
			result["has_next"] = hasNext
			if group.Level == 2 {
				result["has_next"] = true
			}
		}

	}
	rcache.SetFightTime(group.Leader, now)
	rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-apResult.Hp, mainPet.ZbAttr.Mp-apResult.Mp)
	rpc_group.SetGpcDeHp(group, gpc.Hp-bpResult.Hp)
	return fightInfo, msg

}

func (fs *FightService) processGroupFightVictory(group *group_helper.UserGroup) []gin.H {

	fs.OptSvc.Begin()
	defer fs.OptSvc.Commit()
	autoFlag := rcache.GetAutoFightFlag(group.Leader)
	exp := group.GetExp
	if autoFlag == rcache.AutoFightJb {
		exp = int(float64(exp) * 1.2)
	} else if autoFlag == rcache.AutoFightYb {
		exp = int(float64(exp) * 1.5)
	}

	start := false
	for _, m := range group.Member {
		if m.Id == group.FightUserId {
			start = true
		}
		if start {

		}
	}
	results := []gin.H{}
	for _, m := range group.Member {
		getExp := exp
		getMoney := group.GetMoney
		if m.Id == group.FightUserId {
			getExp = int(float64(getExp) * 1.2)
			getMoney = int(float64(getMoney) * 1.2)
		}
		user := fs.OptSvc.UserSrv.GetUserById(m.Id)
		pet := fs.OptSvc.PetSrv.GetPetById(user.Mbid)
		petUpgrade := fs.OptSvc.PetSrv.IncreaseExp2Pet(pet, getExp)
		fs.OptSvc.UserSrv.IncreaseJb(m.Id, group.GetMoney)

		propNameList := []string{}
		getProps := group.GetProps[m.Id]
		carryPropCnt := fs.OptSvc.PropSrv.GetCarryPropsCnt(user.ID)
		if carryPropCnt < user.BagPlace && len(getProps) > 0 {
			propNameList = fs.OptSvc.PropSrv.AddPropByPids(user, carryPropCnt, getProps)
		}
		results = append(results, gin.H{
			"id":          m.Id,
			"nickname":    m.Nickname,
			"jb":          getMoney,
			"prop":        propNameList,
			"exp":         getExp,
			"pet_upgrade": petUpgrade,
		})
	}
	return results
}

func (fs *FightService) GetGroupCardData(userId int) (result gin.H, msg string) {
	result = gin.H{
		"result": false,
		"cards":  nil,
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		msg = "未组队不可进入！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "未组队不可进入！"
		return
	}
	now := time.Now()
	if !(group.Level > 0 && group.Process == 0 && (group.StartCardTime == 0 || int(now.Unix())-group.StartCardTime < 30)) {
		msg = "还没到翻牌的时间！"
		return
	}
	if len(group.CardAwards) == 0 {

		defer common2.GroupEnterCard(group.UUId)

		if userId != group.Leader {
			msg = "不是队长进入翻牌！"
			return
		}

		rpc_group.SetFbStartCardTime(group, int(now.Unix()))
		cards := common.GetYiwangCards(group.Multiple, false)
		sjCards := common.GetYiwangCards(group.Multiple, true)
		r := rand.New(rand.NewSource(now.UnixNano()))
		r.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
		r.Shuffle(len(sjCards), func(i, j int) {
			sjCards[i], sjCards[j] = sjCards[j], sjCards[i]
		})
		CardAwards := []*group_helper.CardInfo{}
		for _, a := range cards[:5] {
			CardAwards = append(CardAwards, &group_helper.CardInfo{
				TarotCardId: a.Id,
				Content:     "",
				UserId:      0,
			})
		}
		for _, a := range sjCards[:5] {
			CardAwards = append(CardAwards, &group_helper.CardInfo{
				TarotCardId: a.Id,
				Content:     "",
				UserId:      0,
			})
		}
		rpc_group.SetFbCardAwards(group, CardAwards)
		result["result"] = true
		return
	} else {
		cardDatas := []gin.H{}
		for i, info := range group.CardAwards {
			if info.UserId > 0 {
				cardDatas = append(cardDatas, gin.H{"id": i, "user": info.Nickname, "content": info.Content})
			}
		}
		result["result"] = true
		result["cards"] = cardDatas
		return
	}
}

func (fs *FightService) DoGroupCard(userId, position int) (result gin.H, msg string) {
	result = gin.H{
		"result":        false,
		"position":      0,
		"nickname":      "",
		"award_content": "",
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		msg = "未组队不可进入！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "未组队不可进入！"
		return
	}
	now := time.Now()
	if !(group.Level > 0 && group.Process == 0 && (group.StartCardTime == 0 || int(now.Unix())-group.StartCardTime < 30)) {
		msg = "还没到翻牌的时间！"
		return
	}
	if len(group.CardAwards) == 0 {
		msg = "翻牌未开始"
		return
	}
	for i, c := range group.CardAwards {
		if (position >= 5 && i >= 5) || (position < 5 && i < 5) {
			if c.UserId == userId {
				msg = "您已翻过该类牌！"
				return
			}
		}
	}
	info := group.CardAwards[position]
	if info.UserId > 0 {
		msg = "该牌已被其他玩家翻过！"
		return
	}
	if position >= 5 {
		if !fs.OptSvc.UserSrv.DecreaseSj(userId, 100) {
			msg = "您的水晶不足！"
			return
		}
	}
	result["position"] = position
	result["result"] = true

	m := group.GetMember(userId)
	info.UserId = userId
	info.Nickname = m.Nickname

	defer common2.GroupUpdateCard(group.UUId, gin.H{"code": 200, "msg": msg, "data": gin.H{"position": position, "nickname": info.Nickname, "content": info.Content}})
	effect := strings.SplitN(common.GetYiwangCard(info.TarotCardId).Effect, ":", 2)
	switch effect[0] {
	case "giveitems":
		for _, s := range strings.Split(effect[1], ",") {
			items := strings.Split(s, ":")
			if len(items) >= 4 {
				pid := com.StrTo(items[0]).MustInt()
				num := com.StrTo(items[1]).MustInt()
				rate := com.StrTo(items[2]).MustInt()
				annouce := com.StrTo(items[3]).MustInt()
				if rand.Intn(rate) == 0 {
					fs.OptSvc.PropSrv.AddOrCreateProp(userId, pid, num, true)
					result["nickname"] = m.Nickname
					mprop := common.GetMProp(pid)
					content := fmt.Sprintf("%s %d 个", mprop.Name, num)
					result["award_content"] = content
					if annouce == 1 {
						AnnounceAll(m.Nickname, fmt.Sprintf("获得遗忘宫殿第%d关的奖励：%s", group.Level, content))
					}
					info.Content = content
					return
				}
			}
		}

		break
	case "exp_add":
		num := com.StrTo(effect[1]).MustInt()
		fs.OptSvc.PetSrv.IncreaseExp2MainPet(userId, num)
		content := fmt.Sprintf("经验 %d", num)
		info.Content = content
		result["nickname"] = m.Nickname
		result["award_content"] = content
		return
	}
	info.Content = "无！"
	rpc_group.SetFbCardAward(group, position, info)
	return
}

func (fs *FightService) GetAllGroupCardData(userId int) (result gin.H, msg string) {
	result = gin.H{
		"result": false,
		"cards":  nil,
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		msg = "未组队不可进入！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "未组队不可进入！"
		return
	}
	now := time.Now()
	if !(group.Level > 0 && group.Process == 0 && (group.StartCardTime == 0 || int(now.Unix())-group.StartCardTime < 30)) {
		msg = "还没到翻牌的时间！"
		return
	}
	if len(group.CardAwards) == 0 {
		msg = "还没到翻牌的时间！"
		return
	}
	defer common2.GroupEnterCard(group.UUId)
	cardDatas := []gin.H{}
	for i, info := range group.CardAwards {
		data := gin.H{"id": i, "user": "", "content": info.Content}
		if info.UserId > 0 {
			data["user"] = info.Nickname
		} else {
			effect := strings.SplitN(common.GetYiwangCard(info.TarotCardId).Effect, ":", 2)
			switch effect[0] {
			case "giveitems":
				find := false
				for _, s := range strings.Split(effect[1], ",") {
					items := strings.Split(s, ":")
					if len(items) >= 4 {
						pid := com.StrTo(items[0]).MustInt()
						num := com.StrTo(items[1]).MustInt()
						rate := com.StrTo(items[2]).MustInt()
						if rand.Intn(rate) == 0 {
							find = true
							mprop := common.GetMProp(pid)
							content := fmt.Sprintf("%s %d 个", mprop.Name, num)
							info.Content = content
							break
						}
					}
				}
				if !find {
					info.Content = "无！"
				}
				break
			case "exp_add":
				num := com.StrTo(effect[1]).MustInt()
				content := fmt.Sprintf("经验 %d", num)
				info.Content = content
				break
			default:
				info.Content = "无！"
				break
			}

		}
		cardDatas = append(cardDatas, data)
	}
	result["result"] = true
	result["cards"] = cardDatas
	return
}

func (fs *FightService) GetGroupBossCard(userId int) (result gin.H, msg string) {
	result = gin.H{
		"result": false,
		"cards":  nil,
		"leader": 0,
		"member": nil,
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		msg = "未组队不可进入！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "未组队不可进入！"
		return
	}
	if group.End {
		msg = "本次副本已结束，队伍自动解散！"
		rpc_group.DropGroup(group.UUId, group.Leader)
		return
	}
	now := time.Now()
	if group.Level != 2 || group.StartCardTime == 0 || int(now.Unix())-group.StartCardTime < 30 {
		msg = "还没到翻牌的时间！"
		return
	}
	result["leader"] = group.Leader
	result["cards"] = group.BossCardEffect
	memberDatas := []gin.H{}
	for _, m := range group.Member {
		user := fs.OptSvc.UserSrv.GetUserById(m.Id)
		mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
		fs.GetZbAttr(mainPet, nil)
		petStatus := rcache.GetPetStatus(mainPet.ID)
		petHp := mainPet.ZbAttr.Hp - petStatus.DeHp
		petMp := mainPet.ZbAttr.Mp - petStatus.DeMp
		if petHp < 1 || petMp < 1 {
			if petHp < 1 {
				petHp = 1
			}
			if petMp < 1 {
				petMp = 0
			}
			rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-petHp, mainPet.ZbAttr.Mp-petMp)
		}
		memberDatas = append(memberDatas, gin.H{
			"id":         m.Id,
			"nickname":   m.Nickname,
			"img":        fmt.Sprintf("face%s.gif", user.Headimg),
			"pet_level":  mainPet.Level,
			"pet_hp":     petHp,
			"pet_max_hp": mainPet.Hp,
			"pet_mp":     petMp,
			"pet_max_mp": mainPet.Mp,
		})
	}
	common2.GroupEnterBossCard(group.UUId)
	result["member"] = memberDatas
	return
}

func (fs *FightService) DoGroupBossCard(userId, position int) (result gin.H, msg string) {
	result = gin.H{
		"result":      false,
		"enter_fight": false,
	}
	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		msg = "未组队不可进入！"
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "未组队不可进入！"
		return
	}
	now := time.Now()
	if group.Level != 2 || group.StartCardTime == 0 || int(now.Unix())-group.StartCardTime < 30 {
		msg = "还没到翻牌的时间！"
		return
	}
	for _, c := range group.BossCardEffect {
		if c.Position == position {
			msg = "此牌已翻过！"
			return
		}
	}
	bossCards := common.GetYiwangBossCards(group.Multiple)
	card := bossCards[rand.Intn(len(bossCards))]

	result["result"] = true
	effect := strings.SplitN(card.Effect, ":", 2)
	switch effect[0] {
	case "all_giveitems":
		for _, m := range group.Member {
			for _, s := range strings.Split(effect[1], ",") {
				items := strings.Split(s, ":")
				if len(items) >= 3 {
					pid := com.StrTo(items[0]).MustInt()
					num := com.StrTo(items[1]).MustInt()
					rate := com.StrTo(items[2]).MustInt()
					if rand.Intn(rate) == 0 {
						fs.OptSvc.PropSrv.AddOrCreateProp(m.Id, pid, num, true)
						mprop := common.GetMProp(pid)
						content := fmt.Sprintf("获得 %s %d 个", mprop.Name, num)
						common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": content})
						break
					}
				}
			}
		}

		break
	case "all_hp_less", "all_hp_add":
		ratio := com.StrTo(strings.ReplaceAll(effect[1], "%", "")).MustFloat64()
		memberDatas := []gin.H{}
		for _, m := range group.Member {
			user := fs.OptSvc.UserSrv.GetUserById(m.Id)
			mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
			fs.GetZbAttr(mainPet, nil)
			petStatus := rcache.GetPetStatus(mainPet.ID)
			petHp := mainPet.ZbAttr.Hp - petStatus.DeHp
			petMp := mainPet.ZbAttr.Mp - petStatus.DeMp
			if effect[0] == "all_hp_less" {
				petHp = petHp - int(float64(petHp)*ratio/100)
			} else {
				petHp = petHp + int(float64(petHp)*ratio/100)
			}

			if petHp < 1 || petMp < 0 {
				if petHp < 1 {
					petHp = 1
				}
				if petMp < 0 {
					petMp = 0
				}
			}
			if petHp > mainPet.ZbAttr.Hp {
				petHp = mainPet.ZbAttr.Hp
			}
			rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-petHp, mainPet.ZbAttr.Mp-petMp)
			memberDatas = append(memberDatas, gin.H{
				"id":         m.Id,
				"nickname":   m.Nickname,
				"img":        fmt.Sprintf("face%s.gif", user.Headimg),
				"pet_level":  mainPet.Level,
				"pet_hp":     petHp,
				"pet_max_hp": mainPet.ZbAttr.Hp,
				"pet_mp":     petMp,
				"pet_max_mp": mainPet.ZbAttr.Mp,
			})
		}
		var content string
		if effect[0] == "all_hp_less" {
			content = fmt.Sprintf("全体减少HP %s", effect[1])
		} else {
			content = fmt.Sprintf("全体增加HP %s", effect[1])
		}
		common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": content})
		result["member"] = memberDatas
		break
	case "all_money":
		num := com.StrTo(effect[1]).MustInt()
		for _, m := range group.Member {
			fs.OptSvc.UserSrv.IncreaseJb(m.Id, num)
		}
		content := fmt.Sprintf("全体增加金币 %d", num)
		common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": content})
		break
	case "all_fight":
		// 遇到怪物，进入战斗
		gpcId := com.StrTo(effect[1]).MustInt()
		rpc_group.SetGpcs(group, append(group.Gpc, gpcId))
		common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": "遇到怪物，开始战斗！"})
		result["enter_fight"] = true
		break
	case "hit_one":
		newMembers := []*group_helper.MemberStatus{}
		member := group.Member[rand.Intn(len(group.Member)-1)+1]
		var content string
		for _, m := range group.Member {
			if m.Id == member.Id {
				content = fmt.Sprintf("遇上恶魔，%s 被强制踢出副本。", m.Nickname)
			} else {
				newMembers = append(newMembers, m)
			}
		}
		rpc_group.SetMemberStatus(group, newMembers)
		memberDatas := []gin.H{}
		for _, m := range group.Member {
			user := fs.OptSvc.UserSrv.GetUserById(m.Id)
			mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
			fs.GetZbAttr(mainPet, nil)
			petStatus := rcache.GetPetStatus(mainPet.ID)
			petHp := mainPet.ZbAttr.Hp - petStatus.DeHp
			petMp := mainPet.ZbAttr.Mp - petStatus.DeMp
			if petHp < 1 || petMp < 0 {
				if petHp < 1 {
					petHp = 1
				}
				if petMp < 0 {
					petMp = 0
				}
				rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-petHp, mainPet.ZbAttr.Mp-petMp)
			}
			memberDatas = append(memberDatas, gin.H{
				"id":         m.Id,
				"nickname":   m.Nickname,
				"img":        fmt.Sprintf("face%s.gif", user.Headimg),
				"pet_level":  mainPet.Level,
				"pet_hp":     petHp,
				"pet_max_hp": mainPet.Hp,
				"pet_mp":     petMp,
				"pet_max_mp": mainPet.Mp,
			})
		}
		result["member"] = memberDatas
		common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": content, "member": memberDatas})

		if len(newMembers) == 1 {
			common2.NoticeTips(group.Leader, "副本人数不足，强制结束副本！")
			rpc_group.DropGroup(group.UUId, group.Leader)
		}
		break
	case "hit_all":
		// 解散队伍
		common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": "运气太差，遇上恶魔，你们将被强制t出副本，请明日再来吧，挑战副本失败!"})
		rpc_group.DropGroup(group.UUId, group.Leader)
		break
	case "all_exp_add":
		num := com.StrTo(effect[1]).MustInt()
		for _, m := range group.Member {
			fs.OptSvc.PetSrv.IncreaseExp2MainPet(m.Id, num)
		}
		common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": fmt.Sprintf("全体获得%d点经验！", num)})

		break
	case "":
		gpcId := com.StrTo(card.Boss).MustInt()
		if gpcId > 0 {
			rpc_group.SetGpcs(group, append(group.Gpc, gpcId))
			rpc_group.SetFbEnd(group, true)
			common2.GroupUpdateBossCard(group.UUId, gin.H{"img": card.Img, "position": position, "award_content": "遇到Boss，开始战斗！"})
			result["enter_fight"] = true
		}
	default:
		result["result"] = false
		msg = "卡片设定错误"
		return
	}
	rpc_group.AddFbBossCardEffects(group, &group_helper.CardEffectInfo{
		Position: position,
		Img:      card.Img,
	})
	return
}

func (fs *FightService) SetGroupStatus(userId, status int) (bool, string) {

	groupId := rpc_group.GetGroupID(userId)
	if groupId == "" {
		return false, "您未加入队伍！"
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		return false, "您未加入队伍！"
	}
	if group.IsLeader(userId) {
		return false, "队长不可离队！"
	}
	if status < 0 {
		for _, m := range group.Member {
			if m.Id == userId {
				rpc_group.SetUserStatus(group, userId, !m.Ready)
				return true, "设置成功！"
			}
		}
		return false, "您未加入队伍！"
	} else if status == 0 {
		rpc_group.SetUserStatus(group, userId, false)
	} else {
		rpc_group.SetUserStatus(group, userId, true)
	}
	return true, "设置成功！"
}

func (fs *FightService) SetUserUnReady(userId int, groupId string) {
	group := rpc_group.GetGroup(groupId)
	if group != nil {
		rpc_group.SetUserStatus(group, userId, false)
	}
}
