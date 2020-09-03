package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"math/rand"
	"pokemon/common/rcache"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"strconv"
	"strings"
	"time"
)

type NpcServices struct {
	BaseService
}

func NewNpcServices(osrc *OptService) *NpcServices {
	us := &NpcServices{}
	us.SetOptSrc(osrc)
	return us
}

// 摩卡志
func (us *NpcServices) GetCardRecord(userId int) []*models.CardRecord {
	var rs []*models.CardRecord
	us.GetDb().Where("uid=?", userId).Find(&rs)
	return rs
}

func (us *NpcServices) GetCardSeriesDatas() []gin.H {
	var data []gin.H
	for _, ss := range common.GetAllCardSeries() {
		data = append(data, gin.H{
			"id":   ss.ID,
			"name": ss.Name,
		})
	}
	return data
}

func (us *NpcServices) GetCardSeriesData(userId, sid int) (datas []gin.H) {
	series := common.GetCardSeries(sid)
	if series == nil {
		return
	}
	var records []*struct {
		Id   int    `gorm:"column:id"`
		Name string `gorm:"column:name"`
		Sum  int    `gorm:"column:sum"`
	}
	us.GetDb().Raw("select props.id as id, props.name as name, record.sum as sum from t_card_user as record left join props on record.card_pid=props.id where record.card_pid in (?) and record.uid=?", common.GetMPropIdsByName(series.CardList), userId).Scan(&records)

	var tmpProp *models.MProp
	for _, name := range series.CardList {
		data := gin.H{
			"name": name,
			"img":  "",
		}
		num := 0
		for _, p := range records {
			if p.Name == name {
				num = p.Sum
				tmpProp = common.GetMProp(p.Id)
				data["img"] = tmpProp.Img
				break
			}
		}
		data["num"] = num
		datas = append(datas, data)
	}
	return datas
}

func (us *NpcServices) GetCardPrizeDatas(userId int) []gin.H {
	var datas []gin.H
	var records []*models.CardRecord
	us.GetDb().Where("uid=?", userId).Find(&records)
	pid2record := make(map[int]int)
	for _, r := range records {
		pid2record[r.CardPid] = r.Sum
	}
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	prizeRecords := strings.Split(userInfo.FHasGetPrize, ",")
	for _, p := range common.GetAllCardPrize() {
		status := 0
		if com.IsSliceContainsStr(prizeRecords, strconv.Itoa(p.ID)) {
			status = 1
		}
		needs := []string{}
		for _, n := range p.NeedList {
			if status == 0 {
				if pid := common.GetMPropIdByName(n.Name); pid > 0 && pid2record[pid] < n.Num {
					status = -1
				}
			}

			needs = append(needs, fmt.Sprintf("需要 %s %d 个", n.Name, n.Num))
		}
		awards := []string{}
		for _, n := range p.PrizeList {
			p := common.GetMProp(n.Pid)
			awards = append(awards, fmt.Sprintf("%s %d 个", p.Name, n.Num))
		}
		datas = append(datas, gin.H{
			"id":     p.ID,
			"title":  p.Title,
			"need":   needs,
			"awards": awards,
			"status": status,
		})
	}
	return datas
}

func (us *NpcServices) GetCardPrize(userId, prizeId int) (bool, string) {
	prizeItem := common.GetCardPrize(prizeId)
	if prizeItem == nil {
		return false, "没有这个奖励选项！"
	}
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	hasPrizes := strings.Split(userInfo.FHasGetPrize, ",")
	if com.IsSliceContainsStr(strings.Split(userInfo.FHasGetPrize, ","), strconv.Itoa(prizeId)) {
		return false, "该奖励已领过！"
	}
	var records []*models.CardRecord
	us.GetDb().Where("uid=?", userId).Find(&records)
	pid2record := make(map[int]int)
	for _, r := range records {
		pid2record[r.CardPid] = r.Sum
	}
	for _, n := range prizeItem.NeedList {
		if pid := common.GetMPropIdByName(n.Name); pid > 0 && pid2record[pid] < n.Num {
			return false, "领取条件未满足！"
		}

	}
	user := us.OptSvc.UserSrv.GetUserById(userId)
	if needPlace := len(prizeItem.PrizeList) - (user.BagPlace - us.OptSvc.PropSrv.GetCarryPropsCnt(userId)); needPlace > 0 {
		return false, fmt.Sprintf("背包空间不足！请至少准备%d的空间", needPlace)
	}
	for _, n := range prizeItem.PrizeList {
		us.OptSvc.PropSrv.AddProp(userId, n.Pid, n.Num, false)
	}
	hasPrizes = append(hasPrizes, strconv.Itoa(prizeId))
	us.GetDb().Exec("update player_ext set F_has_get_prize=? where uid=?", strings.Join(hasPrizes, ","), userId)
	return true, "领取奖励成功！"
}

func (us *NpcServices) UpdateCardUserTitles(user *models.User) {
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(user.ID)
	hasTitles := strings.Split(userInfo.FHasTitle, ",")
	var records []*models.CardRecord
	us.GetDb().Where("uid=?", user.ID).Find(&records)
	pid2record := make(map[int]int)
	for _, r := range records {
		pid2record[r.CardPid] = r.Sum
	}
	for _, title := range common.GetAllCardTitle() {
		if com.IsSliceContainsStr(hasTitles, strconv.Itoa(title.ID)) {
			continue
		}
		enableGet := true
		for _, name := range strings.Split(title.NeedCard, ",") {
			pid := common.GetMPropIdByName(name)
			if pid > 0 && pid2record[pid] == 0 {
				enableGet = false
				break
			}
		}
		if !enableGet {
			continue
		} else {
			us.GetDb().Exec("update player_ext set F_Has_Title=? where uid =?", strings.Join(append(strings.Split(userInfo.FHasTitle, ","), strconv.Itoa(title.ID)), ","), user.ID)
			AnnounceAll(user.Nickname, "获得了新的称号-----"+title.Name)
		}
	}
}

func (us *NpcServices) GetCardTitleDatas(userId int) []gin.H {
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	hasTitles := strings.Split(userInfo.FHasTitle, ",")
	datas := []gin.H{}
	var data gin.H
	for _, title := range common.GetAllCardTitle() {
		data = gin.H{
			"id":         title.ID,
			"name":       title.Name,
			"img":        title.Img,
			"get_method": title.NeedDep,
			"effect":     title.EffectDep,
		}

		if userInfo.NowAchievementTitle == title.CodeName {
			data["status"] = 1
		} else if com.IsSliceContainsStr(hasTitles, strconv.Itoa(title.ID)) {
			data["status"] = 0
		} else {
			data["status"] = -1
		}
		datas = append(datas, data)
	}
	return datas
}

func (us *NpcServices) UseCardTitle(userId, titleId int) (bool, string) {

	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	hasTitles := strings.Split(userInfo.FHasTitle, ",")
	if !com.IsSliceContainsStr(hasTitles, strconv.Itoa(titleId)) {
		return false, "您还没有获得此称号！"
	}
	title := common.GetCardTitleById(titleId)
	if title == nil {
		return false, "称号不存在！"
	}
	pets := us.OptSvc.PetSrv.GetAllPets(userId)
	for _, p := range pets {
		us.OptSvc.FightSrv.DelZbAttr(p.ID)
	}
	us.GetDb().Exec("update player_ext set now_Achievement_title=? where uid =?", title.CodeName, userId)
	return true, "使用称号成功！"
}

func (us *NpcServices) CancelCardTitle(userId, titleId int) (bool, string) {

	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	hasTitles := strings.Split(userInfo.FHasTitle, ",")
	if !com.IsSliceContainsStr(hasTitles, strconv.Itoa(titleId)) {
		return false, "您还没有获得此称号！"
	}

	title := common.GetCardTitleById(titleId)
	if title == nil {
		return false, "称号不存在！"
	}
	if userInfo.NowAchievementTitle != title.CodeName {
		return false, "您未使用此称号！"
	}
	pets := us.OptSvc.PetSrv.GetAllPets(userId)
	for _, p := range pets {
		us.OptSvc.FightSrv.DelZbAttr(p.ID)
	}
	us.GetDb().Exec("update player_ext set now_Achievement_title='' where uid =?", userId)
	return true, "取消称号成功！"
}

// 扫雷

// 扫雷奖励信息
func (us *NpcServices) GetUserSaoleiAward(userId int) map[int]*models.SaoLeiAwardInfo {

	awards, err := rcache.GetSaoleiAward(userId)
	if err == nil {
		return awards
	}
	for i := 1; i < 10; i++ {
		prizeSetting := common.GetWelcome(fmt.Sprintf("sl_prize_best_%d", i))
		if prizeSetting != nil {
			prizes := strings.Split(prizeSetting.Content, ",")
			prizeId := prizes[rand.Intn(len(prizes))]
			if prop := common.GetMProp(com.StrTo(prizeId).MustInt()); prop != nil {
				awards[i] = &models.SaoLeiAwardInfo{Id: prop.ID, Name: prop.Name, Img: prop.Img}
			}
		}
	}
	rcache.SetSaoleiAward(userId, awards)
	return awards
}

// 扫雷-刷新奖励
func (us *NpcServices) UpdateSaoLeiAward(userId int) (bool, string) {
	if !us.OptSvc.PropSrv.DecrPropByPid(userId, 4019, 1) {
		return false, "没有刷新卡了！"
	}
	rcache.DelSaoleiAward(userId)
	return true, "刷新成功！"
}

func (us *NpcServices) UpdateSaoLeiLevel(userId, newLevel int) {
	us.GetDb().Exec("update player_ext set F_saolei_points=? where uid=?", newLevel, userId)
}

func (us *NpcServices) UpdateSaoLeiAddLevel(userId int) {
	us.GetDb().Exec("update player_ext set F_saolei_points=F_saolei_points+1 where uid=?", userId)
}

// 扫雷-开始扫雷
func (us *NpcServices) StartSaoLei(userId, position int) (string, gin.H) {
	result := gin.H{"enbale_sl": true, "result": true}
	level, enableSaolei := us.OptSvc.UserSrv.GetSaoleiStatus(userId)
	result["level"] = level
	if !enableSaolei {
		result["enbale_sl"] = false
		result["result"] = false
		return "您已没有扫雷资格，是否消耗闯关卡进入扫雷", result
	}
	bestAwards := us.GetUserSaoleiAward(userId)
	otherAwardsSetting := common.GetWelcome(fmt.Sprintf("sl_prize_other_%d", level))

	allOtherAwards := []gin.H{}

	successRateSetting := common.GetWelcome(fmt.Sprintf("sl_probability_%d", level))
	luckNum := rand.Intn(90) + 1
	successRateItems := strings.Split(successRateSetting.Content, ",")

	goodFlag := false
	dieFlag := false
	for _, rateStr := range successRateItems {
		rateItems := strings.Split(rateStr, ":")
		rates := strings.Split(rateItems[1], "-")
		if rateItems[0] == "good" {
			if luckNum >= com.StrTo(rates[0]).MustInt() && luckNum <= com.StrTo(rates[1]).MustInt() {
				goodFlag = true
				break
			}
		} else if rateItems[0] == "die" {
			if luckNum >= com.StrTo(rates[0]).MustInt() && luckNum <= com.StrTo(rates[1]).MustInt() {
				dieFlag = true
				break
			}
		}
	}
	bestNum := 1
	dieNum := level - 1
	otherNum := 9 - bestNum - dieNum
	var resultInfo gin.H

	if goodFlag {
		// 获得最好的东西

		getPid := bestAwards[level].Id
		if !us.OptSvc.PropSrv.AddProp(userId, getPid, 1, true) {
			result["result"] = false
			return "背包空间不足！", result
		}
		if level < 9 {
			us.UpdateSaoLeiAddLevel(userId)
		} else {
			us.UpdateSaoLeiLevel(userId, 1)
			rcache.DelSaoleiAward(userId)
		}
		mprop := common.GetMProp(getPid)
		user := us.OptSvc.UserSrv.GetUserById(userId)
		bestNum -= 1
		resultInfo = gin.H{"die": false, "name": mprop.Name, "img": mprop.Img}
		SelfGameLog(userId, fmt.Sprintf("扫雷:通过第%d关,获得极品奖励：%s", level, mprop.Name), 254)
		AnnounceAll(user.Nickname, fmt.Sprintf(" 通过扫雷第%d关,得到本关最极品奖励:%s", level, mprop.Name))
	} else if dieFlag {
		// 踩到地雷
		dieNum -= 1
		resultInfo = gin.H{"die": true}
		rcache.SetSaoleiTodayUser(userId, 1)
		rcache.SetSaoleiTicketUser(userId, 0)
		rcache.SetSaoleiDieUserLevel(userId, level)
		us.UpdateSaoLeiLevel(userId, 1)
	}

	otherLuckeyNum := rand.Intn(100) + 1
	getPid := 0
	for _, s := range strings.Split(otherAwardsSetting.Content, ",") {
		items := strings.Split(s, ":")
		randItems := strings.Split(items[1], "-")
		if !goodFlag && !dieFlag && otherLuckeyNum >= com.StrTo(randItems[0]).MustInt() && otherLuckeyNum <= com.StrTo(randItems[1]).MustInt() {
			if !us.OptSvc.PropSrv.AddProp(userId, com.StrTo(items[0]).MustInt(), 1, true) {
				//return false
				result["result"] = false
				return "背包空间不足！", result
			} else {
				// 获取普通奖励，在这里处理
				getPid = com.StrTo(items[0]).MustInt()
				mprop := common.GetMProp(getPid)
				otherNum -= 1
				resultInfo = gin.H{"die": false, "name": mprop.Name, "img": mprop.Img}
				us.UpdateSaoLeiAddLevel(userId)
				SelfGameLog(userId, fmt.Sprintf("扫雷:通过第%d关,获得普通奖励：%s", level, mprop.Name), 254)
			}
		} else {
			mprop := common.GetMProp(com.StrTo(items[0]).MustInt())
			allOtherAwards = append(allOtherAwards, gin.H{"die": false, "name": mprop.Name, "img": mprop.Img, "id": mprop.ID})
		}
	}
	resultAwards := make(map[int]gin.H)
	otherIndexs := []int64{}
	resultIndex := []int{}
	for i := 1; i <= 9; i++ {
		if i == position {
			continue
		}
		resultIndex = append(resultIndex, i)
		if bestNum > 0 {
			mprop := common.GetMProp(bestAwards[level].Id)
			resultAwards[i] = gin.H{"die": false, "name": mprop.Name, "img": mprop.Img, "id": mprop.ID}
			bestNum -= 1
			continue
		}
		if otherNum > 0 {
			randIndex := int64(rand.Intn(len(allOtherAwards)))
			for com.IsSliceContainsInt64(otherIndexs, randIndex) {
				randIndex = int64(rand.Intn(len(allOtherAwards)))
			}
			otherIndexs = append(otherIndexs, randIndex)
			resultAwards[i] = allOtherAwards[randIndex]
			otherNum -= 1
			continue
		}
		if dieNum > 0 {
			resultAwards[i] = gin.H{"die": true}
		}
	}

	// 打乱结果
	for i := 0; i < len(resultIndex); i++ {
		randIndex := rand.Intn(len(resultIndex))
		resultAwards[resultIndex[randIndex]], resultAwards[resultIndex[i]] = resultAwards[resultIndex[i]], resultAwards[resultIndex[randIndex]]
	}

	resultAwards[position] = resultInfo
	result["result_awards"] = resultAwards
	result["get_fhk"] = false
	if rand.Intn(30)+1 == 30 {
		if us.OptSvc.PropSrv.AddProp(userId, 4038, 1, true) {
			result["get_fhk"] = true
		}
	}
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	result["level"] = userInfo.FSaoleiPoints
	return "", result
}

// 扫雷-开始闯关
func (us *NpcServices) UseSaoleiTicketInto(userId int) (bool, string) {
	if _, ok := us.OptSvc.UserSrv.GetSaoleiStatus(userId); ok {
		return false, "无需使用闯关卡！"
	}
	if !us.OptSvc.PropSrv.DecrPropByPid(userId, 4045, 1) {
		return false, "闯关卡数量卡不足！"
	}
	rcache.SetSaoleiTicketUser(userId, 1)
	return true, "使用闯关卡成功！"
}

// 扫雷-复活
func (us *NpcServices) EasterSaoLei(userId int) (bool, string) {
	lastLevel, err := rcache.GetSaoleiDieUserUserLevel(userId)
	if err != nil {
		return false, "使用复活卡失败，玩家并未死亡!"
	}
	if lastLevel >= 1 && lastLevel <= 9 {
		if !us.OptSvc.PropSrv.DecrPropByPid(userId, 4038, 1) {
			return false, "复活卡数量不足！"
		}
		us.UpdateSaoLeiLevel(userId, lastLevel)
		rcache.SetSaoleiTicketUser(userId, 1)
		return true, "复活成功！"
	}
	rcache.DelSaoleiDieUserUserLevel(userId)
	return false, "使用复活卡失败，玩家并未死亡!"

}

// 扫雷道具信息
// 返回：扫雷闯关卡、复活卡、刷新卡数量
func (us *NpcServices) GetSaoleiPropNum(userId int) (int, int, int) {
	cgkId, fhkId, sxkId := 4045, 4038, 4019
	cgkSum, fhkSum, sxkSum := 0, 0, 0
	props := []models.UProp{}
	us.GetDb().Where("pid in (?, ?, ?) and uid = ? and sums>0", cgkId, fhkId, sxkId, userId).Find(&props)
	for _, prop := range props {
		if prop.Pid == cgkId {
			cgkSum = prop.Sums
		} else if prop.Pid == fhkId {
			fhkSum = prop.Sums
		} else if prop.Pid == sxkId {
			sxkSum = prop.Sums
		}
	}
	return cgkSum, fhkSum, sxkSum
}

// 皇宫奖励信息
func (us *NpcServices) KingAwards() gin.H {
	awardSetting := common.GetWelcome("holiday_prize")
	if awardSetting == nil {
		return gin.H{}
	}
	awardData := gin.H{}
	awards := strings.Split(awardSetting.Content, "|")
	if awards[0] == "0" {
		awardData["day"] = []gin.H{}
	} else {
		dayAward := []gin.H{}
		pItems := strings.Split(awards[0], ",")
		for _, pitem := range pItems {
			items := strings.Split(pitem, "*")
			prop := common.GetMProp(com.StrTo(items[0]).MustInt())
			if prop != nil && len(items) > 1 {
				dayAward = append(dayAward, gin.H{
					"name":    prop.Name,
					"vary_id": prop.VaryName,
					"num":     com.StrTo(items[1]).MustInt(),
					"id":      prop.ID,
				})
			}
		}
		awardData["day"] = dayAward
	}

	if awards[1] == "0" {
		awardData["week"] = []gin.H{}
	} else {
		weekAward := []gin.H{}
		pItems := strings.Split(awards[1], ",")
		for _, pitem := range pItems {
			items := strings.Split(pitem, "*")
			prop := common.GetMProp(com.StrTo(items[0]).MustInt())
			if prop != nil && len(items) > 1 {
				weekAward = append(weekAward, gin.H{
					"name":    prop.Name,
					"vary_id": prop.VaryName,
					"num":     com.StrTo(items[1]).MustInt(),
					"id":      prop.ID,
				})
			}
		}
		awardData["week"] = weekAward
	}
	pItems := strings.Split(awards[0], ";")
	now := time.Now()
	holidayAward := []gin.H{}
	for _, s := range pItems {
		dateItems := strings.Split(s, ":")

		if com.StrTo(dateItems[0][:4]).MustInt() == now.Year() && com.StrTo(dateItems[0][4:6]).MustInt() == int(now.Month()) && com.StrTo(dateItems[0][6:8]).MustInt() == now.Day() {
			pItems := strings.Split(dateItems[1], ",")
			for _, pitem := range pItems {
				items := strings.Split(pitem, "*")
				prop := common.GetMProp(com.StrTo(items[0]).MustInt())
				if prop != nil && len(items) > 1 {
					holidayAward = append(holidayAward, gin.H{
						"name":    prop.Name,
						"vary_id": prop.VaryName,
						"num":     com.StrTo(items[1]).MustInt(),
						"id":      prop.ID,
					})
				}
			}
		}
	}
	awardData["holiday"] = holidayAward

	return awardData

}

// 皇宫-砸蛋
func (us *NpcServices) Zadan(userId, position int, danType string) (bool, string, int, []gin.H) {
	var pid int
	var key, DName string
	awards := []gin.H{}
	if danType == "1" {
		// 金蛋id:3757
		pid = 3757
		key = "golden_eggs"
		DName = "金蛋"
	} else if danType == "2" {
		// 银蛋id :3758
		pid = 3758
		key = "silver_eggs"
		DName = "银蛋"
	} else if danType == "3" {
		// 铜蛋id :3759
		pid = 3759
		key = "copper_eggs"
		DName = "铜蛋"
	} else {
		return false, "参数出错！", 0, awards
	}
	prop := us.OptSvc.PropSrv.GetPropByPid(userId, pid, false)
	if prop == nil || prop.Sums == 0 {
		return false, "蛋券数量不足！", 0, awards
	}
	us.OptSvc.Begin()
	defer us.OptSvc.Rollback()
	if !us.OptSvc.PropSrv.DecrPropById(prop.ID, 1) {
		return false, "蛋券数量不足！", 0, awards
	}
	eggSetting := common.GetWelcome(key)
	if eggSetting == nil {
		return false, "设置出错！", prop.Sums, awards
	}
	var luckeyNum int
	unLuckeyNum := rand.Intn(101)
	if unLuckeyNum >= 85 {
		if danType == "1" {
			// 金蛋, 15%概率一定为玉露
			luckeyNum = 3001
		} else if danType == "2" {
			// 银蛋, 15%概率一定为500w月饼
			luckeyNum = 1001
		} else {
			luckeyNum = rand.Intn(10049)
		}
	} else {
		luckeyNum = rand.Intn(10049)
	}
	getPid := 0
	getNum := 0
	annouceFlag := false
	allAwards := []gin.H{}
	for _, s := range strings.Split(eggSetting.Content, ",") {
		if items := strings.Split(s, ":"); len(items) > 4 {
			pid := com.StrTo(items[0]).MustInt()
			if randItems := strings.Split(items[4], "-"); len(randItems) > 1 {
				if luckeyNum >= com.StrTo(randItems[0]).MustInt() && luckeyNum <= com.StrTo(randItems[1]).MustInt() {
					getPid = pid
					getNum = com.StrTo(items[1]).MustInt()
					if com.StrTo(items[2]).MustInt() == 1 {
						annouceFlag = true
					}
				} else {
					if com.StrTo(items[3]).MustInt() == 1 {
						mprop := common.GetMProp(pid)
						allAwards = append(allAwards, gin.H{"name": mprop.Name, "num": com.StrTo(items[1]).MustInt()})
					}
				}
			}
		}
	}
	if getPid == 0 {
		return false, "砸蛋结果为空，返回蛋券！", prop.Sums, awards
	}
	if !us.OptSvc.PropSrv.AddProp(userId, getPid, getNum, true) {
		return false, "背包空间不足！", prop.Sums, awards
	}
	for i := 0; i < len(allAwards)-1; i++ {
		j := len(allAwards) - i
		_ranNum := rand.Intn(j)
		allAwards[_ranNum], allAwards[j-1] = allAwards[j-1], allAwards[_ranNum]
	}
	iflag := 0
	mprop := common.GetMProp(getPid)
	for i := 0; i < 5; i++ {
		if i == position {
			awards = append(awards, gin.H{"name": mprop.Name, "num": getNum})
		} else {
			awards = append(awards, allAwards[iflag])
			iflag++
		}
	}
	if annouceFlag {
		user := us.OptSvc.UserSrv.GetUserById(userId)
		AnnounceAll(user.Nickname, fmt.Sprintf("参加了幸运砸%s活动，并幸运的获得了%s %d个", DName, mprop.Name, getNum))
	}
	us.OptSvc.Commit()
	return true, fmt.Sprintf("获得了%s %d个", mprop.Name, getNum), prop.Sums - 1, awards
}

func (us *NpcServices) GetKingAwards(userId int, awardType string) (bool, string) {
	awardData := us.KingAwards()
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	now := time.Now()
	getAwardItems := strings.Split(userInfo.PrizeItems, "|")
	if awardType == "1" {
		// 领取日常奖励
		day_award_status := false
		//fmt.Printf("getAwardItems:%s\n", userInfo.PrizeItems)
		if getAwardItems[0] != "" {
			if lastPrizeDay, err := utils.YmdStrParseTime(getAwardItems[0]); err == nil {
				//fmt.Printf("last date:%d-%d-%d\n", lastPrizeDay.Year(), lastPrizeDay.Month(), lastPrizeDay.Day())
				//fmt.Printf("now date:%d-%d-%d\n", now.Year(), now.Month(), now.Day())
				if now.Year() == lastPrizeDay.Year() && now.Month() == lastPrizeDay.Month() && now.Day() == lastPrizeDay.Day() {
					day_award_status = true
				}
			} else {
				fmt.Printf("error:%s\n", err)
			}
		}
		if day_award_status {
			return false, "您今日已领取过奖励了！"
		}
		awards := awardData["day"].([]gin.H)
		us.OptSvc.Begin()
		defer us.OptSvc.Rollback()
		for _, a := range awards {
			if !us.OptSvc.PropSrv.AddProp(userId, a["id"].(int), a["num"].(int), true) {
				return false, "背包空间不足！"
			}
		}
		getAwardItems[0] = utils.TimeFormatYmd(now)
		us.GetDb().Model(userInfo).Update(gin.H{"prize_every_day": strings.Join(getAwardItems, "|")})
		us.OptSvc.Commit()
		return true, "领取日常奖励成功！"
	} else if awardType == "2" {
		// 领取周末奖励

		if !(now.Weekday() == 0 || now.Weekday() == 6) {
			return false, "今天不是周末！"
		}
		week_award_status := false
		if len(getAwardItems) > 1 && getAwardItems[1] != "" {
			if lastPrizeDay, err := utils.YmdStrParseTime(getAwardItems[1]); err == nil {
				year, week := lastPrizeDay.ISOWeek()
				nyear, nweek := now.ISOWeek()
				if year == nyear && week == nweek {
					week_award_status = true
				}
			}
		}
		if week_award_status {
			return false, "您本周已领取过奖励了！"
		}
		awards := awardData["week"].([]gin.H)
		us.OptSvc.Begin()
		defer us.OptSvc.Rollback()
		for _, a := range awards {
			if !us.OptSvc.PropSrv.AddProp(userId, a["id"].(int), a["num"].(int), true) {
				return false, "背包空间不足！"
			}
		}
		if len(getAwardItems) < 2 {
			us.GetDb().Model(userInfo).Update(gin.H{"prize_every_day": getAwardItems[0] + "|" + utils.TimeFormatYmd(now)})
		} else {
			getAwardItems[1] = utils.TimeFormatYmd(now)
			us.GetDb().Model(userInfo).Update(gin.H{"prize_every_day": strings.Join(getAwardItems, "|")})
		}
		us.OptSvc.Commit()
		return true, "领取周末奖励成功！"

	} else if awardType == "3" {
		awards := awardData["holiday"].([]gin.H)
		if len(awards) == 0 {
			return false, "今天没有假日奖励可领取！"
		}
		holiday_award_status := false
		if len(getAwardItems) > 2 && getAwardItems[2] != "" {
			if lastPrizeDay, err := utils.YmdStrParseTime(getAwardItems[2]); err == nil {
				if now.Year() == lastPrizeDay.Year() && now.Month() == lastPrizeDay.Month() && now.Day() == lastPrizeDay.Day() {
					holiday_award_status = true
				}
			}
		}
		if holiday_award_status {
			return false, "您已领取过今日节假日奖励了！"
		}

		us.OptSvc.Begin()
		defer us.OptSvc.Rollback()
		for _, a := range awards {
			if !us.OptSvc.PropSrv.AddProp(userId, a["id"].(int), a["num"].(int), true) {
				return false, "背包空间不足！"
			}
		}
		if len(getAwardItems) < 3 {
			if len(getAwardItems) == 2 {
				us.GetDb().Model(userInfo).Update(gin.H{"prize_every_day": getAwardItems[0] + "|" + utils.TimeFormatYmd(now)})
			} else {
				us.GetDb().Model(userInfo).Update(gin.H{"prize_every_day": getAwardItems[0] + "||" + utils.TimeFormatYmd(now)})
			}

		} else {
			getAwardItems[2] = utils.TimeFormatYmd(now)
			us.GetDb().Model(userInfo).Update(gin.H{"prize_every_day": strings.Join(getAwardItems, "|")})
		}
		us.OptSvc.Commit()
		return true, "领取假日奖励成功！"
	}
	return false, "领取奖励参数出错！"
}

// 蛋券数量
// 返回：{"gold":   0,
//		"silver": 0,
//		"copper": 0,}
func (us *NpcServices) DanQuanCnt(userId int) gin.H {
	props := []models.UProp{}
	us.GetDb().Where("pid in (3757, 3758, 3759) and sums>0 and uid=?", userId).Find(&props)
	danquanData := gin.H{
		"gold":   0,
		"silver": 0,
		"copper": 0,
	}
	for _, prop := range props {
		if prop.Pid == 3757 {
			danquanData["gold"] = prop.Sums
		} else if prop.Pid == 3758 {
			danquanData["silver"] = prop.Sums
		} else if prop.Pid == 3759 {
			danquanData["copper"] = prop.Sums
		}
	}
	return danquanData
}

// 宠物神殿道具信息
func (us *NpcServices) GetPetSdPropInfo(userId int) gin.H {
	// 进化成长保护石Id：3501
	// 抽取成长道具：[3221, 3356, 3370, 3383]
	// 神圣进化添加物：effect 中含有zjsxdj_，varyname=7
	// 转生属性添加物：varyname=19
	// 神圣转生添加物：varyname=23
	// 合成、转生添加物：varyname=8
	sdData := gin.H{}
	props := us.OptSvc.PropSrv.GetCarryProps(userId, false)
	jh_protect_props := []gin.H{}
	cq_props := []gin.H{}
	ss_jh_props := []gin.H{}
	zs_attr_props := []gin.H{}
	sszs_attr_props := []gin.H{}
	zs_protect_props := []gin.H{}
	hc_protect_props := []gin.H{}
	for _, prop := range props {
		prop.GetM()
		propData := gin.H{"name": prop.MModel.Name, "id": prop.ID, "sum": prop.Sums}
		if prop.Pid == 3501 {
			jh_protect_props = append(jh_protect_props, propData)
			continue
		}
		if com.IsSliceContainsInt64([]int64{3221, 3356, 3370, 3383}, int64(prop.Pid)) {
			cq_props = append(cq_props, propData)
			continue
		}
		if prop.MModel.VaryName == 7 && strings.Index(prop.MModel.Effect, "zjsxdj_") > -1 {
			ss_jh_props = append(ss_jh_props, propData)
			continue
		}
		if prop.MModel.VaryName == 19 {
			zs_attr_props = append(zs_attr_props, propData)
			continue
		}
		if prop.MModel.VaryName == 23 {
			sszs_attr_props = append(sszs_attr_props, propData)
			continue
		}
		if prop.MModel.VaryName == 8 && prop.MModel.Effect != "" {
			useAges := strings.Split(prop.MModel.Usages, ":")
			if useAges[0] == "涅盘" {
				zs_protect_props = append(zs_protect_props, propData)
				continue
			} else {
				hc_protect_props = append(hc_protect_props, propData)
			}

		}
	}
	sdData["jh_protect_props"] = jh_protect_props
	sdData["cq_props"] = cq_props
	sdData["ss_jh_props"] = ss_jh_props
	sdData["sszs_attr_props"] = sszs_attr_props
	sdData["zs_attr_props"] = zs_attr_props
	sdData["zs_protect_props"] = zs_protect_props
	sdData["hc_protect_props"] = hc_protect_props
	return sdData
}

// 铁匠铺

func (us *NpcServices) GetTjpShopGood(update bool) gin.H {

	if !update {
		goods, err := rcache.GetTJPGoods()
		if err == nil {
			return goods
		}
	}
	shopData := gin.H{}

	shopProps := []models.MProp{}
	us.GetDb().Where("(prestige>0 or buy>0) && varyname=9 && yb=0").Find(&shopProps)
	jbData := []gin.H{}
	wwData := []gin.H{}
	for _, prop := range shopProps {
		need := "无"
		if prop.Requires != "" {
			need = strings.ReplaceAll(strings.ReplaceAll(prop.Requires, "lv", "等级"), "wx", "五行")
		}
		propInfo := gin.H{
			"id":      prop.ID,
			"name":    prop.Name,
			"vary_id": prop.VaryName,
			"need":    need,
		}
		if prop.SellJb > 0 {
			propInfo["price"] = prop.SellJb
			jbData = append(jbData, propInfo)
		}
		if prop.Prestige > 0 {
			propInfo["price"] = prop.Prestige
			wwData = append(wwData, propInfo)
		}
	}
	shopData = gin.H{
		"jb_list": jbData,
		"ww_list": wwData,
	}
	rcache.SetTJPGoods(shopData)
	return shopData
}

// 获取装备分解次数
func (us *NpcServices) GetZbFJTimes(userId int) int {
	times, err := rcache.GetZbfjTimes(userId)
	if err != nil {
		return 5
	}

	return times
}

// 分解装备
func (us *NpcServices) FenjieZb(userId, propId int) (bool, string) {
	prop := us.OptSvc.PropSrv.GetProp(userId, propId, false)
	if prop == nil {
		return false, "道具不存在！"
	}
	if prop.GetM(); prop.MModel.VaryName != 9 || prop.Sums == 0 || prop.Zbing > 0 {
		return false, "道具不存在！"
	}
	lefttimes := us.GetZbFJTimes(userId)
	if lefttimes == 0 {
		return false, "今日分解次数已达上限！"
	}

	fjSetting := common.GetWelcome("biodegradable_equipment")
	fjPositions := strings.Split(fjSetting.Content, ",")
	if !com.IsSliceContainsStr(fjPositions, strconv.Itoa(prop.MModel.Position)) {
		// 可分解
		return false, "该装备不可分解！"
	}
	successRateSetting := common.GetWelcome(fmt.Sprintf("fj_%d_success_rate", prop.MModel.PropsColor))
	if successRateSetting == nil {
		return false, "该装备不可分解"
	}
	rateItems := strings.Split(successRateSetting.Content, ",")
	randNum := rand.Intn(100) + 1
	getPid := 0
	getNum := 0
	for _, itemStr := range rateItems {
		items := strings.Split(itemStr, ":")
		randItems := strings.Split(items[2], "-")
		if len(randItems) > 1 {
			if randNum >= com.StrTo(randItems[0]).MustInt() && randNum <= com.StrTo(randItems[1]).MustInt() {
				getPid = com.StrTo(items[0]).MustInt()
				numItems := strings.Split(items[1], "-")
				if len(numItems) > 1 {
					startNum := com.StrTo(numItems[0]).MustInt()
					endNum := com.StrTo(numItems[1]).MustInt()
					getNum = rand.Intn(endNum-startNum+1) + startNum
				} else {
					getNum = com.StrTo(numItems[0]).MustInt()
				}
			}
		}
	}
	if getPid == 0 {
		// 分解失败
		rcache.SetZbfjTimes(userId, lefttimes-1)
		SelfGameLog(userId, fmt.Sprintf("装备分解:失去物品id:%s,物品名称:%s,分解失败", prop.ID, prop.MModel.Name), 22)
		return false, "分解失败，失去装备 " + prop.MModel.Name
	} else {
		us.OptSvc.Begin()
		defer us.OptSvc.Rollback()
		if !us.OptSvc.PropSrv.DecrPropById(prop.ID, 1) {
			return false, "该道具不存在！"
		}
		if !us.OptSvc.PropSrv.AddProp(userId, getPid, getNum, true) {
			return false, "背包空间不足！"
		}
		rcache.SetZbfjTimes(userId, lefttimes-1)
		mprop := common.GetMProp(getPid)
		SelfGameLog(userId, fmt.Sprintf("装备分解:失去物品id:%s,物品名称:%s,分解失败,得到物品:%s*%d", prop.ID, prop.MModel.Name, mprop.Name, getNum), 22)
		us.OptSvc.Commit()
		return false, fmt.Sprintf("分解成功，获得道具 %s * %d", mprop.Name, getNum)
	}
}

// 强化装备内置成功率
var QiangHuaEquipSuccessRates = []string{"6,100", "6,300", "6,600", "5,1000", "5,1500", "5,2000", "4,3000", "4,3500", "4,5000", "3,7000", "3,10000", "3,15000", "2,20000", "2,30000", "1,50000"}

// 强化装备
func (us *NpcServices) QiangHuaEquip(userId, propId, fzPropId int) (bool, string) {
	prop := us.OptSvc.PropSrv.GetProp(userId, propId, false)
	prop.GetM()
	if prop.MModel.VaryName != 9 || prop.Sums == 0 || prop.Zbing > 0 {
		return false, "不存在该道具！"
	}
	if prop.MModel.PlusFlag != 1 {
		return false, "该道具不可强化！"
	}
	nowLevel := 0
	if prop.PlusTmsEft != "" {
		if items := strings.Split(prop.PlusTmsEft, ","); len(items) > 1 {
			nowLevel = com.StrTo(items[0]).MustInt() + 1
		}
	}
	if nowLevel >= 15 {
		return false, "该装备强化已达满级！"
	}
	us.OptSvc.Begin()
	defer us.OptSvc.Rollback()
	if prop.MModel.PlusPid != 0 && !us.OptSvc.PropSrv.DecrPropByPid(userId, prop.MModel.PlusPid, 1) {
		return false, "强化材料不足！"
	}
	randNum := rand.Intn(11)
	luckyNum := 6
	needMoney := 1000
	successItems := strings.Split(QiangHuaEquipSuccessRates[nowLevel], ",")
	if len(successItems) > 1 {
		luckyNum = com.StrTo(successItems[0]).MustInt()
		needMoney = com.StrTo(successItems[1]).MustInt()
	}

	logNote := fmt.Sprintf("强化装备：%s(%d,镶嵌效果：%s), 强化等级：%d->%d", prop.MModel.Name, prop.ID, prop.FHoleInfo, nowLevel, nowLevel+1)

	if !us.OptSvc.UserSrv.DecreaseJb(userId, needMoney) {
		return false, "强化所需金币不足！"
	}

	baodengFlag := false // 失败时保存装备
	baodiFlag := false   // 失败时保存装备与属性
	if fzPropId != 0 {
		fzProp := us.OptSvc.PropSrv.GetProp(userId, fzPropId, false)
		if fzProp == nil || !us.OptSvc.PropSrv.DecrProp(userId, fzPropId, 1) {
			return false, "强化辅助道具不足！"
		}
		fzProp.GetM()
		if fzProp.MModel.Effect != "" {
			effectItems := strings.Split(fzProp.MModel.Effect, ":")
			if effectItems[0] == "suc" {
				luckyNum += 1
			} else if effectItems[0] == "100suc" {
				if items := strings.Split(effectItems[1], ","); len(items) > 1 && nowLevel < com.StrTo(items[1]).MustInt() {
					luckyNum = 10
				}
			} else if effectItems[0] == "baodi" {
				baodiFlag = true
			} else if effectItems[0] == "baodeng" {
				baodengFlag = true
			}
		}
		logNote += fmt.Sprintf(", 辅助道具：%s", fzProp.MModel.Name)
	}
	QhEffects := strings.Split(prop.MModel.PlusGet, ",")
	resultMsg := ""
	if randNum <= luckyNum {
		// 强化成功
		us.GetDb().Model(prop).Update(gin.H{"plus_tms_eft": fmt.Sprintf("%d,%s", nowLevel, QhEffects[nowLevel])})
		resultMsg = "强化结果：成功！"
	} else {
		// 强化失败
		if baodiFlag {
			// 强化降级
			nowLevel -= 2
			if nowLevel > 0 {
				us.GetDb().Model(prop).Update(gin.H{"plus_tms_eft": fmt.Sprintf("%d,%s", nowLevel, QhEffects[nowLevel])})
			} else {
				us.GetDb().Model(prop).Update(gin.H{"plus_tms_eft": ""})
			}
			resultMsg = "强化结果：失败！装备保留，强化属性降两级"

		} else if !baodengFlag {
			// 删除装备
			us.GetDb().Delete(prop)
			resultMsg = "强化结果：失败！装备消失"
		} else {
			resultMsg = "强化结果：失败！装备保留"
		}
	}
	logNote += ", " + resultMsg
	SelfGameLog(userId, logNote, 5)
	us.OptSvc.Commit()
	return true, resultMsg

}

// 装备的强化要求
func (us *NpcServices) QiangHuaInfo(userId, propId int) (gin.H, string) {
	prop := us.OptSvc.PropSrv.GetProp(userId, propId, false)
	result := gin.H{"enable_qh": false, "prop_name": "", "jb": 0}
	if prop == nil {
		return result, "道具不存在！"
	}
	prop.GetM()
	if prop.MModel.PlusFlag == 1 {
		nowLevel := 0
		if prop.PlusTmsEft != "" {
			if items := strings.Split(prop.PlusTmsEft, ","); len(items) > 1 {
				nowLevel = com.StrTo(items[0]).MustInt() + 1
			}
		}
		if nowLevel >= 15 {
			return result, "该装备强化已达满级！"
		}
		needMoney := 1000
		successItems := strings.Split(QiangHuaEquipSuccessRates[nowLevel], ",")
		if len(successItems) > 1 {
			needMoney = com.StrTo(successItems[1]).MustInt()
		}

		if prop.MModel.PlusPid != 0 {
			needProp := common.GetMProp(prop.MModel.PlusPid)
			result["prop_name"] = needProp.Name
		}
		result["jb"] = needMoney
		result["enable_qh"] = true
	}

	return result, ""
}

// 合成水晶、镶嵌装备
func (us *NpcServices) MergeProps(userId, id1, id2, fzid int) (bool, string) {
	prop1 := us.OptSvc.PropSrv.GetProp(userId, id1, false)
	var prop2 *models.UProp
	if id1 == id2 {
		prop2 = prop1
	} else {
		prop2 = us.OptSvc.PropSrv.GetProp(userId, id2, false)
	}
	if prop1 == nil || prop2 == nil || prop1.Sums < 1 || prop2.Sums < 1 {
		return false, "选取道具不存在！"
	}
	prop1.GetM()
	prop2.GetM()
	if prop1.MModel.VaryName == 25 && prop2.MModel.VaryName == 25 {
		// 合成水晶
		if prop1.Pid != prop2.Pid {
			return false, "合成水晶材料必须相同！"
		}
		effectItems := strings.Split(prop1.MModel.Effect, ",")
		if effectItems[0] == "full" {
			return false, "该道具已经满级，无法再进行合成！"
		}
		us.OptSvc.Begin()
		defer us.OptSvc.Rollback()
		if prop1.ID == prop2.ID {
			if prop1.Sums < 2 || !us.OptSvc.PropSrv.DecrPropById(prop1.ID, 2) {
				return false, "合成水晶材料数量不足！"
			}
		} else {
			if !us.OptSvc.PropSrv.DecrPropById(prop1.ID, 1) || !us.OptSvc.PropSrv.DecrPropById(prop2.ID, 1) {
				return false, "合成水晶材料数量不足！"
			}
		}
		baodiFlag := false
		AnnouceFlag := false

		logNote := fmt.Sprintf("水晶合成：添加物1 %s，添加物2 %s", prop1.MModel.Name, prop2.MModel.Name)
		level := com.StrTo(strings.Split(prop1.MModel.Name, "级")[0]).MustInt()
		if level >= 3 {
			AnnouceFlag = true
		}
		if fzid != 0 {
			fzProp := us.OptSvc.PropSrv.GetProp(userId, fzid, false)
			if fzProp == nil || fzProp.Sums < 1 {
				return false, "添加辅助材料数量不足！"
			}
			fzProp.GetM()
			fzEffectItems := strings.Split(fzProp.MModel.Effect, ":")
			if len(fzEffectItems) < 2 || fzEffectItems[0] != "bd" {
				return false, "添加辅助材料无效！"
			}
			items := strings.Split(fzEffectItems[1], "-")
			if level == 0 {
				return false, "此添加辅助材料不起作用，请更换或清除辅助材料！"
			}
			if len(items) == 2 {
				if !(level >= com.StrTo(items[0]).MustInt() && level <= com.StrTo(items[1]).MustInt()) {
					return false, "此添加辅助材料不起作用，请更换或清除辅助材料！"
				}
			} else {
				if level != com.StrTo(items[0]).MustInt() {
					return false, "此添加辅助材料不起作用，请更换或清除辅助材料！"
				}
			}
			baodiFlag = true
			logNote += fmt.Sprintf(", 添加辅助物：%s", fzProp.MModel.Name)
		}

		mergeItems := strings.Split(effectItems[0], ":")
		if len(mergeItems) < 3 {
			return false, "合成出错！道具无法合成！"
		}
		successRate := com.StrTo(strings.ReplaceAll(mergeItems[1], "%", "")).MustInt()
		randNum := rand.Intn(100) + 1
		var resultMsg string
		if randNum <= successRate {
			// 合成成功
			newPropId := com.StrTo(mergeItems[2]).MustInt()
			if !us.OptSvc.PropSrv.AddProp(userId, newPropId, 1, true) {
				return false, "背包空间不足！"
			}

			newProp := us.OptSvc.PropSrv.GetPropByPid(userId, newPropId, false)
			newProp.GetM()
			resultMsg = fmt.Sprintf("合成结果：成功合成 %s", newProp.MModel.Name)
			if AnnouceFlag {
				user := us.OptSvc.UserSrv.GetUserById(userId)
				color := ""
				if newProp.MModel.PropsColor == 3 {
					color = "red"
				} else if newProp.MModel.PropsColor == 4 {
					color = "green"
				} else if newProp.MModel.PropsColor == 5 {
					color = "#EDC028"
				}
				AnnounceAll(user.Nickname, fmt.Sprintf("成功合成<span style=color:%s><b>【<a onclick=showTip3(%d,0,1,2) onmouseout=UnTip3() style=cursor:pointer;color:%s;>%s</a>】</b></span>", color, newProp.ID, color, newProp.MModel.Name))
			}
		} else {
			if baodiFlag {
				us.OptSvc.PropSrv.AddPropSums(prop1.Sums, 1)
				resultMsg = fmt.Sprintf("合成结果：失败，保留道具 %s*1", prop1.MModel.Name)
			} else {
				resultMsg = fmt.Sprintf("合成结果：失败，添加道具消失")
			}
		}

		us.OptSvc.Commit()
		SelfGameLog(userId, logNote+", "+resultMsg, 5)
		return true, resultMsg
	} else if (prop1.MModel.VaryName == 25 && prop2.MModel.VaryName == 9) || (prop1.MModel.VaryName == 9 && prop2.MModel.VaryName == 25) {
		// 镶嵌水晶
		var zbProp, sjProp *models.UProp
		if prop1.MModel.VaryName == 9 {
			zbProp = prop1
			sjProp = prop2
		} else {
			zbProp = prop2
			sjProp = prop1
		}
		if zbProp.Zbing > 0 {
			return false, "该装备已装备在宠物身上，无法进行镶嵌！"
		}
		if zbProp.MModel.PlusNum == 0 {
			return false, "该装备水晶卡槽不足！！"
		}
		if zbProp.FHoleInfo != "" {
			holeInfos := strings.Split(zbProp.FHoleInfo, ",")
			if len(holeInfos) >= zbProp.MModel.PlusNum {
				return false, "该装备水晶卡槽不足！！"
			}
		}
		us.OptSvc.Begin()
		defer us.OptSvc.Rollback()
		if !us.OptSvc.PropSrv.DecrPropById(sjProp.ID, 1) {
			return false, "镶嵌水晶数量不足！！"
		}
		if sjProp.MModel.Requires != "" {
			requireItems := strings.Split(sjProp.MModel.Requires, ",")
			for _, require := range requireItems {
				items := strings.Split(require, ":")
				if items[0] == "postion" && len(items) > 1 {
					if !com.IsSliceContainsStr(strings.Split(items[1], "|"), strconv.Itoa(zbProp.MModel.Position)) {
						return false, "镶嵌水晶要求装备部位不符！"
					}
				}
			}
		}
		effectItems := strings.Split(sjProp.MModel.Effect, ",")
		if len(effectItems) < 2 {
			return false, "镶嵌道具不是水晶，无法进行镶嵌！"
		}
		effectItems = strings.Split(effectItems[1], ":")
		if len(effectItems) < 2 || effectItems[0] != "xq" {
			return false, "镶嵌道具不是水晶，无法进行镶嵌！"
		}
		effects := strings.Split(effectItems[1], "|")
		randNum := rand.Intn(100) + 1
		attrType := ""
		attrData := ""
		for _, str := range effects {
			items := strings.Split(str, "_")
			rateItems := strings.Split(items[2], "-")
			if randNum >= com.StrTo(rateItems[0]).MustInt() && randNum <= com.StrTo(rateItems[1]).MustInt() {
				attrType = items[0]
				attrData = items[1]
				break
			}
		}
		logNote := fmt.Sprintf("镶嵌装备：装备 %s(%d), 水晶 %s, ", zbProp.MModel.Name, zbProp.ID, sjProp.MModel.Name)
		if attrType != "" {
			us.GetDb().Model(zbProp).Update(gin.H{"F_item_hole_info": attrType + ":" + attrData})
			resultMsg := "镶嵌结果："
			switch attrType {
			case "ac":
				resultMsg += "攻击增加:" + attrData
				break
			case "crit":
				resultMsg += "会心一击发动几率增加:" + attrData
				break
			case "shjs":
				resultMsg += "伤害加深:" + attrData
				break
			case "dxsh":
				resultMsg += "伤害抵消:" + attrData
				break
			case "hp":
				resultMsg += "HP上限增加:" + attrData
				break
			case "mp":
				resultMsg += "MP上限增加:" + attrData
				break
			case "mc":
				resultMsg += "防御增加:" + attrData
				break
			case "hits":
				resultMsg += "命中增加:" + attrData
				break
			case "miss":
				resultMsg += "闪避增加:" + attrData
				break
			case "szmp":
				resultMsg += "伤害的" + attrData + "转化为mp"
				break
			case "sdmp":
				resultMsg += "伤害的" + attrData + "以mp抵消"
				break
			case "speed":
				resultMsg += "攻击速度:" + attrData
				break
			case "hitsmp":
				resultMsg += "命中吸取伤害的" + attrData + "转化为自身MP"
				break
			case "hitshp":
				resultMsg += "命中吸取伤害的" + attrData + "转化为自身HP"
				break
			}
			us.OptSvc.Commit()
			SelfGameLog(userId, logNote+resultMsg, 5)
			return true, resultMsg
		} else {
			return false, "水晶数据出错，无法进行镶嵌！"
		}
	} else {
		return false, "选取道具无法进行合成或镶嵌！"
	}
}

// 家族
func (us *NpcServices) GetAllFamilyData() []gin.H {
	datas := []gin.H{}
	var families []*struct {
		Id            int    `gorm:"column:id"`
		Name          string `gorm:"column:name"`
		Honor         int    `gorm:"column:honor"`
		Level         int    `gorm:"column:level"`
		Count         int    `gorm:"column:count"`
		PresidentName string `gorm:"column:president_name"`
	}
	us.GetDb().Raw(`select guild.id as id, guild.name as name, guild.honor as honor, guild.level as level, guild.number_of_member as count, player.nickname as president_name 
from guild 
left join player on guild.president_id=player.id 
order by guild.honor desc`).
		Scan(&families)
	for _, f := range families {
		datas = append(datas, gin.H{
			"id":        f.Id,
			"name":      f.Name,
			"president": f.PresidentName,
			"honor":     f.Honor,
			"level":     f.Level,
			"count":     f.Count,
		})
	}
	return datas
}

func (us *NpcServices) GetMyFamily(userId int) *models.Family {
	member := us.GetFamilyMember(userId)
	if member == nil {
		return nil
	}
	return us.GetFamily(member.FamilyId)
}

func (us *NpcServices) GetFamilyMember(userId int) *models.FamilyMember {
	member := &models.FamilyMember{}
	us.GetDb().Where("member_id=?", userId).First(member)
	if member.UserId == 0 {
		return nil
	}
	return member
}

func (us *NpcServices) GetFamily(familyId int) *models.Family {
	family := &models.Family{}
	us.GetDb().Where("id=?", familyId).First(family)
	if family.Id == 0 {
		return nil
	}
	return family
}

func (us *NpcServices) GetFamilySet(level int) *models.FamilySetting {
	familySet := &models.FamilySetting{}
	us.GetDb().Where("level=?", level).First(familySet)
	if familySet.Id > 0 {
		return familySet
	}
	return nil
}

func (us *NpcServices) GetFamilyStoreData(userId, shopLevel int) []gin.H {
	datas := []gin.H{}

	props := []*models.MProp{}
	us.GetDb().Where("contribution > 0 or honor > 0 and guild_level<=?", shopLevel).Find(&props)
	for _, p := range props {
		datas = append(datas, gin.H{
			"id":           p.ID,
			"name":         p.Name,
			"vary_id":      p.VaryName,
			"honor":        p.Honor,
			"contribution": p.Contribution,
		})
	}
	return datas
}

func (us *NpcServices) GetFamilyInfo(familyId, userId int) gin.H {
	authority := -1
	family := us.GetFamily(familyId)
	if family == nil {
		return nil
	}
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority < 1 || member.FamilyId != familyId {
		authority = -1
	} else {
		authority = member.Authority
	}
	members := []*struct {
		Id        int    `gorm:"column:id" json:"id"`
		Nickname  string `gorm:"column:nickname" json:"nickname"`
		Level     int    `gorm:"column:level" json:"level"`
		Czl       string `gorm:"column:czl" json:"czl"`
		Authority int    `gorm:"column:authority" json:"authority"`
	}{}
	selectAuthority := 1
	if authority > 1 {
		selectAuthority = 0
	}
	us.GetDb().Raw(`SELECT player.id as id,player.nickname as nickname,userbb.level as level,userbb.czl as czl,guild_members.priv as authority 
FROM guild_members 
left join player on guild_members.member_id=player.id
left join userbb on player.mbid=userbb.id
WHERE guild_members.guild_id = ? and guild_members.priv>=? ORDER BY authority DESC`, familyId, selectAuthority).
		Scan(&members)
	president := ""
	for _, m := range members {
		if m.Authority == 3 {
			president = m.Nickname
			break
		}
	}

	deposits := []string{}
	familyBag := []*models.FamilyBag{}
	us.GetDb().Where("guild_id=?", familyId).Find(&familyBag)
	for _, bag := range familyBag {
		prop := common.GetMProp(bag.Pid)
		deposits = append(deposits, fmt.Sprintf("%s x %d", prop.Name, bag.Sums))
	}
	return gin.H{
		"id":          familyId,
		"name":        family.Name,
		"honor":       family.Honor,
		"president":   president,
		"deposits":    strings.Join(deposits, ","),
		"level":       family.Level,
		"create_time": utils.FormatTime(time.Unix(int64(family.CreateTime), 0)),
		"max_member":  us.GetFamilySet(family.Level).MaxMemberNumber,
		"introduce":   family.Info,
		"members":     members,
		"authority":   authority,
	}
}

func (us *NpcServices) CreateFamily(userId int, name, info string) (bool, string) {
	family := us.GetMyFamily(userId)
	if family != nil {
		return false, "您已经加入到其它家族，不能创建！"
	}
	user := us.OptSvc.UserSrv.GetUserById(userId)
	if user.Vip < 10 {
		return false, "您的积分不足10，不能创建！"
	}
	if !us.OptSvc.PropSrv.DecrPropByPid(userId, 2494, 1) {
		return false, "您没有家族令牌，不能创建！"
	}
	us.GetDb().Model(user).Update(gin.H{"vip": user.Vip - 10})
	now := utils.NowUnix()
	family = &models.Family{}
	us.GetDb().Model(family).Where("name=?", name).First(family)
	if family.Id > 0 {
		return false, "家族名称已存在！"
	}

	family = &models.Family{
		Name:           name,
		Info:           info,
		CreatorIdStr:   strconv.Itoa(userId),
		PresidentId:    userId,
		Honor:          0,
		Level:          1,
		ShopLevel:      1,
		NumberOfMember: 1,
		CreateTime:     now,
	}
	us.GetDb().Create(family)
	member := &models.FamilyMember{
		UserId:       userId,
		FamilyId:     family.Id,
		JoinTime:     now,
		Authority:    3,
		Contribution: 0,
		Honor:        0,
	}
	us.GetDb().Create(member)
	return true, "创建成功！"
}

func (us *NpcServices) ApplyFamily(userId, familyId int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member != nil {
		if member.Authority == 0 {
			return false, "您已经申请加入其它或此家族，不能再申请！"
		} else if member.Authority > 0 {
			return false, "您已经加入到其它家族，不能创建！"
		}
	}
	family := &models.Family{}
	us.GetDb().Where("id=?", familyId).First(family)
	if family.Id == 0 {
		return false, "申请家族不存在！"
	}
	familySet := us.GetFamilySet(family.Id)
	if familySet.MaxMemberNumber <= family.NumberOfMember {
		return false, "家族人数已满，不能再申请！"
	}
	now := time.Now()
	member = &models.FamilyMember{
		UserId:       userId,
		FamilyId:     familyId,
		JoinTime:     int(now.Unix()),
		Authority:    0,
		Contribution: 0,
		Honor:        0,
	}
	us.GetDb().Create(member)
	user := us.OptSvc.UserSrv.GetUserById(userId)
	msg := fmt.Sprintf("%s 请求加入您的家族，请速去处理吧！", user.Nickname)
	familyManagers := []*models.FamilyMember{}
	us.GetDb().Where("priv>2 and guild_id=?", familyId).Find(&familyManagers)
	for _, m := range familyManagers {
		us.OptSvc.SysSrv.SendEmail(m.UserId, now, msg)
	}
	return true, "申请成功！正在等待回复！"
}

func (us *NpcServices) ReplyApplyFamily(authorityId, applyId int, pass bool) (bool, string) {
	authorityMember := us.GetFamilyMember(authorityId)
	if authorityMember == nil || authorityMember.Authority == 0 {
		return false, "您还没有加入家族！"
	}
	if authorityMember.Authority < 2 {
		return false, "您没有权利处理申请！"
	}
	applyMember := us.GetFamilyMember(applyId)
	if applyMember == nil || authorityMember.FamilyId != applyMember.FamilyId {
		return false, "该玩家并没有申请加入您的家族！"
	}
	if applyMember.Authority != 0 {
		return false, "该玩家已经进入家族中！"
	}
	if pass {
		us.OptSvc.SysSrv.SendEmail(applyId, time.Now(), fmt.Sprintf("家族【%s】通过您的入会请求！"))
		us.GetDb().Model(applyMember).Update(gin.H{"priv": 1})
		us.GetDb().Model(&models.Family{}).Where("id=?", authorityMember.FamilyId).Update(gin.H{"number_of_member": gorm.Expr("number_of_member+1")})
		return true, "成功批准申请！"
	} else {
		us.OptSvc.SysSrv.SendEmail(applyId, time.Now(), fmt.Sprintf("家族【%s】拒绝了您的入会请求！"))
		us.GetDb().Delete(applyMember)
		return true, "成功拒绝申请！"
	}
}

func (us *NpcServices) ManageAuthority(authorityId, changeId int, authority int) (bool, string) {
	authorityMember := us.GetFamilyMember(authorityId)
	if authorityMember == nil || authorityMember.Authority == 0 {
		return false, "您还没有加入家族！"
	}
	if authorityId == changeId {
		return false, "您无法管理自身！"
	}
	changeMember := us.GetFamilyMember(changeId)
	if changeMember == nil || authorityMember.FamilyId != changeMember.FamilyId {
		return false, "该玩家并没有加入您的家族！"
	}
	if authorityMember.Authority <= changeMember.Authority {
		return false, "您没有资格管理该玩家!"
	}
	if authorityMember.Authority < authority {
		return false, "您没有资格做出管理！"
	}
	if changeMember.Authority == authority {
		return false, "不可重复设置权限！"
	}
	if authority == 3 {
		pCnt := 1
		us.GetDb().Model(&models.FamilyMember{}).Where("priv=3").Count(&pCnt)
		if pCnt >= 2 {
			return false, "家族会长不能超过2个！"
		}
		us.GetDb().Model(changeMember).Update(gin.H{"priv": 3})
		us.OptSvc.SysSrv.SendEmail(changeId, time.Now(), fmt.Sprintf("【家族】您被授权为家族会长！"))
	} else if authority == 2 {
		pCnt := 1
		us.GetDb().Model(&models.FamilyMember{}).Where("priv=2").Count(&pCnt)
		if pCnt >= 4 {
			return false, "家族长老不能超过4个！"
		}
		us.GetDb().Model(changeMember).Update(gin.H{"priv": 2})
		us.OptSvc.SysSrv.SendEmail(changeId, time.Now(), fmt.Sprintf("【家族】您被授权为家族长老！"))
	} else if authority == 1 {
		us.GetDb().Model(changeMember).Update(gin.H{"priv": 1})
		us.OptSvc.SysSrv.SendEmail(changeId, time.Now(), fmt.Sprintf("【家族】您被降为普通会员！"))
		return true, "权限调整成功！"
	}
	return false, "参数错误！"
}

func (us *NpcServices) FireFamilyMember(authorityId, changeId int) (bool, string) {

	authorityMember := us.GetFamilyMember(authorityId)
	if authorityMember == nil || authorityMember.Authority == 0 {
		return false, "您还没有加入家族！"
	}
	if authorityId == changeId {
		return false, "您无法管理自身！"
	}
	changeMember := us.GetFamilyMember(changeId)
	if changeMember == nil || authorityMember.FamilyId != changeMember.FamilyId || authorityMember.Authority < 1 {
		return false, "该玩家并没有加入您的家族！"
	}
	if authorityMember.Authority <= changeMember.Authority {
		return false, "您没有资格管理该玩家!"
	}
	family := us.GetFamily(authorityMember.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	us.OptSvc.SysSrv.SendEmail(changeId, time.Now(), fmt.Sprintf("您被从家族【%s】移出！", family.Name))
	us.GetDb().Delete(changeMember)

	return true, "成功移出会员！"
}

func (us *NpcServices) ExitFamily(userId int) (bool, string) {

	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	msg := ""
	if member.Authority == 3 {
		pCnt := 1
		us.GetDb().Model(&models.FamilyMember{}).Where("priv=3").Count(&pCnt)
		if pCnt == 1 {
			return false, "家族会长只有一个，不可退出！"
		}
		if family.PresidentId == userId {
			member = &models.FamilyMember{}
			us.GetDb().Where("priv=? and guild_id=? and member_id!=?", 3, family.Id, userId).First(member)
			if member.UserId > 0 {
				us.GetDb().Model(family).Update(gin.H{"president_id": member.UserId})
			} else {
				return false, "家族会长只有一个，不可退出！"
			}
		}
		user := us.OptSvc.UserSrv.GetUserById(userId)
		msg = fmt.Sprintf("【家族】会长 %s 退出家族！", user.Nickname)
	}
	us.GetDb().Delete(member)
	if msg != "" {

		mids := []*struct {
			Id int
		}{}
		us.GetDb().Raw("select member_id as id from guild_members where guild_id=? and priv>0", family.Id).Scan(&mids)
		now := time.Now()
		for _, id := range mids {
			if id.Id == userId {
				continue
			}
			us.OptSvc.SysSrv.SendEmail(id.Id, now, msg)
		}
	}
	us.GetDb().Model(family).Update(gin.H{"number_of_member": gorm.Expr("number_of_member-1")})
	return true, "退出家族成功！"

}

func (us *NpcServices) DisbandFamily(userId int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}
	if member.Authority < 3 {
		return false, "您没有权限解散家族！"
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	mids := []*struct {
		Id int
	}{}
	us.GetDb().Raw("select member_id as id from guild_members where guild_id=? and priv>0", family.Id).Scan(&mids)
	now := time.Now()
	msg := fmt.Sprintf("您所在的家族【%s】已被解散！", family.Name)
	for _, id := range mids {
		if id.Id == userId {
			continue
		}
		us.OptSvc.SysSrv.SendEmail(id.Id, now, msg)
	}
	us.GetDb().Where("guild_id=?", family.Id).Delete(&models.FamilyMember{})
	us.GetDb().Where("id=?", family.Id).Delete(&models.Family{})
	us.GetDb().Where("guild_id=?", family.Id).Delete(&models.FamilyBag{})
	// 还要清理战场记录
	return true, "成功解散家族！"
}

func (us *NpcServices) GetFamilyUpgradeInfo(userId int) (info gin.H, msg string) {
	info = gin.H{
		"next_level":  0,
		"need_honor":  0,
		"need_member": 0,
		"need_prop":   nil,
	}
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		msg = "您还没有加入到家族中！"
		return
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		msg = "您还没有加入到家族中！"
		return
	}
	set := us.GetFamilySet(family.Level)
	info["next_level"] = family.Level + 1
	info["need_honor"] = set.NeedHonor
	info["need_member"] = set.NeedMemberNumber
	propInfo := []string{}
	familyBag := []*models.FamilyBag{}
	us.GetDb().Where("guild_id=?", family.Id).Find(&familyBag)
	for _, s := range strings.Split(set.NeedProps, ",") {
		items := strings.Split(s, "|")
		if len(items) > 1 {
			pid := com.StrTo(items[0]).MustInt()
			NeedSum := com.StrTo(items[1]).MustInt()
			HasSum := 0
			for _, bag := range familyBag {
				if bag.Pid == pid {
					HasSum = bag.Sums
					break
				}
			}
			prop := common.GetMProp(pid)
			propInfo = append(propInfo, fmt.Sprintf("%s:%d/%d", prop.Name, HasSum, NeedSum))
		}
	}
	info["need_prop"] = propInfo
	return
}

func (us *NpcServices) UpgradeFamily(userId int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}
	if member.Authority < 3 {
		return false, "您没有权限升级家族！"
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	set := us.GetFamilySet(family.Level)
	if family.NumberOfMember < set.NeedMemberNumber {
		return false, "家族成员不够!"
	}
	if family.Honor < set.NeedHonor {
		return false, "家族荣誉不够!"
	}

	familyBag := []*models.FamilyBag{}
	us.GetDb().Where("guild_id=?", family.Id).Find(&familyBag)
	for _, s := range strings.Split(set.NeedProps, ",") {
		items := strings.Split(s, "|")
		if len(items) > 1 {
			pid := com.StrTo(items[0]).MustInt()
			NeedSum := com.StrTo(items[1]).MustInt()
			HasSum := 0
			for _, bag := range familyBag {
				if bag.Pid == pid {
					HasSum = bag.Sums
					break
				}
			}
			if HasSum < NeedSum {
				return false, "物品不够!"
			}
		}
	}
	us.GetDb().Delete(familyBag)
	us.GetDb().Model(family).Update(gin.H{"level": family.Level + 1})
	return true, "升级成功！"
}

func (us *NpcServices) UpgradeFamilyStore(userId int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}
	if member.Authority < 3 {
		return false, "您没有权限升级家族商店！"
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	if family.ShopLevel >= family.Level {
		return false, "商店等级不能高于家族等级！"
	}

	us.OptSvc.Begin()
	defer us.OptSvc.Rollback()
	set := us.GetFamilySet(family.Level)

	for _, s := range strings.Split(set.NeedItemsForShop, ",") {
		items := strings.Split(s, ":")
		if len(items) > 1 {
			pid := com.StrTo(items[0]).MustInt()
			NeedSum := com.StrTo(items[1]).MustInt()
			if !us.OptSvc.PropSrv.DecrPropByPid(userId, pid, NeedSum) {
				return false, "物品不够!"
			}
		}
	}
	us.GetDb().Model(family).Update(gin.H{"shop_level": family.ShopLevel + 1})
	us.OptSvc.Commit()
	return true, "升级商店成功！"
}

func (us *NpcServices) PurchaseFamilyStoreProp(userId, propId, num int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}

	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	prop := common.GetMProp(propId)
	if prop == nil {
		return false, "商品不存在！"
	}
	if (prop.Honor == 0 && prop.Contribution == 0) || prop.GuildLevel > family.ShopLevel {
		return false, "商品不存在！"
	}
	if prop.Vary == 2 && num > 1 {
		return false, "不可叠加道具一次最多只能购买一个！"
	}
	us.OptSvc.Begin()
	defer us.OptSvc.Rollback()
	if us.GetDb().Model(member).Where("honor>=? and contribution>=?", prop.Honor*num, prop.Contribution*num).Update(gin.H{
		"honor":        gorm.Expr("honor-?", prop.Honor*num),
		"contribution": gorm.Expr("contribution-?", prop.Contribution*num),
	}).RowsAffected == 0 {
		return false, "您的荣誉或贡献不足！"
	}
	if !us.OptSvc.PropSrv.AddProp(userId, propId, num, true) {
		return false, "背包空间不足！"
	}
	return true, "购买成功！"
}

func (us *NpcServices) FamilyDonate(userId, propId, num int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	prop := us.OptSvc.PropSrv.GetProp(userId, propId, false)
	if prop == nil {
		return false, "所选道具不存在！"
	}
	if prop.Sums < num {
		return false, "所填道具数量过多！"
	}
	set := us.GetFamilySet(family.Level)
	familyBagProp := &models.FamilyBag{}
	us.GetDb().Where("guild_id=? and pid=?", family.Id, prop.Pid).First(&familyBagProp)
	needNum := 0
	giveHonor := 0
	for _, s := range strings.Split(set.NeedProps, ",") {
		items := strings.Split(s, "|")
		if len(items) > 1 {
			pid := com.StrTo(items[0]).MustInt()
			NeedSum := com.StrTo(items[1]).MustInt()
			if pid == prop.ID {
				HasSum := 0
				if familyBagProp.Id > 0 {
					HasSum = familyBagProp.Sums
				}
				if HasSum >= NeedSum {
					return false, "您要捐赠的物品已经足够了！"
				}
				needNum = needNum - HasSum
				giveHonor = com.StrTo(items[2]).MustInt()
				break
			}
		}
	}
	if needNum == 0 {
		return false, "升到下一级不需要您捐献这个物品！"
	}
	us.OptSvc.PropSrv.DecrProp(userId, prop.ID, num)
	us.GetDb().Model(member).Update(gin.H{"contribution": member.Contribution + num*giveHonor})
	if familyBagProp.Id > 0 {
		us.GetDb().Model(familyBagProp).Update(gin.H{"sums": gorm.Expr("sums+?", num)})
	} else {
		familyBagProp = &models.FamilyBag{
			FamilyId: family.Id,
			Pid:      prop.Pid,
			Sums:     num,
		}
	}
	return true, fmt.Sprintf("捐献成功！获得贡献：%d", num*giveHonor)
}

func (us *NpcServices) GetFamilyWelfare(userId int) (bool, string) {
	member := us.GetFamilyMember(userId)
	if member == nil || member.Authority == 0 {
		return false, "您还没有加入到家族中！"
	}
	family := us.GetFamily(member.FamilyId)
	if family == nil {
		return false, "您还没有加入到家族中！"
	}
	user := us.OptSvc.UserSrv.GetUserById(userId)
	if user.BagPlace-us.OptSvc.PropSrv.GetCarryPropsCnt(userId) < 3 {
		return false, "背包空间不足！请留下至少3格的空间"
	}
	userInfo := us.OptSvc.UserSrv.GetUserInfoById(userId)
	now := utils.TimeFormatYmd(time.Now())
	if userInfo.GetWelfareDate >= com.StrTo(now).MustInt() {
		return false, "今天已经领过一次奖励了！"
	}
	set := us.GetFamilySet(family.Level)
	msgs := []string{}
	for _, s := range strings.Split(set.Welfare, ",") {
		items := strings.Split(s, ":")
		if len(items) > 2 {
			pid := com.StrTo(items[0]).MustInt()
			Rate := com.StrTo(items[1]).MustInt()
			Num := com.StrTo(items[2]).MustInt()
			if utils.RandInt(1, Rate) == 1 {
				if p, ok := us.OptSvc.PropSrv.AddOrCreateProp(userId, pid, Num, false); ok {
					p.GetM()
					msgs = append(msgs, fmt.Sprintf(" %s %d个", p.MModel.Name, Num))
				}
			}
		}
	}
	us.GetDb().Model(userInfo).Update(gin.H{"get_welfare_time": now})
	return true, fmt.Sprintf("获得道具 %s", strings.Join(msgs, ","))

}
