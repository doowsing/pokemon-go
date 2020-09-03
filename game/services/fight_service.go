package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/psampaz/slice"
	"github.com/unknwon/com"
	"math/rand"
	"pokemon/common/rcache"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"strconv"
	"strings"
)

const NoKeyMapPunishmentRate = 3

type FightService struct {
	BaseService
}

func NewFightService(osrc *OptService) *FightService {
	us := &FightService{}
	us.SetOptSrc(osrc)
	return us
}

func (fs *FightService) DelZbAttr(upid int) {
	rcache.ClearPetAttribute(upid)
}

func (fs *FightService) GetZbAttr(upet *models.UPet, zbs []models.UProp) {
	if speAttr := rcache.GetPetAttribute(upet.ID); speAttr != nil {
		upet.ZbAttr = speAttr
		//upet.Ac = speAttr.Ac
		//upet.Mc = speAttr.Mc
		//upet.Hits = speAttr.Hits
		//upet.Miss = speAttr.Miss
		//upet.Speed = speAttr.Speed
		//upet.Hp = speAttr.Hp
		//upet.Mp = speAttr.Mp
		//fmt.Printf("zbs attrs:%s\n", upet.SpeAttr)
		return
	}
	if zbs == nil {
		zbs = fs.OptSvc.PropSrv.GetPZbs(upet.ID)
	}

	upet.ZbAttr.Ac = upet.Ac
	upet.ZbAttr.Mc = upet.Mc
	upet.ZbAttr.Hits = upet.Hits
	upet.ZbAttr.Miss = upet.Miss
	upet.ZbAttr.Speed = upet.Speed
	upet.ZbAttr.Hp = upet.Hp
	upet.ZbAttr.Mp = upet.Mp
	zbAttr := fs.OptSvc.PropSrv.CountZbAttr(zbs)
	for t, v := range zbAttr {
		switch t {
		case "ac":
			upet.ZbAttr.Ac += int(v)
			break
		case "mc":
			upet.ZbAttr.Mc += int(v)
			break
		case "hp":
			upet.ZbAttr.Hp += int(v)
			break
		case "mp":
			upet.ZbAttr.Mp += int(v)
			break
		case "speed":
			upet.ZbAttr.Speed += int(v)
			break
		case "hits":
			upet.ZbAttr.Hits += int(v)
			break
		case "miss":
			upet.ZbAttr.Miss += int(v)
			break
		case "time", "crit", "dxsh", "sdmp", "hitshp", "hitsmp", "shjs", "szmp":
			upet.ZbAttr.Special[t] = v
			break
		}
	}
	for t, v := range zbAttr {

		if !strings.Contains(t, "rate") {
			continue
		}
		switch strings.ReplaceAll(t, "rate", "") {
		case "ac":
			upet.ZbAttr.Ac = int(float64(upet.ZbAttr.Ac) * (1 + v))
			break
		case "mc":
			upet.ZbAttr.Mc = int(float64(upet.ZbAttr.Mc) * (1 + v))
			break
		case "hp":
			upet.ZbAttr.Hp = int(float64(upet.ZbAttr.Hp) * (1 + v))
			break
		case "mp":
			upet.ZbAttr.Mp = int(float64(upet.ZbAttr.Mp) * (1 + v))
			break
		case "speed":
			upet.ZbAttr.Speed = int(float64(upet.ZbAttr.Speed) * (1 + v))
			break
		case "hits":
			upet.ZbAttr.Hits = int(float64(upet.ZbAttr.Hits) * (1 + v))
			break
		case "miss":
			upet.ZbAttr.Miss = int(float64(upet.ZbAttr.Miss) * (1 + v))
			break
		}
	}

	// 称号
	userInfo := fs.OptSvc.UserSrv.GetUserInfoById(upet.Uid)
	if userInfo.NowAchievementTitle != "" {
		cardTitle := fs.OptSvc.UserSrv.GetCardTile(userInfo.NowAchievementTitle)
		if cardTitle != nil {
			upet.ZbAttr.Ac += cardTitle.Ac
			upet.ZbAttr.Mc += cardTitle.Mc
			upet.ZbAttr.Hits += cardTitle.Hits
			upet.ZbAttr.Miss += cardTitle.Miss
			upet.ZbAttr.Speed += cardTitle.Speed
			upet.ZbAttr.Hp += cardTitle.Hp
			upet.ZbAttr.Mp += cardTitle.Mp
			upet.ZbAttr.Ac = int(float64(upet.ZbAttr.Ac) * (1 + float64(cardTitle.AcRate)*0.01))
			upet.ZbAttr.Mc = int(float64(upet.ZbAttr.Mc) * (1 + float64(cardTitle.McRate)*0.01))
			upet.ZbAttr.Hits = int(float64(upet.ZbAttr.Hits) * (1 + float64(cardTitle.HitsRate)*0.01))
			upet.ZbAttr.Miss = int(float64(upet.ZbAttr.Miss) * (1 + float64(cardTitle.MissRate)*0.01))
			upet.ZbAttr.Speed = int(float64(upet.ZbAttr.Speed) * (1 + float64(cardTitle.SpeedRate)*0.01))
			upet.ZbAttr.Hp = int(float64(upet.ZbAttr.Hp) * (1 + float64(cardTitle.HpRate)*0.01))
			upet.ZbAttr.Mp = int(float64(upet.ZbAttr.Mp) * (1 + float64(cardTitle.MpRate)*0.01))
			upet.ZbAttr.Special["time"] += float64(cardTitle.Time)
			upet.ZbAttr.Special["money"] += float64(cardTitle.AddMoney)
			upet.ZbAttr.Special["dxsh"] += float64(cardTitle.Dxsh)
			upet.ZbAttr.Special["sdmp"] += float64(cardTitle.Sdmp)
			upet.ZbAttr.Special["hitshp"] += float64(cardTitle.HitsHp)
			upet.ZbAttr.Special["hitsmp"] += float64(cardTitle.HitsMp)
			upet.ZbAttr.Special["shjs"] += float64(cardTitle.Shjs)
			upet.ZbAttr.Special["szmp"] += float64(cardTitle.Szmp)
		}
	}
	//fmt.Printf("zbs attrs:%s\n", upet.SpeAttr)
	rcache.SetPetAttribute(upet.ID, upet.ZbAttr)
}

func (fs *FightService) OpenMap(userId, mapId int) (bool, string) {
	user := fs.OptSvc.UserSrv.GetUserById(userId)
	openMap := strings.Split(user.OpenMap, ",")
	for _, v := range openMap {
		if v == strconv.Itoa(mapId) {
			return false, "该地图已经打开了!"
		}
	}
	keyProp := &models.MProp{}
	fs.GetDb().Where("effect=? and varyname=?", "openmap:"+strconv.Itoa(mapId), 13).First(keyProp)
	if keyProp.ID == 0 || !fs.OptSvc.PropSrv.DecrPropByPid(userId, keyProp.ID, 1) {
		return false, "找不到地图对应的钥匙！"
	}
	openMap = append(openMap, strconv.Itoa(mapId))
	fs.GetDb().Model(user).Update(gin.H{"openmap": strings.Join(openMap, ",")})
	return true, "开启地图成功！"
}

func (fs *FightService) GetFbRecord(userId, mapId int) *models.RecordFb {
	record := &models.RecordFb{}
	if fs.GetDb().Where("uid=? and inmap=?", userId, mapId).Find(record).RowsAffected > 0 {
		return record
	}
	return nil
}

func (fs *FightService) GetNormalMapGpc(userId int, mmap *models.Map) *models.Gpc {
	for repeateTimes := 0; repeateTimes < 10; repeateTimes++ {
		randLv := rand.Intn(len(mmap.Lv2GpcIds))
		var Gpcs []int
		i := 0
		for _, list := range mmap.Lv2GpcIds {
			if i == randLv {
				Gpcs = list
				break
			}
			i++
		}

		randId := Gpcs[rand.Intn(len(Gpcs))]
		gpc := common.GetGpc(randId)
		now := utils.NowUnix()
		if gpc.Boss == 3 {
			bossRecord := &models.BossRecord{}
			fs.GetDb().Where("gid=?", gpc.ID).First(bossRecord)
			if bossRecord.Id > 0 {
				if !((bossRecord.Glock == 0 && bossRecord.Dtime+3600 < now) || (bossRecord.Glock == 1 && bossRecord.Rtime+600 < now)) {
					if len(mmap.GpcIds) < 2 {
						return nil
					}
					continue
				}
				fs.GetDb().Model(bossRecord).Update(gin.H{
					"glock":    1,
					"fightuid": userId,
					"rtime":    now,
				})
			} else {
				bossRecord.GpcId = gpc.ID
				bossRecord.Glock = 1
				bossRecord.FightUid = userId
				bossRecord.Rtime = now
				fs.GetDb().Create(bossRecord)
			}
			user := fs.OptSvc.UserSrv.GetUserById(userId)
			AnnounceAll(user.Nickname, fmt.Sprintf("遇上了沉睡中的[%s]，勇士请赶快去消灭它吧！", gpc.Name))
		}

		return gpc
	}
	return nil
}

func (fs *FightService) IntoMap(userId, mapId int) (mapInfo gin.H, msg string) {
	//
	msg = ""
	mapInfo = gin.H{
		"result":      false,
		"id":          0,
		"title":       "",
		"description": "",
		"level":       "",
		"gpcs":        "",
		"is_multiple": false, //是否是多难度的
	}
	mmap := common.GetMMap(mapId)
	if mmap == nil {
		msg = "地图不存在！"
		return
	}
	if !utils.IsNormalMap(mapId) && !utils.IsTTMap(mapId) {
		msg = "地图未开放！"
		return
	}

	user := fs.OptSvc.UserSrv.GetUserById(userId)
	if utils.IsNeedKeyMap(mapId) && !com.IsSliceContainsStr(strings.Split(user.OpenMap, ","), strconv.Itoa(mapId)) {
		msg = "地图未开启！"
		return
	}
	if utils.IsTTMap(mapId) {
		userInfo := fs.OptSvc.UserSrv.GetUserInfoById(userId)
		now := utils.ToDayStartUnix()
		if now > userInfo.TgLastTime {
			// 今日没打过
			fs.GetDb().Model(userInfo).Update(gin.H{"tglasttime": utils.NowUnix(), "tgt": 0})
			userInfo.Tgt = 0
			fs.NewTTRecord(userId, 0)
			rcache.SetTTFlag(userId, rcache.UserTTNone)
		}
		mapInfo["level"] = userInfo.Tgt + 1
	} else {
		mapInfo["level"] = mmap.Level
	}
	rcache.DelAutoFightFlag(userId)
	rcache.DelAutoFightSkill(userId)
	mapInfo["result"] = true
	mapInfo["id"] = mmap.ID
	mapInfo["title"] = mmap.Name
	mapInfo["description"] = mmap.Description
	mapInfo["gpcs"] = mmap.GpcList
	if utils.IsMultipleMap(mapId) {
		mapInfo["is_multiple"] = true
	}
	rcache.SetInMap(userId, mapId)

	return
}

func (fs *FightService) IntoFbMap(userId, mapId int) (mapInfo gin.H, msg string) {
	//
	msg = ""
	mapInfo = gin.H{
		"result":      false,
		"id":          0,
		"title":       "",
		"description": "",
		"level":       "",
		"gpcs":        "",
		"gpc_cnt":     0,
		"next_gpc":    "", //下一只怪物
		"time":        0,  //冷却时间
		"progress":    1,
	}
	mmap := common.GetMMap(mapId)
	if mmap == nil {
		msg = "地图不存在！"
		return
	}
	if !utils.IsFbMap(mapId) {
		msg = "地图未开放！"
		return
	}
	fbSet := utils.GetFbSet(mapId)
	if fbSet == nil {
		msg = "地图未开放！"
		return
	}
	now := utils.NowUnix()
	fbRecord := fs.GetFbRecord(userId, mapId)
	if fbRecord == nil {
		fbRecord = &models.RecordFb{
			Uid:      userId,
			GpcId:    fbSet.Gpcs[0],
			LeftTime: fbSet.CoolTime,
			InMap:    mapId,
			SrcTime:  now,
		}
		if fs.GetDb().Create(fbRecord).RowsAffected < 1 {
			msg = "系统错误！"
			return
		}
	}
	if fbRecord.SrcTime+fbSet.CoolTime < now {
		fbRecord.GpcId = fbSet.Gpcs[0]
		fbRecord.SrcTime = now
		fs.GetDb().Model(fbRecord).Update(gin.H{
			"gwid":    fbRecord.GpcId,
			"srctime": fbRecord.SrcTime,
		})
	}
	if fbRecord.GpcId == 0 {
		mapInfo["time"] = (fbRecord.SrcTime + fbRecord.LeftTime) - now
		fbRecord.GpcId = fbSet.Gpcs[0]
	}
	for i, id := range fbSet.Gpcs {
		if fbRecord.GpcId == id {
			mapInfo["progress"] = i + 1
		}
	}
	rcache.DelAutoFightFlag(userId)
	rcache.DelAutoFightSkill(userId)
	mapInfo["result"] = true
	mapInfo["id"] = mmap.ID
	mapInfo["title"] = mmap.Name
	mapInfo["description"] = mmap.Description
	mapInfo["level"] = mmap.Level
	mapInfo["gpcs"] = mmap.GpcList
	mapInfo["gpc_cnt"] = len(fbSet.Gpcs)
	mapInfo["next_gpc"] = common.GetGpc(fbRecord.GpcId).Name
	rcache.SetInMap(userId, mapId)
	return

}

func (fs *FightService) StartFight(userId, multiple int) (fightInfo gin.H, msg string) {
	msg = ""
	fightInfo = gin.H{
		"result":               false,
		"waittime":             0,
		"is_fb":                false,
		"need_sj":              0,
		"auto_start":           false,
		"user":                 nil,
		"catch_ball":           nil,
		"interval_auto_attack": 10,
		"mapid":                0,
		"multiple":             multiple,
		"pet":                  nil,
		"gpc":                  nil,
	}

	mapId := rcache.GetInMap(userId)
	autoFlag := rcache.GetAutoFightFlag(userId)
	if autoFlag > 0 {
		fightInfo["auto_start"] = true
	}
	fightInfo["mapid"] = mapId
	now := utils.NowUnix()
	if t := rcache.GetFightCoolTime(userId, now); t > 0 {
		fightInfo["waittime"] = t
		msg = "战斗等待中！"
		return
	}
	if autoFlag == rcache.AutoFightJb {
		if !fs.OptSvc.UserSrv.DecreaseAutoJb(userId) {
			rcache.DelAutoFightFlag(userId)
		}
	} else if autoFlag == rcache.AutoFightYb {
		if !fs.OptSvc.UserSrv.DecreaseAutoYb(userId) {
			rcache.DelAutoFightFlag(userId)
		}
	}

	var mmap *models.Map
	var gpc *models.Gpc
	if utils.IsNormalMap(mapId) || utils.IsTTMap(mapId) {
		if utils.IsMultipleMap(mapId) && multiple < 3 {
			mmap = common.GetMMap(mapId + multiple)
			gpc = fs.GetNormalMapGpc(userId, mmap)
		} else if utils.IsTTMap(mapId) {
			mmap = common.GetMMap(mapId)
			var needSj bool
			needSj, gpc = fs.GetTTGpc(userId)
			if needSj {
				fightInfo["need_sj"] = 200
				return
			}
		} else if utils.IsHuPoWuMap(mapId) {

		} else if utils.IsYiWangMap(mapId) {
			msg = "该地图只能组队战斗！"
			return
		} else {
			mmap = common.GetMMap(mapId)
			gpc = fs.GetNormalMapGpc(userId, mmap)
		}

	} else if utils.IsFbMap(mapId) {
		fightInfo["is_fb"] = true
		fbSet := utils.GetFbSet(mapId)
		if fbSet == nil {
			msg = "地图未开放！"
			return
		}
		now := utils.NowUnix()
		fbRecord := fs.GetFbRecord(userId, mapId)
		if fbRecord == nil {
			fbRecord = &models.RecordFb{
				Uid:      userId,
				GpcId:    fbSet.Gpcs[0],
				LeftTime: fbSet.CoolTime,
				InMap:    mapId,
				SrcTime:  now,
			}
			if fs.GetDb().Create(fbRecord).RowsAffected < 1 {
				msg = "系统错误！"
				return
			}
		} else if fbRecord.SrcTime+fbSet.CoolTime < now {
			fbRecord.GpcId = fbSet.Gpcs[0]
			fbRecord.SrcTime = now
			fs.GetDb().Model(fbRecord).Update(gin.H{
				"gwid":    fbRecord.GpcId,
				"srctime": fbRecord.SrcTime,
			})
		}
		if fbRecord.GpcId < 1 {
			msg = "副本未开启"
			return
		}
		mmap = common.GetMMap(mapId)
		gpc = common.GetGpc(fbRecord.GpcId)
	}

	if mmap == nil {
		msg = "地图不存在！"
		return
	}
	if gpc == nil {
		msg = "怪物不存在！"
		return
	}

	rcache.SetFightTime(userId, now)
	user := fs.OptSvc.UserSrv.GetUserById(userId)
	mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
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
	if utils.IsFbMap((mapId)) {
		if mainPet.Level < mmap.Levels[0] {
			msg = "等级不足！"
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
	balls := []gin.H{}
	for _, p := range fs.OptSvc.PropSrv.GetCarryProps(userId, false) {
		p.GetM()
		if p.MModel.VaryName == 3 && strings.Contains(p.MModel.Effect, strconv.Itoa(gpc.ID)) {
			balls = append(balls, gin.H{
				"id":   p.ID,
				"name": p.MModel.Name,
				"sum":  p.Sums,
			})
		}
	}

	rcache.SetFightStatus(userId, gpc.ID, multiple, 0)
	fightInfo["interval_auto_attack"] = 10 - int(mainPet.ZbAttr.Special["time"])
	fightInfo["result"] = true
	fightInfo["user"] = gin.H{
		"nickname": user.Nickname,
		//"img":      fmt.Sprintf("1%s.gif", user.Headimg),
	}
	fightInfo["catch_ball"] = balls
	fightInfo["multiple"] = multiple
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
	fightInfo["gpc"] = gin.H{
		"id":          gpc.ID,
		"name":        gpc.Name,
		"level":       gpc.Level,
		"hp":          gpc.Hp,
		"stand_img":   gpc.ImgStand,
		"attack_img1": gpc.ImgAck,
		"attack_img2": gpc.ImgAck,
	}
	return
}

func (fs *FightService) GetGpcRandSkill(gpc *models.Gpc) *models.Uskill {
	if len(gpc.Skills) == 0 {
		return nil
	}
	skillItem := gpc.Skills[rand.Intn(len(gpc.Skills))]
	mSkill := common.GetMskill(skillItem.Sid)
	lv := skillItem.Level - 1
	uSkill := &models.Uskill{
		Level:       skillItem.Level,
		Uhp:         mSkill.UHpItem[lv],
		Ump:         mSkill.UMpItem[lv],
		MModel:      mSkill,
		AckValue:    mSkill.AckItem[lv],
		PlusValue:   mSkill.PlusItem[lv],
		EffectValue: mSkill.EffectItem[lv],
	}
	return uSkill

}

func (fs *FightService) Attack(userId, skillId int) (fightInfo gin.H, msg string) {
	msg = ""
	result := gin.H{
		"finish":      false,
		"finish_msg":  "",
		"exp":         0,
		"money":       0,
		"zb_money":    0,
		"pet_upgrade": false,
		"prop":        "",
		"has_next":    false,
		"in_tt":       false,
		"in_fb":       false,
		"in_group":    false,
	}
	fightInfo = gin.H{
		"success":  false,
		"result":   result,
		"mapid":    0,
		"multiple": 0,
		"pet":      nil,
		"gpc":      nil,
	}

	mapId := rcache.GetInMap(userId)
	fightInfo["mapid"] = mapId
	fightStatus := rcache.GetFightStatus(userId)
	if fightStatus == nil || fightStatus.GpcId == 0 {
		result["finish_msg"] = "战斗失效！1"
		return
	}
	fightInfo["multiple"] = fightStatus.Multiple

	now := utils.NowUnix()
	if rcache.GetAttackCoolTime(userId, now) > 0 {
		result["finish_msg"] = "操作过快！"
		return
	}
	if now-rcache.GetFightTime(userId) > 20 {
		result["finish_msg"] = "战斗失效！12"
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

	gpc := common.GetGpc(fightStatus.GpcId)

	gpcHp := gpc.Hp - fightStatus.DeHp
	if petHp < 1 || gpcHp < 1 {
		result["finish_msg"] = "战斗失效！2"
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
		result["finish_msg"] = "技能无效！"
		return
	}
	skill.GetM()
	if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
		result["finish_msg"] = "技能无效！"
		return
	}
	if skill.Ump > petMp {
		fightInfo["success"] = false
		return
	}

	var userInfo *models.UserInfo
	if utils.IsTTMap(mapId) {
		userInfo = fs.OptSvc.UserSrv.GetUserInfoById(userId)
		if userInfo.Tgt > 30 {
			gpc.Skills = append(gpc.Skills, &struct {
				Sid   int
				Level int
			}{Sid: 302, Level: 10})
		}
	}

	gpcSkill := fs.GetGpcRandSkill(gpc)
	if gpcSkill == nil {
		result["finish_msg"] = "战斗错误！"
		return
	}

	// 开始计算伤害
	success, apResult, bpResult := fs.GetFightResult(mainPet, skill, petStatus, gpc, gpcSkill, fightStatus.DeHp)
	if !success {
		msg = "宠物蓝量不足以使用此技能！"
		fightInfo["success"] = false
		return
	}
	if apResult == nil || bpResult == nil {
		result["finish_msg"] = "技能类型出错！"
		return
	}
	fightInfo["success"] = true
	fightInfo["pet"] = apResult
	fightInfo["gpc"] = bpResult
	mmap := common.GetMMap(mapId)
	if bpResult.Die {
		result["finish"] = true

		// 战斗胜利处理
		hasNext, petUpgrade, finishMsg, propNameList, isBagMax := fs.processFightVictory(user, userInfo, mainPet, mmap, gpc, int(mainPet.ZbAttr.Special["money"]))
		result["money"] = gpc.Money
		result["exp"] = gpc.Exp
		result["zb_money"] = int(mainPet.ZbAttr.Special["money"])
		result["finish_msg"] = finishMsg
		result["has_next"] = hasNext
		result["pet_upgrade"] = petUpgrade
		if petUpgrade {
			fs.GetZbAttr(mainPet, nil)
		}
		if len(propNameList) > 0 {
			result["prop"] = strings.Join(propNameList, ",")
		} else {
			result["prop"] = "无！"
		}
		if isBagMax {
			result["prop"] = result["prop"].(string) + "(背包已满！)"
		}

	} else if apResult.Die {
		result["finish"] = true
		result["finish_msg"] = "战斗失败！"

		// 战斗失败处理
		hasNext, finishMsg := fs.processFightDefeat(user, mmap)
		result["finish_msg"] = finishMsg
		result["has_next"] = hasNext
	}
	rcache.SetFightTime(userId, now)
	//fmt.Printf("本轮攻击结束：hp:%d,dehp:%d\n", mainPet.ZbAttr.Hp, mainPet.ZbAttr.Hp-apResult.Hp)
	rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-apResult.Hp, mainPet.ZbAttr.Mp-apResult.Mp)
	rcache.SetFightStatus(userId, gpc.ID, fightStatus.Multiple, gpc.Hp-bpResult.Hp)
	return fightInfo, msg

}

func (fs *FightService) missAttack(apHits, bpMiss int) bool {
	var num int
	if apHits > bpMiss*3 {
		num = 9
	} else if apHits > bpMiss*2 {
		num = 8
	} else if apHits > bpMiss {
		num = 8
	} else if apHits*2 > bpMiss {
		num = 7
	} else if apHits*3 > bpMiss {
		num = 7
	} else {
		num = 6
	}
	if num >= rand.Intn(10) {
		return false
	}
	return true
}

func (fs *FightService) hitsAddAc(apHits, bpMiss int) float64 {
	addNum := float64(apHits-bpMiss) / 100
	if addNum < 1 {
		addNum = 1
	} else if addNum > 1.5 {
		addNum = 1.5
	}
	return addNum
}

func (fs *FightService) isCritical(criticalRate float64) bool {
	return rand.Float64() < criticalRate
}

// 返回：是否缺蓝无法攻击，ap战斗结果，bp战斗结果
func (fs *FightService) GetFightResult(ap *models.UPet, apSkill *models.Uskill, apStatus *rcache.PetStatus, bp *models.Gpc, bpSkill *models.Uskill, gpcDeHp int) (bool, *APAttackResult, *BPAttackResult) {
	apResult := &APAttackResult{
		Effect: &struct {
			XX   int `json:"xx"`
			XM   int `json:"xm"`
			SHJS int `json:"shjs"`
			XHDX int `json:"xhdx"`
			XHFT int `json:"xhft"`
		}{},
	}
	bpResult := &BPAttackResult{}
	bpResult.BeMiss = fs.missAttack(bp.Hits, ap.ZbAttr.Miss)
	if apSkill == nil {
		return true, nil, nil
	}
	apSkill.GetM()
	apResult.SkillName = apSkill.MModel.Name
	if apSkill.MModel.Category == 1 {
		if apResult.Mp = ap.ZbAttr.Mp - apStatus.DeMp - apSkill.Ump; apResult.Mp < 0 {
			return false, nil, nil
		}
		if apResult.BeMiss = fs.missAttack(ap.ZbAttr.Hits, bp.Miss); !apResult.BeMiss {

			// 基础伤害
			apResult.Attack = int(float64(ap.ZbAttr.Ac+apSkill.AckValue)*(1+float64(apSkill.PlusValue)/100)) - bp.Mc
			//fmt.Printf("attack1 :%d\n", apResult.Attack)

			// 命中加成
			apResult.Attack = int(float64(apResult.Attack) * fs.hitsAddAc(ap.ZbAttr.Hits, bp.Miss))
			//fmt.Printf("attack2 :%d\n", apResult.Attack)

			// 伤害浮动
			apResult.Attack = int(float64(apResult.Attack) * (1 + float64(utils.RandInt(-10, 5))/100))
			//fmt.Printf("attack3 :%d\n", apResult.Attack)

			// 暴击伤害
			if apResult.Critical = fs.isCritical(ap.ZbAttr.Special["crit"]); apResult.Critical {
				apResult.Attack *= 2
			}
			if apResult.Attack < 1 {
				apResult.Attack = 1
			}

		}
	} else if apSkill.MModel.Category == 3 {
		if apResult.Mp = ap.ZbAttr.Mp - apStatus.DeMp - apSkill.Ump; apResult.Mp < 0 {
			return false, nil, nil
		}
		// 数据库记得是负值，所以要取反
		apResult.Reply = -apSkill.Uhp

	} else {
		return true, nil, nil
	}

	bpResult.SkillName = bpSkill.MModel.Name
	if bpSkill.MModel.Category == 1 {
		if bpResult.BeMiss = fs.missAttack(bp.Hits, ap.ZbAttr.Miss); !bpResult.BeMiss {
			// 基础伤害
			bpResult.Attack = int(float64(bp.Ac+bpSkill.AckValue)*(1+float64(bpSkill.PlusValue)/100)) - ap.ZbAttr.Mc

			// 命中加成
			bpResult.Attack = int(float64(bpResult.Attack) * fs.hitsAddAc(bp.Hits, ap.ZbAttr.Miss))

			// 伤害浮动
			bpResult.Attack = int(float64(bpResult.Attack) * (1 + float64(utils.RandInt(-10, 5))/100))
			if bpResult.Attack < 1 {
				bpResult.Attack = 1
			}
		}
	} else {
		return true, nil, nil
	}

	// 技能特效处理
	if apSkill.EffectValue != nil {
		switch apSkill.EffectValue.Key {
		case "hitshp":
			apResult.Effect.XX += int(float64(apResult.Attack) * apSkill.EffectValue.Value)
			break
		case "dxsh":
			apResult.Effect.XHDX += int(float64(bpResult.Attack) * apSkill.EffectValue.Value)
			break
		case "shjs":
			apResult.Effect.SHJS += int(float64(apResult.Attack) * apSkill.EffectValue.Value)
			break

		}
	}

	// 特效处理
	for k, v := range ap.ZbAttr.Special {
		switch k {
		case "dxsh":
			apResult.Effect.XHDX += int(float64(bpResult.Attack) * v)
			break
		case "hitshp":
			apResult.Effect.XX += int(float64(apResult.Attack) * v)
			break
		case "hitsmp":
			apResult.Effect.XM += int(float64(apResult.Attack) * v)
			break
		case "shjs":
			apResult.Effect.SHJS += int(float64(apResult.Attack) * v)
			break
		case "sdmp":
			//伤害的v转化为MP
			break
		case "szmp":
			//伤害的v以MP抵消
			break
		}
	}
	if bpResult.Hp = bp.Hp - gpcDeHp - (apResult.Attack + apResult.Effect.SHJS); bpResult.Hp <= 0 {
		bpResult.Hp = 0
		bpResult.Die = true
		bpResult.Attack = 0
		bpResult.Critical = false
		apResult.Effect.XHDX = 0
	}
	if apResult.Hp = ap.ZbAttr.Hp - apStatus.DeHp - bpResult.Attack + (apResult.Effect.XX + apResult.Effect.XHDX); !bpResult.Die && apResult.Hp <= 0 {
		apResult.Hp = 0
		apResult.Die = true
	}
	apResult.Mp += apResult.Effect.XM
	if apResult.Hp > ap.ZbAttr.Hp {
		apResult.Hp = ap.ZbAttr.Hp
	}
	if apResult.Mp > ap.ZbAttr.Mp {
		apResult.Mp = ap.ZbAttr.Mp
	}
	return true, apResult, bpResult
}

func (fs *FightService) processFightVictory(user *models.User, userInfo *models.UserInfo, pet *models.UPet, mmap *models.Map, gpc *models.Gpc, zbMoney int) (hasNext, petUpgrade bool, finishMsg string, propNameList []string, isBagMax bool) {

	fs.OptSvc.Begin()
	defer fs.OptSvc.Commit()
	var ratePids []*models.RatePid
	ratePids = append(ratePids, gpc.Drops...)

	IsOpenMap := true
	if utils.IsNeedKeyMap(mmap.ID) {
		IsOpenMap = slice.ContainsString(strings.Split(user.OpenMap, ","), strconv.Itoa(mmap.ID))
	}

	// 自动战斗奖励
	autoFlag := rcache.GetAutoFightFlag(user.ID)
	exp := gpc.Exp
	if autoFlag == rcache.AutoFightJb {
		exp = int(float64(exp) * 1.2)
		if user.AutoFightTimeM == 0 {
			hasNext = false
			finishMsg = "金币自动战斗次数已用完！"
		}
	} else if autoFlag == rcache.AutoFightYb {
		exp = int(float64(exp) * 1.5)
		if user.AutoFightTimeYb == 0 {
			hasNext = false
			finishMsg = "元宝自动战斗次数已用完！"
		}
	}
	money := gpc.Money + zbMoney

	if !IsOpenMap {
		exp = int(float64(exp) / NoKeyMapPunishmentRate)
		money = int(float64(exp) / NoKeyMapPunishmentRate)
	}

	petUpgrade = fs.OptSvc.PetSrv.IncreaseExp2Pet(pet, exp)
	fs.OptSvc.UserSrv.IncreaseJb(user.ID, money)

	if utils.IsFbMap(mmap.ID) {
		// 副本处理，检查是否有下一只怪，若没有则公告
		go FinishFightTask(user.ID, gpc.ID)
		fbRecord := fs.GetFbRecord(user.ID, mmap.ID)
		fbSet := utils.GetFbSet(mmap.ID)
		progress := 0
		for i, id := range fbSet.Gpcs {
			if fbRecord.GpcId == id {
				progress = i
			}
		}
		if progress == len(fbSet.Gpcs)-1 {
			// 副本结束
			fs.GetDb().Model(fbRecord).Update(gin.H{"gwid": 0, "srctime": utils.NowUnix()})
			AnnounceAll(user.Nickname, fmt.Sprintf("成功通过副本 %s！", mmap.Name))
			hasNext = false
			finishMsg = "恭喜您成功通过副本！"
		} else {
			fs.GetDb().Model(fbRecord).Update(gin.H{"gwid": fbSet.Gpcs[progress+1]})
			hasNext = true
		}

	} else if utils.IsTTMap(mmap.ID) {
		// 通天塔处理，记录下一只怪，检查是否需要升到下一阶塔，检查是否通报，检查是否有下一只怪，检查下一只怪是否需要水晶
		hasNext = true
		record := rcache.GetTTRecord(user.ID)
		gpcGroup := common.GetGpcGroup(record.GpcGroupId)
		if record.Index < len(gpcGroup.GpcList)-1 {

			record.Index += 1
			rcache.SetTTRecord(user.ID, record)
		} else {
			ratePids = append(ratePids, gpcGroup.DropList...)
			if userInfo.Tgt+1 >= 35 && userInfo.Tgt+1%5 == 0 {
				defer AnnounceAll(user.Nickname, fmt.Sprintf("成功通过通天塔第%d层，获得%s。", userInfo.Tgt+1, strings.Join(propNameList, ",")))

			}
			if userInfo.Tgt+1 >= 50 {
				hasNext = false
				fs.GetDb().Model(userInfo).Update(gin.H{"tgt": 0, "tgttime": utils.NowUnix()})
				rcache.DelTTRecord(user.ID)
			} else {
				fs.GetDb().Model(userInfo).Update(gin.H{"tgt": gorm.Expr("tgt+1")})
				fs.NewTTRecord(user.ID, userInfo.Tgt+1)
			}
		}
	} else if utils.IsHuPoWuMap(mmap.ID) {
		// 琥珀屋处理，检查是否有下一只怪
	} else if utils.IsYiWangMap(mmap.ID) {
		// 遗忘宫殿处理，检查是否有下一只怪，暂时不在这里处理
	} else {
		// 普通地图处理，检查是否打败BOSS，是否需要公告
		// 特殊时段地图掉落
		go FinishFightTask(user.ID, gpc.ID)
		hasNext = true
		if gpc.Boss == 3 {
			bossRecord := &models.BossRecord{}
			fs.GetDb().Where("gid=?", gpc.ID).First(bossRecord)
			if bossRecord.Id > 0 {
				fs.GetDb().Model(bossRecord).Update(gin.H{
					"glock":    0,
					"fightuid": user.ID,
					"dtime":    utils.NowUnix(),
				})
			} else {
				bossRecord.GpcId = gpc.ID
				bossRecord.Glock = 0
				bossRecord.FightUid = user.ID
				bossRecord.Dtime = utils.NowUnix()
				fs.GetDb().Create(bossRecord)
			}
			AnnounceAll(user.Nickname, fmt.Sprintf("消灭了BOSS[%s]，获得了大量宝物！", gpc.Name))
		}
	}
	//fmt.Printf("rand len:%d\n", len(ratePids))

	carryPropCnt := fs.OptSvc.PropSrv.GetCarryPropsCnt(user.ID)
	if carryPropCnt < user.BagPlace {
		propNameList = fs.OptSvc.PropSrv.AddPropRand(user, carryPropCnt, ratePids, IsOpenMap)
	}

	return

}

func (fs *FightService) processFightDefeat(user *models.User, mmap *models.Map) (hasNext bool, finishMsg string) {
	if utils.IsFbMap(mmap.ID) {
		hasNext = true
	} else if utils.IsTTMap(mmap.ID) {
		// 通天塔处理，记录下一只怪，检查是否需要升到下一阶塔，检查是否通报，检查是否有下一只怪，检查下一只怪是否需要水晶
		hasNext = false
		fs.OptSvc.UserSrv.UpdateUserInfo(user.ID)(gin.H{"tgt": 0, "tgttime": utils.NowUnix()})
		rcache.DelTTRecord(user.ID)
		rcache.SetTTFlag(user.ID, rcache.UserTTNone)
	} else if utils.IsHuPoWuMap(mmap.ID) {
		// 琥珀屋处理，检查是否有下一只怪
	} else if utils.IsYiWangMap(mmap.ID) {
		// 遗忘宫殿处理，检查是否有下一只怪，暂时不在这里处理
	} else {
		// 普通地图处理，检查是否打败BOSS，是否需要公告
		// 特殊时段地图掉落
		hasNext = true
	}
	return
}

func (fs *FightService) SetAutoStart(userId int, autoType string) (bool, string) {
	if autoType != "jb" && autoType != "yb" {
		return false, "参数出错！"
	}
	user := fs.OptSvc.UserSrv.GetUserById(userId)
	mapId := rcache.GetInMap(userId)
	if utils.IsTTMap(mapId) || utils.IsHuPoWuMap(mapId) || utils.IsSSZhanchangMap(mapId) {
		return false, "该地图无法进行自动战斗！"
	}
	if autoType == "jb" {
		if user.AutoFightTimeM > 0 {
			rcache.SetAutoFightFlag(userId, rcache.AutoFightJb)
			return true, ""
		}
	} else {
		if user.AutoFightTimeYb > 0 {
			rcache.SetAutoFightFlag(userId, rcache.AutoFightYb)
			return true, ""
		}
	}
	return false, "自动战斗次数不足！"

}

func (fs *FightService) TTUseSj(userId int) (bool, string) {
	userInfo := fs.OptSvc.UserSrv.GetUserInfoById(userId)
	now := utils.ToDayStartUnix()

	if now > userInfo.TgLastTime {
		// 今日没打过
	} else if now > userInfo.TgtTime {
		// 今日没结束过
		if userInfo.Tgt == 30 {
			if fs.OptSvc.UserSrv.DecreaseSj(userId, 200) {
				rcache.SetTTFlag(userId, rcache.UserTT31)
				return true, ""
			}
			return false, "水晶不足！"
		}
	} else {
		// 今日打过，而且结束了
		if fs.OptSvc.UserSrv.DecreaseSj(userId, 200) {
			rcache.SetTTFlag(userId, rcache.UserTTFh)
			return true, ""
		}
		return false, "水晶不足！"
	}
	return false, "无需使用水晶！"
}

func (fs *FightService) IncreaseTTRecordIndex(userId int, record *models.TTRecord) {
	if record.Index < 4 {
		record.Index += 1
		rcache.SetTTRecord(userId, record)
	}
}

func (fs *FightService) NewTTRecord(userId, level int) *models.TTRecord {
	gpcGroups := common.GetGpcGroupByLevel(level + 1)
	if len(gpcGroups) == 0 {
		panic(fmt.Sprintf("gpcGroup boss %d\n no find!", level))
	}
	gpcGroupId := gpcGroups[rand.Intn(len(gpcGroups))]
	record := &models.TTRecord{
		GpcGroupId: gpcGroupId,
		Index:      0,
	}
	rcache.SetTTRecord(userId, record)
	return record
}

// 返回结果：是否需要水晶、gpc
func (fs *FightService) GetTTGpc(userId int) (bool, *models.Gpc) {
	userInfo := fs.OptSvc.UserSrv.GetUserInfoById(userId)
	now := utils.ToDayStartUnix()
	var record *models.TTRecord

	if now > userInfo.TgLastTime {
		// 今日没打过
		fs.GetDb().Model(userInfo).Update(gin.H{"tglasttime": utils.NowUnix(), "tgt": 0})
		record = fs.NewTTRecord(userId, 0)
		rcache.SetTTFlag(userId, rcache.UserTTNone)
	} else if now > userInfo.TgtTime {
		// 今日没结束过
		if userInfo.Tgt == 30 && rcache.GetTTFlag(userId) != rcache.UserTT31 {
			return true, nil
		}
		record = rcache.GetTTRecord(userId)
		if record == nil {
			record = fs.NewTTRecord(userId, userInfo.Tgt)
		}
	} else {
		// 今日打过，而且结束了
		if rcache.GetTTFlag(userId) != rcache.UserTTFh {
			return true, nil
		}
		record = rcache.GetTTRecord(userId)
		if record == nil {
			record = fs.NewTTRecord(userId, userInfo.Tgt)
		}
	}
	//fmt.Printf("gpc list:%v\n", record.Gpcs)
	gpcGroup := common.GetGpcGroup(record.GpcGroupId)
	if record == nil || record.Index >= len(gpcGroup.GpcList) {
		// 没有记录，出错了
		fmt.Printf("not find!%d %d\n", record.Index, len(gpcGroup.GpcList))
		return false, nil
	}
	gpc := common.GetGpc(gpcGroup.GpcList[record.Index])
	return false, gpc
}

func (fs *FightService) GetTTUserRank() []gin.H {
	users := []*struct {
		NickName string `gorm:"column:nickname"`
		Level    int    `gorm:"column:level"`
	}{}
	fs.GetDb().Raw("select player.nickname as nickname, player_ext.tgt as level from player inner join player_ext on player.id=player_ext.uid where player_ext.tglasttime>? order by level desc limit 5", utils.ToDayStartUnix()).Scan(&users)
	userDatas := []gin.H{}
	for _, u := range users {
		userDatas = append(userDatas, gin.H{"nickname": u.NickName, "level": u.Level + 1})
	}
	return userDatas

}

func (fs *FightService) GetMapUsers(userId, mapId int) []gin.H {
	users := rcache.GetMapUserList(mapId)
	userIds := []int{}
	now := utils.NowUnix()
	for _, i := range users {
		if i == userId {
			continue
		}
		if now-rcache.GetInMapTime(i) < 60*10 || now-rcache.GetFightTime(i) < 60*10 {
			userIds = append(userIds, i)
		} else {
			rcache.DelMapUser(i, mapId)
		}
	}

	userDatas := []gin.H{}
	if len(userDatas) > 0 {
		userInfos := []*struct {
			Id       int    `gorm:"column:id"`
			NickName string `gorm:"column:nickname"`
		}{}
		fs.GetDb().Raw("select id, nickname from player where id in ?", userIds).Scan(&userInfos)
		for _, u := range userInfos {
			userDatas = append(userDatas, gin.H{"id": u.Id, "nickname": u.NickName})
		}
	}

	return userDatas
}

func (fs *FightService) Catch(userId, propId int) (bool, string) {
	fightStatus := rcache.GetFightStatus(userId)
	if fightStatus == nil {
		return false, "战斗已失效！"
	}
	prop := fs.OptSvc.PropSrv.GetProp(userId, propId, false)
	if prop == nil || prop.Sums == 0 {
		return false, "精灵球数量不足！"
	}
	prop.GetM()
	if prop.MModel.VaryName != 3 {
		return false, "选取道具并非精灵球！"
	}
	item := strings.Split(prop.MModel.Effect, ":")
	switch item[0] {
	case "get":
		if len(item) < 5 {
			return false, "选取道具并非可用精灵球！"
		}
		enableCatchList := strings.Split(item[1], "|")
		if !com.IsSliceContainsStr(enableCatchList, strconv.Itoa(fightStatus.GpcId)) {
			return false, "该精灵球不可捕捉该怪物！"
		}
		user := fs.OptSvc.UserSrv.GetUserById(userId)
		if fs.OptSvc.PropSrv.GetCarryPropsCnt(userId) >= user.BagPlace {
			return false, "背包空间不足！"
		}
		gpc := common.GetGpc(fightStatus.GpcId)
		bRate := com.StrTo(strings.ReplaceAll(item[2], "%", "")).MustFloat64()
		rate := float64(gpc.CatchRate)/100*float64(fightStatus.DeHp)/float64(gpc.Hp) + bRate
		fs.OptSvc.PropSrv.DecrProp(userId, propId, 1)
		rcache.DelFightStatus(userId)
		if utils.RandInt(1, int(1/rate)) == 1 {
			// 成功
			pid := com.StrTo(item[4]).MustInt()
			p := common.GetMProp(pid)
			fs.OptSvc.PropSrv.AddProp(userId, pid, 1, false)
			if item[3] == "2" {
				AnnounceAll(user.Nickname, fmt.Sprintf("成功的获取了: %s，太爽了！", p.Name))
			}
			return true, fmt.Sprintf("恭喜您获取道具 %s！", p.Name)
		} else {
			// 失败
			return true, "捕捉失败！怪物逃掉啦！"
		}
	case "catch":
		if len(item) < 4 {
			return false, "选取道具并非可用精灵球！"
		}
		gpc := common.GetGpc(fightStatus.GpcId)
		if gpc.CatchBid == 0 {
			return false, "该怪物不可捕捉！"
		}
		enableCatchList := strings.Split(item[1], "|")
		if !com.IsSliceContainsStr(enableCatchList, strconv.Itoa(fightStatus.GpcId)) {
			return false, "该精灵球不可捕捉该怪物！"
		}
		if fs.OptSvc.PetSrv.GetCarryPetCnt(userId) >= 3 {
			return false, "身上携带宠物过多，请至少空出一个位置！"
		}
		bRate := com.StrTo(strings.ReplaceAll(item[2], "%", "")).MustFloat64() / 100
		rate := float64(gpc.CatchRate)/100*float64(fightStatus.DeHp)/float64(gpc.Hp) + bRate
		fs.OptSvc.PropSrv.DecrProp(userId, propId, 1)
		rcache.DelFightStatus(userId)
		if utils.RandInt(1, int(1/rate)) == 1 {
			// 成功
			user := fs.OptSvc.UserSrv.GetUserById(userId)
			fs.OptSvc.PetSrv.CreatPetById(user, gpc.CatchBid)
			if item[3] == "2" {
				AnnounceAll(user.Nickname, fmt.Sprintf("成功的捕捉到了 %s ，太有才了！", gpc.Name))
			}
			return true, "恭喜您捕获一只新宠物！"
		} else {
			// 失败

			return true, "捕捉失败！"
		}
	default:
		return false, "选取道具并非可用精灵球！"
	}

}

func (fs *FightService) SetAutoSkill(userId, skillId int) (bool, string) {
	skill := fs.OptSvc.PetSrv.GetSkill(skillId)
	if skill == nil {
		return false, "技能不存在！"
	}
	skill.GetM()
	if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
		return false, "技能非战斗技能！"
	}
	rcache.SetAutoFightSkill(userId, skillId)
	return true, "设置技能成功！"
}

func (fs *FightService) StartFb(userId, mapId int, useSj bool) (bool, string) {
	mmap := common.GetMMap(mapId)
	if mmap == nil {
		return false, "地图不存在！"
	}
	if !utils.IsFbMap(mapId) {
		return false, "地图非副本地图！"
	}
	pid, sj := 0, 0
	for _, s := range strings.Split(mmap.Need, ",") {
		items := strings.Split(s, ":")
		if len(items) > 1 {
			if items[0] == "needitem" {
				pid = com.StrTo(items[1]).MustInt()
			} else if items[0] == "sj" {
				sj = com.StrTo(items[1]).MustInt()
			}
		}
	}
	fbSet := utils.GetFbSet(mapId)
	if fbSet == nil {
		return false, "地图未开放！"
	}
	now := utils.NowUnix()
	fbRecord := fs.GetFbRecord(userId, mapId)
	if fbRecord == nil {
		fbRecord = &models.RecordFb{
			Uid:      userId,
			GpcId:    fbSet.Gpcs[0],
			LeftTime: fbSet.CoolTime,
			InMap:    mapId,
			SrcTime:  now,
		}
		fs.GetDb().Create(fbRecord)
		return false, "该副本可直接进入！"
	} else if fbRecord.SrcTime+fbSet.CoolTime < now {
		fbRecord.GpcId = fbSet.Gpcs[0]
		fbRecord.SrcTime = now
		fs.GetDb().Model(fbRecord).Update(gin.H{
			"gwid":    fbRecord.GpcId,
			"srctime": fbRecord.SrcTime,
		})
		return false, "该副本可直接进入！"
	}
	if useSj {
		if sj > 0 {
			if fs.OptSvc.UserSrv.DecreaseSj(userId, sj) {
				fbRecord.GpcId = fbSet.Gpcs[0]
				fbRecord.SrcTime = now
				fs.GetDb().Model(fbRecord).Update(gin.H{
					"gwid":    fbRecord.GpcId,
					"srctime": fbRecord.SrcTime,
				})
				return true, "使用水晶刷新副本成功！"
			} else {
				return false, "水晶数量不足！"
			}
		} else {
			return false, "该副本不可用水晶刷新！"
		}
	} else {
		if pid > 0 {
			if fs.OptSvc.PropSrv.DecrPropByPid(userId, pid, 1) {
				fbRecord.GpcId = fbSet.Gpcs[0]
				fbRecord.SrcTime = now
				fs.GetDb().Model(fbRecord).Update(gin.H{
					"gwid":    fbRecord.GpcId,
					"srctime": fbRecord.SrcTime,
				})
				return true, "使用通行证刷新副本成功！"
			} else {
				return false, "您的通行证数量不足！"
			}
		} else {
			return false, "该副本不可用通行证刷新！"
		}
	}

}

type APAttackResult struct {
	Hp        int    `json:"hp"`
	Mp        int    `json:"mp"`
	Die       bool   `json:"die"`
	SkillName string `json:"skill_name"`
	BeMiss    bool   `json:"be_miss"`
	Critical  bool   `json:"critical"`
	Attack    int    `json:"attack"`
	Reply     int    `json:"reply"`
	Effect    *struct {
		XX   int `json:"xx"`
		XM   int `json:"xm"`
		SHJS int `json:"shjs"`
		XHDX int `json:"xhdx"`
		XHFT int `json:"xhft"`
	} `json:"effect"`
}

type BPAttackResult struct {
	Hp        int    `json:"hp"`
	Die       bool   `json:"die"`
	SkillName string `json:"skill_name"`
	BeMiss    bool   `json:"be_miss"`
	Critical  bool   `json:"critical"`
	Attack    int    `json:"attack"`
}
