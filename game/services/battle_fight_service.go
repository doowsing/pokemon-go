package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"math/rand"
	"pokemon/common/persistence"
	"pokemon/common/rcache"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"strconv"
	"strings"
	"time"
)

func (fs *FightService) GetSSBattleUserList() gin.H {
	ziranUsers := []*struct {
		Nickname string
		Num      int
	}{}
	fs.GetDb().Raw(`SELECT b.curjgvalue as num,p.nickname as nickname
								      FROM player as p,battlefield_user as b
									 WHERE p.id=b.uid and b.pos=1 and b.curjgvalue>0
									 ORDER BY b.curjgvalue desc
									 LIMIT 0,10`).Scan(&ziranUsers)
	anyeUsers := []*struct {
		Nickname string
		Num      int
	}{}
	fs.GetDb().Raw(`SELECT b.curjgvalue as num,p.nickname as nickname
								      FROM player as p,battlefield_user as b
									 WHERE p.id=b.uid and b.pos=2 and b.curjgvalue>0
									 ORDER BY b.curjgvalue desc
									 LIMIT 0,10`).Scan(&anyeUsers)
	ziranUserDatas := []gin.H{}
	for _, u := range ziranUsers {
		ziranUserDatas = append(ziranUserDatas, gin.H{"nickname": u.Nickname, "num": u.Num})
	}
	anyeUserDatas := []gin.H{}
	for _, u := range anyeUsers {
		anyeUserDatas = append(anyeUserDatas, gin.H{"nickname": u.Nickname, "num": u.Num})
	}

	now := utils.ToDayStartUnix()
	ziranCnt := 0
	fs.GetDb().Model(&models.SSBattleUser{}).Where("lastvtime>? and pos=1", now).Count(&ziranCnt)
	anyeCnt := 0
	fs.GetDb().Model(&models.SSBattleUser{}).Where("lastvtime>? and pos=2", now).Count(&anyeCnt)
	return gin.H{
		"anye":      anyeUserDatas,
		"anye_num":  anyeCnt,
		"ziran":     ziranUserDatas,
		"ziran_num": ziranCnt,
	}
}

func GetTimeConfigs(title string) []*models.TimeConfig {
	timeSets := []*models.TimeConfig{}
	persistence.GetOrm().Where("titles=?", title).Find(&timeSets)
	return timeSets
}

func (fs *FightService) SSBattleCheckTime() bool {
	now := time.Now()
	weekday := now.Weekday()
	h := now.Hour()
	mu := now.Minute()
	timeSets := GetTimeConfigs("battle")
	for _, s := range timeSets {
		if com.StrTo(s.Day).MustInt() == int(weekday) {
			startItems := strings.Split(s.StartTime, ":")
			endItems := strings.Split(s.EndTime, ":")
			if (com.StrTo(startItems[0]).MustInt() <= h && com.StrTo(startItems[1]).MustInt() <= mu) && com.StrTo(endItems[0]).MustInt() >= h && com.StrTo(endItems[1]).MustInt() >= mu {
				return true
			}
		}
	}
	return false
}

func (fs *FightService) GetSSBattleInfos() []*models.SSBattle {
	infos := []*models.SSBattle{}
	fs.GetDb().Find(&infos)
	return infos
}

func (fs *FightService) GetSSBattleInfo(id int) *models.SSBattle {
	info := &models.SSBattle{}
	fs.GetDb().Where("id=?", id).Find(&info)
	return info
}

func (fs *FightService) GetSSBattleUser(uid int) *models.SSBattleUser {
	userRecord := &models.SSBattleUser{}
	fs.GetDb().Where("uid=?", uid).First(userRecord)
	if userRecord.Id > 0 {
		return userRecord
	}
	return nil
}

func (fs *FightService) SSBattleEnter(userId int, isZiran bool) (string, gin.H) {
	if !fs.SSBattleCheckTime() {
		return "战场未开启！", gin.H{"result": false}
	}
	data := fs.GetSSBattleUserList()
	battleInfos := fs.GetSSBattleInfos()
	fs.SSBattleCheckStartNew(battleInfos, utils.NowUnix())
	index := 0
	if !isZiran {
		index = 1
	}
	thisNum, anotherNum := data["ziran_num"].(int), data["anye_num"].(int)
	if !isZiran {
		thisNum, anotherNum = anotherNum, thisNum
	}
	overNum := battleInfos[index].BfMlNum

	user := fs.OptSvc.UserSrv.GetUserById(userId)
	mainPet := fs.OptSvc.PetSrv.GetPetById(user.Mbid)

	if int(mainPet.CC) < battleInfos[index].BfLevelLimit {
		return "您的宠物成长不足进入战场！", gin.H{"result": false}
	}

	levels := ""
	levelDatas := []string{}
	victoryJgValue, victoryValue, failedJgValue, failedValue := 0, 0, 0, 0
	find := false
	for _, str := range strings.Split(battleInfos[index].LevelGet, ",") {
		items := strings.Split(strings.ReplaceAll(str, "|", ":"), ":")
		czlItems := strings.Split(items[0], "-")
		min := com.StrTo(czlItems[0]).MustInt()
		max := com.StrTo(czlItems[1]).MustFloat64()
		if !find && int(mainPet.CC) >= min && mainPet.CC <= max {
			find = true
			levels = items[0]
			victoryJgValue = com.StrTo(items[1]).MustInt()
			victoryValue = com.StrTo(items[2]).MustInt()
			failedJgValue = com.StrTo(items[3]).MustInt()
			failedValue = com.StrTo(items[4]).MustInt()
		}
		levelDatas = append(levelDatas, items[0])
	}
	if !find {
		return "您的宠物成长匹配战场空间！需在10~499成长", gin.H{"result": false}
	}

	userRecord := fs.GetSSBattleUser(userId)
	toDay := utils.ToDayStartUnix()
	if userRecord == nil || toDay > userRecord.LastVTime {
		// 未进入过
		// 检测进入资格，成长足够，人数条件达成
		if thisNum >= battleInfos[index].MaxUser {
			return "本阵营人数已满！", gin.H{"result": false}
		}
		if thisNum-overNum >= anotherNum {
			return fmt.Sprintf("我方当前人数超过对方至少 %d 人，已足够剿灭对方，请您稍后再来！", overNum), gin.H{"result": false}
		}
		if userRecord == nil {
			// 创建
			userRecord = &models.SSBattleUser{
				Uid:          userId,
				Pos:          index + 1,
				Bid:          user.Mbid,
				JgValue:      0,
				Levels:       levels,
				AddJgValue:   victoryJgValue,
				AckValue:     victoryValue,
				FailJgValue:  failedJgValue,
				FailAckValue: failedValue,
				LastVTime:    utils.NowUnix(),
				DoubleJg:     0,
				Tops:         0,
				CurJgValue:   0,
				BoxNum:       0,
				Nscf:         0,
				SubHp:        0,
				AddHp:        0,
			}
			fs.GetDb().Create(userRecord)
		} else {
			// 更新

			fs.GetDb().Model(userRecord).Update(gin.H{
				"levels":       levels,
				"addjgvalue":   victoryJgValue,
				"ackvalue":     victoryValue,
				"failjgvalue":  failedJgValue,
				"failackvalue": failedValue,
				"bid":          user.Mbid,
				"lastvtime":    utils.NowUnix(),
				"doublejg":     0,
				"pos":          index + 1,
				"tops":         0,
				"jgvalue":      gorm.Expr("curjgvalue+jgvalue"),
				"curjgvalue":   0,
				"boxnum":       0,
				"nscf":         0,
				"subhp":        0,
				"addhp":        0,
			})
		}
	} else {
		if (isZiran && userRecord.Pos != 1) || (!isZiran && userRecord.Pos != 2) {
			return "这里不欢迎间谍", gin.H{"result": false}
		}
		fs.GetDb().Model(userRecord).Update(gin.H{
			"levels":       levels,
			"addjgvalue":   victoryJgValue,
			"ackvalue":     victoryValue,
			"failjgvalue":  failedJgValue,
			"failackvalue": failedValue,
			"bid":          user.Mbid,
			"lastvtime":    utils.NowUnix(),
		})
	}
	data["result"] = true
	data["ziran_hp"], data["anye_hp"] = battleInfos[0].Hp, battleInfos[1].Hp
	data["ziran_max_hp"], data["anye_max_hp"] = battleInfos[0].SrcHp, battleInfos[1].SrcHp
	data["levels"] = levelDatas
	return "", data
}

func (fs *FightService) SSBattleStartFight(userId int, fightLevel string) (fightInfo gin.H, msg string) {
	msg = ""
	fightInfo = gin.H{
		"result":               false,
		"waittime":             0,
		"user":                 nil,
		"interval_auto_attack": 10,
		"pet":                  nil,
		"gpc":                  nil,
	}
	if !fs.SSBattleCheckTime() {
		msg = "战场未开启！"
		return
	}

	battleInfos := fs.GetSSBattleInfos()
	fs.SSBattleCheckStartNew(battleInfos, utils.NowUnix())

	rcache.SetInMap(userId, 0)
	now := utils.NowUnix()
	if t := rcache.GetFightCoolTime(userId, now); t > 0 {
		fightInfo["waittime"] = t
		msg = "战斗等待中！"
		return
	}

	var gpc *models.Gpc
	otherBattleUsers := []*models.SSBattleUser{}
	fs.GetDb().Where("levels=?", fightLevel).Find(&otherBattleUsers)
	if len(otherBattleUsers) == 0 {
		msg = "没有找到对手，请稍后重试！"
		return
	}
	gUserRecord := otherBattleUsers[rand.Intn(len(otherBattleUsers))]
	gPet := fs.OptSvc.PetSrv.GetPetById(gUserRecord.Bid)
	if gPet == nil {
		msg = "没有找到对手，请稍后重试！"
		return
	}
	gPet.GetM()
	gpc = &models.Gpc{
		ID:       gPet.ID,
		Name:     gPet.MModel.Name,
		Wx:       gPet.MModel.Wx,
		Level:    gPet.Level,
		Hp:       gPet.Hp,
		Mp:       gPet.Mp,
		Ac:       gPet.Ac,
		Mc:       gPet.Mc,
		Hits:     gPet.Hits,
		Speed:    gPet.Speed,
		Miss:     gPet.Miss,
		Skill:    "1:1",
		ImgStand: gPet.MModel.ImgStand,
		ImgAck:   gPet.MModel.ImgAck,
	}

	rcache.SetFightTime(userId, now)
	user := fs.OptSvc.UserSrv.GetUserById(userId)
	mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()

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

	rcache.SetFightStatus(userId, gpc.ID, 0, 0)
	fightInfo["interval_auto_attack"] = 10 - int(mainPet.ZbAttr.Special["time"])
	fightInfo["result"] = true
	fightInfo["user"] = gin.H{
		"nickname": user.Nickname,
	}
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

func (fs *FightService) SSBattleAttack(userId int) (fightInfo gin.H, msg string) {
	msg = ""
	result := gin.H{
		"finish":     false,
		"finish_msg": "",
		"jungong":    0,
	}
	fightInfo = gin.H{
		"success": false,
		"result":  result,
		"pet":     nil,
		"gpc":     nil,
	}

	fightStatus := rcache.GetFightStatus(userId)
	if fightStatus == nil || fightStatus.GpcId == 0 {
		result["finish_msg"] = "战斗失效！1"
		return
	}

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

	gPet := fs.OptSvc.PetSrv.GetPetById(fightStatus.GpcId)
	if gPet == nil {
		result["finish_msg"] = "战斗失效！2"
		return
	}
	gPet.GetM()
	gpc := &models.Gpc{
		ID:       gPet.ID,
		Name:     gPet.MModel.Name,
		Wx:       gPet.MModel.Wx,
		Level:    gPet.Level,
		Hp:       gPet.Hp,
		Mp:       gPet.Mp,
		Ac:       gPet.Ac,
		Mc:       gPet.Mc,
		Hits:     gPet.Hits,
		Speed:    gPet.Speed,
		Miss:     gPet.Miss,
		Skill:    "1:1",
		ImgStand: gPet.MModel.ImgStand,
		ImgAck:   gPet.MModel.ImgAck,
	}
	gpcHp := gpc.Hp - fightStatus.DeHp
	if petHp < 1 || gpcHp < 1 {
		result["finish_msg"] = "战斗失效！2"
		return
	}
	var skill *models.Uskill
	skill = fs.OptSvc.PetSrv.GetSkillBySid(mainPet.ID, 1)

	if skill == nil || skill.Bid != mainPet.ID {
		result["finish_msg"] = "技能无效！"
		return
	}
	skill.GetM()
	if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
		result["finish_msg"] = "技能无效！"
		return
	}

	gpc.Skills = []*struct {
		Sid   int
		Level int
	}{{Sid: 1, Level: 10}}
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
	battleUser := fs.GetSSBattleUser(userId)
	if bpResult.Die {
		result["finish"] = true
		getJg := 0
		getJg = int(float64(battleUser.AddJgValue) * (1 - (mainPet.CC-gPet.CC-20)/1000))
		if battleUser.DoubleJg == 1 {
			getJg *= 3
		} else {
			getJg *= 2
		}
		fs.GetDb().Model(battleUser).Update(gin.H{"curjgvalue": gorm.Expr("curjgvalue+?", getJg)})
		enemyId := 1
		if battleUser.Pos == 1 {
			enemyId = 2
		}
		enemyBattleInfo := fs.GetSSBattleInfo(enemyId)
		enemyHp := enemyBattleInfo.Hp - battleUser.AckValue
		if enemyHp < 0 {
			enemyHp = 0
		}

		fs.GetDb().Model(enemyBattleInfo).Update(gin.H{"hp": enemyHp})
		// 战斗胜利处理
		finishMsg := fmt.Sprintf("恭喜您，获得了本次战斗的胜利！您获得了 %d 军工", getJg)
		result["finish_msg"] = finishMsg
		if enemyHp == 0 {
			fs.SSBattleEndAward(false)
		}

	} else if apResult.Die {
		result["finish"] = true
		result["finish_msg"] = "战斗失败！"
		battleInfo := fs.GetSSBattleInfo(battleUser.Pos)
		ourHp := battleInfo.Hp - battleUser.AckValue
		if ourHp < 0 {
			ourHp = 0
		}
		fs.GetDb().Model(battleInfo).Update(gin.H{"hp": ourHp})

		// 战斗失败处理
		finishMsg := "很遗憾，战斗失败！"
		result["finish_msg"] = finishMsg
		if ourHp == 0 {
			fs.SSBattleEndAward(false)
		}
	}
	rcache.SetFightTime(userId, now)
	//fmt.Printf("本轮攻击结束：hp:%d,dehp:%d\n", mainPet.ZbAttr.Hp, mainPet.ZbAttr.Hp-apResult.Hp)
	rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-apResult.Hp, mainPet.ZbAttr.Mp-apResult.Mp)
	rcache.SetFightStatus(userId, gpc.ID, fightStatus.Multiple, gpc.Hp-bpResult.Hp)
	return fightInfo, msg

}

func (fs *FightService) SSBattleCheckStartNew(battleInfos []*models.SSBattle, now int) {
	isEnd := true
	for _, binfo := range battleInfos {
		if binfo.Ends > 0 {
			isEnd = true
		}
	}
	if isEnd {
		if t := rcache.GetSSBattleStartTime(); now-t > 60 {
			// 避免重复开始
			rcache.SetSSBattleStartTime(now)
			fs.GetDb().Model(&models.SSBattle{}).Update(gin.H{"startf": 1, "countf": 0, "success": 0, "ends": 0, "hp": gorm.Expr("srchp")})
			fs.GetDb().Model(&models.SSBattleUser{}).Where("curjgvalue>0").Update(gin.H{
				"tops":       0,
				"jgvalue":    gorm.Expr("jgvalue+curjgvalue"),
				"curjgvalue": 0,
				"nscf":       0,
				"subhp":      0,
				"addhp":      0,
			})
		}

	}
}

// 神圣战场结账
func (fs *FightService) SSBattleEndAward(timeOut bool) {
	battleInfos := fs.GetSSBattleInfos()
	if battleInfos[0].CountF == 0 && battleInfos[1].CountF == 0 {
		// 避免重复发奖
		now := utils.NowUnix()
		today := utils.ToDayStartUnix()
		if t := rcache.GetSSBattleEndTime(); !timeOut && now-t < 60 {
			return
		}
		rcache.SetSSBattleEndTime(now)

		victoryPos := 1
		victoryIndex := 0
		failurePos := 2
		failureIndex := 1
		if battleInfos[1].Hp > battleInfos[2].Hp {
			victoryPos = 2
			victoryIndex = 1
			failurePos = 1
			failureIndex = 0
		}
		// 更新胜者
		fs.GetDb().Model(battleInfos[victoryIndex]).Update(gin.H{"success": 1})

		fs.GetDb().Model(&models.SSBattle{}).Update(gin.H{"startf": 1, "countf": 1, "success": 0, "ends": 1, "hp": gorm.Expr("srchp")})
		AnnounceSys(fmt.Sprintf("本次战场结束，%s被打得溃不成军，%s取得了胜利！", battleInfos[failureIndex].PosName, battleInfos[victoryIndex].PosName))
		AnnounceSys(fmt.Sprintf("本次战场结束，%s被打得溃不成军，%s取得了胜利！", battleInfos[failureIndex].PosName, battleInfos[victoryIndex].PosName))
		AnnounceSys(fmt.Sprintf("本次战场结束，%s被打得溃不成军，%s取得了胜利！", battleInfos[failureIndex].PosName, battleInfos[victoryIndex].PosName))

		victoryUsers := []*models.SSBattleUser{}
		fs.GetDb().Where("pos=? and lastvtime>? and curjgvalue>0", victoryPos, today).Order("curjgvalue desc").Limit(10).Find(&victoryUsers)

		boxNum := 0
		addJg := 0
		for i, u := range victoryUsers {
			boxNum = 0
			addJg = 0
			switch i + 1 {
			case 1:
				boxNum = 10
				addJg = 2000
				break
			case 2, 3:
				boxNum = 6
				addJg = 1500
				break
			case 4, 5, 6:
				boxNum = 4
				addJg = 1000
				break
			case 7, 8, 9, 10:
				boxNum = 2
				addJg = 500
				break
			}
			fs.GetDb().Model(u).Update(gin.H{"tops": i + 1, "boxnum": boxNum, "curjgvalue": gorm.Expr("curjgvalue+?", addJg)})
		}
		failureUsers := []*models.SSBattleUser{}
		fs.GetDb().Where("pos=? and lastvtime>? and curjgvalue>0", failurePos, today).Order("curjgvalue desc").Limit(10).Find(&failureUsers)

		for i, u := range failureUsers {
			boxNum = 0
			addJg = 0
			switch i + 1 {
			case 1:
				boxNum = 5
				addJg = 1000
				break
			case 2, 3:
				boxNum = 3
				addJg = 500
				break
			case 4, 5, 6:
				boxNum = 2
				addJg = 300
				break
			case 7, 8, 9, 10:
				boxNum = 1
				addJg = 100
				break
			}
			fs.GetDb().Model(u).Update(gin.H{"tops": i + 1, "boxnum": boxNum, "curjgvalue": gorm.Expr("curjgvalue+?", addJg)})
		}

	}

}

func (fs *FightService) SSBattleUseProp(userId, t int) (bool, string) {
	if !fs.SSBattleCheckTime() {
		return false, "战场未开启！"
	}
	switch t {
	case 1:
		// 诅咒宝石
		pid := 203
		prop := fs.OptSvc.PropSrv.GetProp(userId, pid, false)
		if prop == nil || prop.Sums == 0 {
			return false, "道具数量不足！"
		}

		if !fs.OptSvc.PropSrv.DecrPropByPid(userId, pid, 1) {
			return false, "道具数量不足！"
		}
		battleUser := fs.GetSSBattleUser(userId)
		now := utils.NowUnix()

		if battleUser.LastVTime-now > 3600 {
			return false, "您没有参与本场战争，无法使用道具！"
		}
		if deT := battleUser.SubHp + 60 - now; deT > 0 {
			return false, fmt.Sprintf("道具使用时间冷却中，请过 %d 秒再试！", deT)
		}
		pos := 1
		if battleUser.Pos == 1 {
			pos = 2
		}
		battleInfo := fs.GetSSBattleInfo(pos)
		if battleInfo.Hp < 1000 {
			return false, fmt.Sprintf("对方女神生命低于 1000 点，无法使用该道具!")
		}
		fs.GetDb().Model(battleInfo).Update(gin.H{"hp": gorm.Expr("hp-100")})
		fs.GetDb().Model(battleUser).Update(gin.H{"curjgvalue": gorm.Expr("curjgvalue+50"), "subhp": now})
		user := fs.OptSvc.UserSrv.GetUserById(userId)
		AnnounceAll(user.Nickname, fmt.Sprintf(" ,使用 <诅咒宝石>,诅咒对方女神,%s 女神HP减少 100 点!", battleInfo.PosName))
		return true, "使用道具成功，军功增加 50 点"
	case 2:
		// 天地树的果实
		pid := 204
		prop := fs.OptSvc.PropSrv.GetProp(userId, pid, false)
		if prop == nil || prop.Sums == 0 {
			return false, "道具数量不足！"
		}

		if !fs.OptSvc.PropSrv.DecrPropByPid(userId, pid, 1) {
			return false, "道具数量不足！"
		}
		battleUser := fs.GetSSBattleUser(userId)
		now := utils.NowUnix()

		if battleUser.LastVTime-now > 3600 {
			return false, "您没有参与本场战争，无法使用道具！"
		}
		if deT := battleUser.AddHp + 600 - now; deT > 0 {
			return false, fmt.Sprintf("道具使用时间冷却中，请过 %d 秒再试！", deT)
		}

		battleInfo := fs.GetSSBattleInfo(battleUser.Pos)
		battleInfo.Hp += 1000
		if battleInfo.Hp > 10000 {
			battleInfo.Hp = 10000
		}
		fs.GetDb().Model(battleInfo).Update(gin.H{"hp": battleInfo.Hp})
		fs.GetDb().Model(battleUser).Update(gin.H{"curjgvalue": gorm.Expr("curjgvalue+500"), "addhp": now})
		user := fs.OptSvc.UserSrv.GetUserById(userId)
		AnnounceAll(user.Nickname, fmt.Sprintf(" ,使用 <天地树的果实>,%s 女神HP恢复 1000 点!", battleInfo.PosName))
		return true, "使用道具成功，军功增加 500 点"
	case 3:
		// 女神圣水

		pid := 205
		prop := fs.OptSvc.PropSrv.GetProp(userId, pid, false)
		if prop == nil || prop.Sums == 0 {
			return false, "道具数量不足！"
		}

		if !fs.OptSvc.PropSrv.DecrPropByPid(userId, pid, 1) {
			return false, "道具数量不足！"
		}
		battleUser := fs.GetSSBattleUser(userId)
		if battleUser.Nscf == 1 {
			return false, "每场活动时，只能使用道具得到一次女神赐福！"
		}
		fs.GetDb().Model(battleUser).Update(gin.H{"doublejg": 1, "nscf": 1})
		break
	}
	return false, "参数出错！"
}

func (fs *FightService) SSBattleStoreData() []gin.H {
	data := []gin.H{}
	goods := []*models.SSBattleGood{}
	fs.GetDb().Find(goods)
	var prop *models.MProp
	for _, g := range goods {
		prop = common.GetMProp(g.Pid)
		data = append(data, gin.H{"id": g.Pid, "name": prop.Name, "price": g.Price})
	}
	return data
}

func (fs *FightService) SSBattleGetAward(userId, t int) (bool, string) {
	pid := 0
	switch t {
	case 1:
		pid = 1059
		break
	case 2:
		pid = 1060
		break
	case 3:
		pid = 1061
		break
	}
	if pid == 0 {
		return false, "参数出错！"
	}
	battleUser := fs.GetSSBattleUser(userId)

	if battleUser == nil || battleUser.BoxNum == 0 {
		return false, "您没进入排名或已经领取奖励！"
	}
	if fs.OptSvc.PropSrv.AddProp(userId, pid, battleUser.BoxNum, true) {
		prop := common.GetMProp(pid)
		fs.GetDb().Create(&models.SSBattleLog{
			Uid:   userId,
			UseJg: 0,
			Type:  "GoldBox",
			Num:   strconv.Itoa(battleUser.BoxNum),
			Pid:   strconv.Itoa(pid),
			Times: utils.NowUnix(),
		})
		return true, fmt.Sprintf("恭喜您，获得 %s * %d !", prop.Name, battleUser.BoxNum)
	} else {
		return false, "您的背包空间不足！"
	}
}

func (fs *FightService) SSBattleConvertExp(userId, num int) (bool, string) {
	if num == 0 {
		return false, "兑换数量必须大于0！"
	}
	battleUser := fs.GetSSBattleUser(userId)

	if battleUser == nil || battleUser.JgValue < num {
		return false, "您的军工数量不足！"
	}
	if fs.GetDb().Model(battleUser).Where("jgvalue>=?", num).Update(gin.H{"jgvalue": gorm.Expr("jgvalue-?", num)}).RowsAffected == 0 {

		return false, "您的军工数量不足！"
	}
	exp := num * 100
	fs.OptSvc.PetSrv.IncreaseExp2MainPet(userId, exp)

	fs.GetDb().Create(&models.SSBattleLog{
		Uid:   userId,
		UseJg: num,
		Type:  "BattleExp",
		Num:   strconv.Itoa(exp),
		Pid:   "",
		Times: utils.NowUnix(),
	})
	return true, fmt.Sprintf("恭喜您，主战宠物获得了 %d 点经验", exp)

}

func (fs *FightService) SSBattleConvertProp(userId, pid, num int) (bool, string) {
	if num == 0 {
		return false, "兑换数量必须大于0！"
	}
	good := &models.SSBattleGood{}
	fs.GetDb().Where("pid=?", pid).First(good)
	if good.Id == 0 {
		return false, "商品不存在！"
	}
	needJg := good.Price * num

	battleUser := fs.GetSSBattleUser(userId)

	if battleUser == nil || battleUser.JgValue < needJg {
		return false, "您的军工数量不足！"
	}

	if fs.GetDb().Model(battleUser).Where("jgvalue>=?", needJg).Update(gin.H{"jgvalue": gorm.Expr("jgvalue-?", needJg)}).RowsAffected == 0 {

		return false, "您的军工数量不足！"
	}
	if fs.OptSvc.PropSrv.AddProp(userId, pid, num, true) {
		prop := common.GetMProp(pid)
		fs.GetDb().Create(&models.SSBattleLog{
			Uid:   userId,
			UseJg: needJg,
			Type:  "BattleProps",
			Num:   strconv.Itoa(num),
			Pid:   strconv.Itoa(pid),
			Times: utils.NowUnix(),
		})
		return true, fmt.Sprintf("获得 %s * %d !", prop.Name, battleUser.BoxNum)
	} else {
		return false, "您的背包空间不足！"
	}

}

func (fs *FightService) FamilyBattleInfo(userId int) gin.H {
	family := fs.OptSvc.NpcSrv.GetMyFamily(userId)
	memberData := []gin.H{}
	if family != nil {
		members := []*struct {
			Nickname string
			Honor    int
		}{}
		fs.GetDb().Raw("SELECT m.honor as honor,u.nickname as u FROM guild_members m inner join player u on m.member_id=u.id WHERE m.guild_id = ? ORDER BY m.honor DESC", family.Id).
			Scan(&members)
		for _, m := range members {
			memberData = append(memberData, gin.H{"nickname": m.Nickname, "honor": m.Honor})
		}
	}

	families := []*models.Family{}
	fs.GetDb().Order("honor DESC").Find(&families)
	familyData := []gin.H{}
	for _, f := range families {
		enable_invite_battle := false
		if family != nil && f.Id != family.Id {
			deL := f.Level - family.Level
			if deL <= 5 || deL >= -5 {
				enable_invite_battle = true
			}
		}
		familyData = append(familyData, gin.H{"name": f.Name, "honor": f.Honor, "enable_invite_battle": enable_invite_battle})
	}

	battleBooks := []gin.H{}
	if family != nil {
		today := utils.ToDayStartUnix()
		desenders := []*struct {
			ChallengerId int `gorm:"column:challenger_id"`
			Name         string
			Flag         int
		}{}
		fs.GetDb().Raw("SELECT gc.challenger_id AS challenger_id,g.name AS name,gc.flags AS flag FROM guild_challenges gc inner join guild g on gc.challenger_id=g.id WHERE gc.defenser_id = ? AND gc.create_time>?", family.Id, today).
			Scan(&desenders)
		challengers := []*struct {
			ChallengerId int `gorm:"column:challenger_id"`
			Name         string
			Flag         int
		}{}
		fs.GetDb().Raw("SELECT gc.challenger_id AS challenger_id,g.name AS name,gc.flags AS flag FROM guild_challenges gc inner join guild g on gc.defenser_id=g.id WHERE gc.challenger_id = ? AND gc.create_time>?", family.Id, today).
			Scan(&challengers)
		for _, gc := range desenders {
			status := "未接受"
			enable_accept := false
			if gc.Flag == 1 {
				status = "已接受"
				enable_accept = true
			}
			battleBooks = append(battleBooks, gin.H{"id": gc.ChallengerId, "name": gc.Name, "status": status, "enable_accept": enable_accept})
		}
		for _, gc := range challengers {
			status := "未接受"
			enable_accept := false
			if gc.Flag == 1 {
				status = "已接受"
			}
			battleBooks = append(battleBooks, gin.H{"id": gc.ChallengerId, "name": gc.Name, "status": status, "enable_accept": enable_accept})
		}
	}
	introduce := "暂未更新"
	introduceSet := common.GetWelcome("guild_battle")
	if introduceSet != nil {
		introduce = introduceSet.Content
	}
	return gin.H{
		"introduce":   introduce,
		"families":    familyData,
		"members":     memberData,
		"battle_book": battleBooks,
	}
}

func (fs *FightService) FamilyBattleCheckTime() bool {

	now := time.Now()
	weekday := now.Weekday()
	h := now.Hour()
	mu := now.Minute()
	timeSets := GetTimeConfigs("guild_battle")
	for _, s := range timeSets {
		if com.StrTo(s.Day).MustInt() == int(weekday) {
			startItems := strings.Split(s.StartTime, ":")
			endItems := strings.Split(s.EndTime, ":")
			if (com.StrTo(startItems[0]).MustInt() <= h && com.StrTo(startItems[1]).MustInt() <= mu) && com.StrTo(endItems[0]).MustInt() >= h && com.StrTo(endItems[1]).MustInt() >= mu {
				return true
			}
		}
	}
	return false
}

func (fs *FightService) FamilyBattleInvite(userId, defenderId int) (bool, string) {
	if fs.FamilyBattleCheckTime() {
		return false, "战场已开始！无法再下战书了"
	}
	member := fs.OptSvc.NpcSrv.GetFamilyMember(userId)
	if member == nil {
		return false, "您未加入家族！"
	}
	if member.Authority != 2 {
		return false, "您无权下战书！"
	}
	if member.FamilyId == defenderId {
		return false, "您不能给自己的家族下战书！"
	}
	today := utils.ToDayStartUnix()
	battleRecord := &models.FamilyChallenge{}
	fs.GetDb().Where("(challenger_id = ? or defenser_id = ? or challenger_id = ? or defenser_id = ?) AND flags = 1 and create_time>=?", defenderId, defenderId, member.FamilyId, member.FamilyId, today).First(battleRecord)
	if battleRecord.Id > 0 {
		return false, "您的家族或者对方家族已经接受战书，不能再下战书了！"
	}
	cnt := 0
	battleRecord.Id = 0
	fs.GetDb().Model(battleRecord).Where("challenger_id=? and create_time>=?", member.FamilyId, today).Count(&cnt)
	if cnt >= 3 {
		return false, "您的家族当前已经发出3份战书，不能再发了！"
	}
	cnt = 0
	battleRecord.Id = 0
	fs.GetDb().Model(battleRecord).Where("defenser_id=? and create_time>=?", defenderId, today).Count(&cnt)
	if cnt >= 3 {
		return false, "该家族已经收到三份战书，不能再下了，试试别的吧"
	}
	cnt = 0
	battleRecord.Id = 0
	fs.GetDb().Model(battleRecord).Where("challenger_id=? and defenser_id=? and create_time>=?", member.FamilyId, defenderId, today).Count(&cnt)
	if cnt >= 1 {
		return false, "您的家族已经对此家族下了战书"
	}
	myFamily := fs.OptSvc.NpcSrv.GetFamily(member.FamilyId)
	defenderFamily := fs.OptSvc.NpcSrv.GetFamily(defenderId)
	if defenderFamily == nil {
		return false, "所选家族不存在！"
	}
	if deL := myFamily.Level - defenderFamily.Level; !(deL <= 5 || deL >= -5) {
		return false, "您只能对等级相差为5的家族下战书！"
	}
	battleRecord = &models.FamilyChallenge{
		ChallengerId: member.FamilyId,
		DefenserId:   defenderId,
		CreateTime:   utils.NowUnix(),
	}
	return true, "下战书成功！"
}

func (fs *FightService) FamilyBattleAccept(userId, challengerId int) (bool, string) {
	if fs.FamilyBattleCheckTime() {
		return false, "战场已开始！无法再下战书了"
	}
	member := fs.OptSvc.NpcSrv.GetFamilyMember(userId)
	if member == nil {
		return false, "您未加入家族！"
	}
	if member.Authority != 2 {
		return false, "您无权接受战书！"
	}
	today := utils.ToDayStartUnix()
	battleRecord := &models.FamilyChallenge{}
	fs.GetDb().Where("(challenger_id = ? or defenser_id = ? or challenger_id = ? or defenser_id = ?) AND flags = 1 and create_time>=?", challengerId, challengerId, member.FamilyId, member.FamilyId, today).First(battleRecord)
	if battleRecord.Id > 0 {
		return false, "您的家族或者对方家族已经接受战书，不能再接受了！"
	}
	battleRecord.Id = 0
	fs.GetDb().Where("challenger_id = ? and defenser_id = ? and create_time>=?", challengerId, member.FamilyId, today).First(battleRecord)
	if battleRecord.Id == 0 {
		return false, "对方并未给您家族下战书！"
	}
	myFamily := fs.OptSvc.NpcSrv.GetFamily(member.FamilyId)
	defenderFamily := fs.OptSvc.NpcSrv.GetFamily(challengerId)
	if defenderFamily == nil {
		return false, "所选家族不存在！"
	}
	if deL := myFamily.Level - defenderFamily.Level; !(deL <= 5 || deL >= -5) {
		return false, "您接受的战书请求等级相差最多为5！"
	}
	fs.GetDb().Model(&models.FamilyChallenge{}).Where("challenger_id = ? and defenser_id = ? and create_time>=?", challengerId, member.FamilyId, today).
		Update(gin.H{"flags": 1})
	fs.GetDb().Where("(challenger_id = ? or defenser_id = ? or challenger_id=? or defenser_id = ?) and create_time>=? and flags=0", challengerId, challengerId, member.FamilyId, member.FamilyId, today).Delete(&models.FamilyChallenge{})
	return true, "接受战书成功！"
}

func (fs *FightService) FamilyBattleStartFight(userId int) (fightInfo gin.H, msg string) {
	msg = ""
	fightInfo = gin.H{
		"result":               false,
		"waittime":             0,
		"user":                 nil,
		"interval_auto_attack": 10,
		"pet":                  nil,
		"gpc":                  nil,
	}
	if !fs.FamilyBattleCheckTime() {
		msg = "战场未开启！"
		return
	}
	now := time.Now()
	y, m, d := now.Date()
	today := int(time.Date(y, m, d, 0, 0, 0, 0, now.Location()).Unix())
	member := fs.OptSvc.NpcSrv.GetFamilyMember(userId)
	if member == nil {
		msg = "您未加入家族！"
		return
	}
	rcache.SetInMap(userId, 0)
	if t := rcache.GetFightCoolTime(userId, int(now.Unix())); t > 0 {
		fightInfo["waittime"] = t
		msg = "战斗等待中！"
		return
	}
	familyBattleBook := &models.FamilyChallenge{}
	fs.GetDb().Where("challenger_id=? and create_time>? and flags=1", member.FamilyId, today).First(familyBattleBook)
	if familyBattleBook.Id == 0 {
		msg = "您的家族并未下战书或者没有家族接受战书！"
		return
	}

	enemyMember := &models.FamilyMember{}
	cnt := 0
	fs.GetDb().Model(enemyMember).Where("guild_id=?", familyBattleBook.DefenserId).Count(&cnt)
	if cnt == 0 {
		msg = "对方家族已解散或者不存在！"
		return
	}
	index := rand.Intn(cnt)
	enemyMember.UserId = 0
	fs.GetDb().Model(enemyMember).Where("guild_id=?", familyBattleBook.DefenserId).Offset(index).First(enemyMember)
	enemy := fs.OptSvc.UserSrv.GetUserById(enemyMember.UserId)

	var gpc *models.Gpc
	gPet := fs.OptSvc.PetSrv.GetPetById(enemy.Mbid)
	if gPet == nil {
		msg = "没有找到对手，请稍后重试！"
		return
	}
	gPet.GetM()
	gpc = &models.Gpc{
		ID:       gPet.ID,
		Name:     gPet.MModel.Name,
		Wx:       gPet.MModel.Wx,
		Level:    gPet.Level,
		Hp:       gPet.Hp,
		Mp:       gPet.Mp,
		Ac:       gPet.Ac,
		Mc:       gPet.Mc,
		Hits:     gPet.Hits,
		Speed:    gPet.Speed,
		Miss:     gPet.Miss,
		Skill:    "1:1",
		ImgStand: gPet.MModel.ImgStand,
		ImgAck:   gPet.MModel.ImgAck,
	}

	rcache.SetFightTime(userId, int(now.Unix()))
	user := fs.OptSvc.UserSrv.GetUserById(userId)
	mainPet := fs.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()

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

	rcache.SetFightStatus(userId, gpc.ID, enemyMember.UserId, 0)
	fightInfo["interval_auto_attack"] = 10 - int(mainPet.ZbAttr.Special["time"])
	fightInfo["result"] = true
	fightInfo["user"] = gin.H{
		"nickname": user.Nickname,
	}
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

func (fs *FightService) FamilyBattleAttack(userId int) (fightInfo gin.H, msg string) {
	msg = ""
	result := gin.H{
		"finish":     false,
		"finish_msg": "",
		"jungong":    0,
	}
	fightInfo = gin.H{
		"success": false,
		"result":  result,
		"pet":     nil,
		"gpc":     nil,
	}

	fightStatus := rcache.GetFightStatus(userId)
	if fightStatus == nil || fightStatus.GpcId == 0 {
		result["finish_msg"] = "战斗失效！1"
		return
	}

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

	gPet := fs.OptSvc.PetSrv.GetPetById(fightStatus.GpcId)
	if gPet == nil {
		result["finish_msg"] = "战斗失效！2"
		return
	}
	gPet.GetM()
	gpc := &models.Gpc{
		ID:       gPet.ID,
		Name:     gPet.MModel.Name,
		Wx:       gPet.MModel.Wx,
		Level:    gPet.Level,
		Hp:       gPet.Hp,
		Mp:       gPet.Mp,
		Ac:       gPet.Ac,
		Mc:       gPet.Mc,
		Hits:     gPet.Hits,
		Speed:    gPet.Speed,
		Miss:     gPet.Miss,
		Skill:    "1:1",
		ImgStand: gPet.MModel.ImgStand,
		ImgAck:   gPet.MModel.ImgAck,
	}
	gpcHp := gpc.Hp - fightStatus.DeHp
	if petHp < 1 || gpcHp < 1 {
		result["finish_msg"] = "战斗失效！2"
		return
	}
	var skill *models.Uskill
	skill = fs.OptSvc.PetSrv.GetSkillBySid(mainPet.ID, 1)

	if skill == nil || skill.Bid != mainPet.ID {
		result["finish_msg"] = "技能无效！"
		return
	}
	skill.GetM()
	if skill.MModel.Category != 1 && skill.MModel.Category != 3 {
		result["finish_msg"] = "技能无效！"
		return
	}

	gpc.Skills = []*struct {
		Sid   int
		Level int
	}{{Sid: 1, Level: 10}}
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
	if bpResult.Die {
		result["finish"] = true
		member := fs.OptSvc.NpcSrv.GetFamilyMember(userId)
		myFamily := fs.OptSvc.NpcSrv.GetFamily(member.FamilyId)
		enemyMember := fs.OptSvc.NpcSrv.GetFamilyMember(fightStatus.Multiple)
		enemyFamily := fs.OptSvc.NpcSrv.GetFamily(enemyMember.FamilyId)

		getHonor := int(10 * (1 + float64(myFamily.Level-enemyFamily.Level)/10))
		fs.GetDb().Model(myFamily).Update(gin.H{"honor": gorm.Expr("honor+?", getHonor)})
		fs.GetDb().Model(member).Update(gin.H{"honor": gorm.Expr("honor+?", getHonor)})
		fs.GetDb().Model(&models.FamilyChallenge{}).Where("challenger_id=? and defenser_id=? and create_time>? and flags=1", myFamily.Id, enemyFamily.Id, utils.ToDayStartUnix()).
			Update(gin.H{"challenger_score": gorm.Expr("challenger_score+1")})

		// 战斗胜利处理
		finishMsg := fmt.Sprintf("胜利方获得荣誉：%d", getHonor)
		result["finish_msg"] = finishMsg

	} else if apResult.Die {
		result["finish"] = true

		member := fs.OptSvc.NpcSrv.GetFamilyMember(userId)
		myFamily := fs.OptSvc.NpcSrv.GetFamily(member.FamilyId)
		enemyMember := fs.OptSvc.NpcSrv.GetFamilyMember(fightStatus.Multiple)
		enemyFamily := fs.OptSvc.NpcSrv.GetFamily(enemyMember.FamilyId)

		getHonor := int(10 * (1 + float64(enemyFamily.Level-myFamily.Level)/10))
		fs.GetDb().Model(enemyFamily).Update(gin.H{"honor": gorm.Expr("honor+?", getHonor)})
		fs.GetDb().Model(enemyMember).Update(gin.H{"honor": gorm.Expr("honor+?", getHonor)})
		fs.GetDb().Model(&models.FamilyChallenge{}).Where("challenger_id=? and defenser_id=? and create_time>? and flags=1", myFamily.Id, enemyFamily.Id, utils.ToDayStartUnix()).
			Update(gin.H{"defenser_score": gorm.Expr("defenser_score+1")})
		result["finish_msg"] = "战斗失败！"

		// 战斗失败处理
		finishMsg := "很遗憾，战斗失败！"
		result["finish_msg"] = finishMsg
	}
	rcache.SetFightTime(userId, now)
	//fmt.Printf("本轮攻击结束：hp:%d,dehp:%d\n", mainPet.ZbAttr.Hp, mainPet.ZbAttr.Hp-apResult.Hp)
	rcache.SetPetStatus(mainPet.ID, mainPet.ZbAttr.Hp-apResult.Hp, mainPet.ZbAttr.Mp-apResult.Mp)
	rcache.SetFightStatus(userId, gpc.ID, fightStatus.Multiple, gpc.Hp-bpResult.Hp)
	return fightInfo, msg

}
