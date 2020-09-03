package services

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"math/rand"
	"pokemon/common/rcache"
	"pokemon/game/common"
	"pokemon/game/models"
	"pokemon/game/repositories"
	common2 "pokemon/game/services/common"
	"pokemon/game/utils"
	"strconv"
	"strings"
)

type PetService struct {
	BaseService
	repo *repositories.PetRepositories
}

func NewPetService(osrc *OptService) *PetService {
	us := &PetService{repo: repositories.NewPetRepositories()}
	us.SetOptSrc(osrc)
	return us
}

// 初始化使用
func (ps *PetService) InitTable() {
	if !ps.repo.HasTable("model_pet") {
		fmt.Println("开始生成表")
		ps.repo.CreateTable(&models.MPet{})
		ps.repo.CreateTable(&models.MSkill{})
		ps.repo.CreateTable(&models.UPet{})
		fmt.Println("生成表结束")
	} else {
		ps.repo.AutoMigrate(&models.MPet{})
		ps.repo.AutoMigrate(&models.MSkill{})
		ps.repo.AutoMigrate(&models.UPet{})
	}
}

func (ps *PetService) GetAllMPet() *[]models.MPet {
	mpets, success := ps.repo.GetALLMPetsFromMysql()
	if success {
		return mpets
	} else {
		return nil
	}
}

func (ps *PetService) GetAllMSkill() *[]models.MSkill {
	mskill, success := ps.repo.GetALLMSkillsFromMysql()
	if success {
		return mskill
	} else {
		return nil
	}
}

func (ps *PetService) GetAllExpList() *[]models.ExpList {
	expList, success := ps.repo.GetALLExp2lvFromMysql()
	if success {
		return expList
	} else {
		return nil
	}
}

func (ps *PetService) GetAllGrowthList() *[]models.Growth {
	GrowthList, success := ps.repo.GetALLGrowthFromMysql()
	if success {
		return GrowthList
	} else {
		return nil
	}
}

func (ps *PetService) GetAllSSZsList() *[]models.SSzsRule {
	RoleList, success := ps.repo.GetALLSSZsFromMysql()
	if success {
		return RoleList
	} else {
		return nil
	}
}

func (ps *PetService) GetAllSSJhList() *[]models.SSjhRule {
	RoleList, success := ps.repo.GetALLSSJhFromMysql()
	if success {
		return RoleList
	} else {
		return nil
	}
}

// 固定数据查询

func (ps *PetService) GetExpByLv(level int) (int, bool) {
	exp := common2.GetNextExp(level)
	if exp == 0 {
		return 0, false
	} else {
		return exp, true
	}
}

func (ps *PetService) GetMskill(skillId int) *models.MSkill {
	mskill, err := rcache.GetMSkill(skillId)
	if err != nil {
		fmt.Println("发生错误：", err)
		return nil
	} else {
		return mskill
	}
}

func (ps *PetService) GetMpet(petId int) *models.MPet {
	mpet := &models.MPet{}
	ps.GetDb().Where("Id=?", petId).First(mpet)

	return mpet
}

func (ps *PetService) GetMpetByName(name string) *models.MPet {
	name = utils.ToGbk(name)
	mpet := &models.MPet{}
	ps.GetDb().Where("Name=?", name).First(mpet)
	return mpet
}

func (ps *PetService) GetMaxLevel(mpet *models.MPet) int {
	var maxLevel int
	if mpet.Wx != common.WxSS {
		maxLevel = common.MaxLevel
	} else {
		maxLevel = common2.GetSSJhRule(mpet.ID).MaxLevel
	}
	return maxLevel
}

// 操作游戏数据API

func (ps *PetService) GetPet(userId, petId int) *models.UPet {
	pet := &models.UPet{}
	ps.GetDb().Where(&models.UPet{ID: petId, Uid: userId}).First(pet)
	if pet.ID > 0 {
		return pet
	}
	return nil
}

func (ps *PetService) GetPetById(petId int) *models.UPet {
	pet := &models.UPet{}
	ps.GetDb().Where(&models.UPet{ID: petId}).First(pet)
	if pet.ID > 0 {
		return pet
	}
	return nil
}

func (ps *PetService) PutIn(userId, petId int) (result gin.H, msg string) {
	// 存到牧场
	pet := ps.GetPet(userId, petId)
	result = gin.H{"result": false}
	if pet == nil {
		msg = "宠物不存在！"
		return
	}
	if pet.Muchang != 0 {
		msg = "宠物已在牧场！"
		return
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	if user.Mbid == petId {
		msg = "不可寄养主宠！"
		return
	}
	if user.McPlace <= ps.GetMuchangPetCnt(userId) {
		msg = "牧场已满！"
		return
	}
	if ps.GetDb().Model(&models.UPet{ID: petId}).Update(map[string]interface{}{"muchang": 1}).RowsAffected > 0 {
		msg = "寄养成功！"
		result["result"] = true
		return
	} else {
		msg = "操作失败！找不到该宠物！"
		return
	}
}

func (ps *PetService) PutOut(userId, petId int, inputpwd string) (result gin.H, msg string) {
	// 携带宠物
	pet := ps.GetPet(userId, petId)
	result = gin.H{"result": false, "need_pass": false}
	if pet == nil {
		msg = "宠物不存在！"
		return
	}
	if pet.Muchang == 0 {
		msg = "宠物已携带！"
		return
	} else if pet.Muchang > 1 {
		msg = "该宠物托管中！"
		return
	}

	if 3 <= ps.GetCarryPetCnt(userId) {
		msg = "最多只能同时携带3个宠物！"
		return
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	if ps.OptSvc.UserSrv.CheckNeedPwd(inputpwd, user.McPwd) {
		result["need_pass"] = true
		msg = "请输入正确的牧场密码后再操作！"
		return
	}
	if ps.GetDb().Model(&models.UPet{ID: petId}).Update(map[string]interface{}{"muchang": 0}).RowsAffected > 0 {
		msg = "携带成功！"
		result["result"] = true
		return
	} else {
		msg = "操作失败！找不到该宠物！"
		return
	}
}

func (ps *PetService) Throw(userId, petId int, inputpwd string) (result gin.H, msg string) {
	// 丢弃牧场宠物
	pet := ps.GetPet(userId, petId)
	result = gin.H{"result": false, "need_pass": false}
	if pet == nil {
		msg = "宠物不存在！"
		return
	}
	if pet.Muchang == 0 {
		msg = "只能放生牧场中的宠物！"
		return
	} else if pet.Muchang > 1 {
		msg = "该宠物托管中！"
		return
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	if ps.OptSvc.UserSrv.CheckNeedPwd(inputpwd, user.McPwd) {
		result["need_pass"] = true
		msg = "请输入正确的牧场密码后再操作！"
		return
	}
	if !ps.OptSvc.UserSrv.DecreaseJb(userId, 10000) {
		msg = "金币不足！放生宠物需要缴纳10000金币的费用！"
		return
	}
	if ps.GetDb().Delete(pet).RowsAffected > 0 {
		msg = "放生成功！"
		result["result"] = true
		return
	} else {
		ps.OptSvc.UserSrv.DecreaseJb(userId, -10000)
		msg = "操作失败！找不到该宠物！"
		return
	}
}

func (ps *PetService) DropPet(pet *models.UPet) bool {
	// 删除宠物全部信息，包括宠物、技能，佩戴装备
	if ps.GetDb().Delete(pet).RowsAffected > 0 {
		ps.GetDb().Where("bid = ?", pet.ID)
		ps.OptSvc.PropSrv.DropPetZb(pet.Uid, pet.ID)
		ps.OptSvc.FightSrv.DelZbAttr(pet.ID)
		rcache.DelPetStatus(pet.ID)
		return true
	} else {
		return false
	}
}

func (ps *PetService) CreatPet(pet *models.UPet) bool {
	ps.GetDb().Create(pet)
	pet.GetM()

	skillitems := strings.Split(pet.SkillList, ",")
	updateFlag := false
	for _, str := range skillitems {
		items := strings.Split(str, ":")
		if len(items) < 2 {
			continue
		}

		if ok, flag := ps.CreatSkill(pet, com.StrTo(items[0]).MustInt(), com.StrTo(items[1]).MustInt()); ok {
			if flag {
				updateFlag = true
			}
		} else {
			return false
		}

	}
	if updateFlag {
		ps.GetDb().Save(pet)
	}
	return true
}

func (ps *PetService) CreatPetById(user *models.User, petId int) (bool, *models.UPet) {
	mpet := common2.GetMpet(petId)
	if mpet == nil {
		return false, nil
	}
	czlArr := strings.Split(mpet.Czl, ",")
	var czl float64
	if len(czlArr) < 2 {
		czl, _ = strconv.ParseFloat(czlArr[0], 64)
	} else {
		czl1 := com.StrTo(czlArr[0]).MustFloat64() + 0.001
		czl2 := com.StrTo(czlArr[1]).MustFloat64() + 0.001
		czl = rand.Float64()*(czl2-czl1) + czl1
		czl = utils.Round(czl, 1)
	}
	fmt.Printf("mpet %s\n", mpet)
	newPet := &models.UPet{
		Bid:         petId,
		Uid:         user.ID,
		Name:        mpet.Name,
		UserName:    user.Nickname,
		Czl:         utils.CzlStr(czl),
		Level:       1,
		Wx:          mpet.Wx,
		Hp:          mpet.Hp,
		Mp:          mpet.Mp,
		SrcHp:       mpet.Hp,
		SrcMp:       mpet.Mp,
		Ac:          mpet.Ac,
		Mc:          mpet.Mc,
		Hits:        mpet.Hits,
		Miss:        mpet.Miss,
		Speed:       mpet.Speed,
		Stime:       ps.NowUnix(),
		NowExp:      0,
		LExp:        100,
		ImgStand:    mpet.ImgStand,
		ImgAck:      mpet.ImgAck,
		ImgDie:      mpet.ImgDie,
		ImgHead:     mpet.ImgHead,
		ImgCard:     mpet.ImgCard,
		ImgEffect:   mpet.ImgEffect,
		Kx:          mpet.Kx,
		RemakeLevel: mpet.ReMakeLevel,
		RemakeId:    mpet.ReMakeId,
		RemakePid:   mpet.ReMakePid,
		ReMakeTimes: 0,
		Muchang:     0,
		SkillList:   mpet.SkillList,
	}
	return ps.CreatPet(newPet), newPet
}

func (ps *PetService) CreatSkill(upet *models.UPet, skillId, skillLv int) (bool, bool) {
	// 返回: 是否更新宠物资料是否出错
	mskill := common2.GetMskill(skillId)
	if mskill == nil {
		return false, false
	}
	newskill := &models.Uskill{
		Bid:   upet.ID,
		Sid:   skillId,
		Name:  mskill.Name,
		Level: skillLv,
		Vary:  mskill.Vary,
		Wx:    mskill.Wx,
	}
	upFlag := false
	lv := skillLv - 1
	if mskill.Vary == "4" {
		str := strings.Split(mskill.ImgEft, ",")[0]
		items := strings.Split(str, ":")
		if len(items) > 1 && com.IsSliceContainsStr([]string{"addhits", "addmc", "addac", "addhp", "addmp"}, items[0]) {
			numStr := strings.ReplaceAll(items[1], "%", "")
			rate := com.StrTo(numStr).MustFloat64() / 100.0
			upFlag = true
			switch strings.ReplaceAll(items[0], "add", "") {
			case "hits":
				upet.Hits = int(float64(upet.Hits) * (1 + rate))
				break
			case "mc":
				upet.Mc = int(float64(upet.Mc) * (1 + rate))
				break
			case "ac":
				upet.Ac = int(float64(upet.Ac) * (1 + rate))
				break
			case "hp":
				upet.Hp = int(float64(upet.Hp) * (1 + rate))
				upet.SrcHp = upet.Hp
				break
			case "mp":
				upet.Mp = int(float64(upet.Mp) * (1 + rate))
				upet.SrcMp = upet.Mp
				break
			default:
				upFlag = false
			}
		}
		newskill.Value = ""
		newskill.Plus = ""
		newskill.Img = str
		newskill.Uhp = 0
		newskill.Ump = 0
	} else if mskill.Vary == "3" {
		newskill.Value = ""
		newskill.Plus = ""
		newskill.Img = ""
		if uhpItems := strings.Split(mskill.Uhp, ","); len(uhpItems) > lv {
			newskill.Uhp = com.StrTo(uhpItems[lv]).MustInt()
		}
		if umpItems := strings.Split(mskill.Ump, ","); len(umpItems) > lv {
			newskill.Ump = com.StrTo(umpItems[lv]).MustInt()
		}
	} else {

		if valueItems := strings.Split(mskill.AckValue, ","); len(valueItems) > lv {
			newskill.Value = valueItems[lv]
		}
		if plusItems := strings.Split(mskill.Plus, ","); len(plusItems) > lv {
			newskill.Plus = plusItems[lv]
		}
		if imgItems := strings.Split(mskill.ImgEft, ","); len(imgItems) > lv {
			newskill.Img = imgItems[lv]
		}
		newskill.Uhp = 0
		if umpItems := strings.Split(mskill.Ump, ","); len(umpItems) > lv {
			newskill.Ump = com.StrTo(umpItems[lv]).MustInt()
		}
	}
	return ps.GetDb().Create(&newskill).RowsAffected > 0, upFlag
}

func (ps *PetService) GetCarryPets(UserId int) []*models.UPet {

	carryPets := []models.UPet{}
	ps.GetDb().Where("uid = ? and muchang = ?", UserId, 0).Order("Level desc, nowexp desc").Find(&carryPets)
	cpets := []*models.UPet{}
	for _, p := range carryPets {
		p1 := p
		cpets = append(cpets, &p1)
	}
	return cpets
}

func (ps *PetService) GetAllPets(UserId int) []*models.UPet {

	pets := []*models.UPet{}
	ps.GetDb().Where("uid=?", UserId).Order("level desc, nowexp desc").Find(&pets)
	return pets
	//results, err := ps.GetDb().Table("userbb ub").Select(`ub.Id, ub.bid, ub.muchang, ub.tgflag, ub.Level, b.Name, b.wx, b.cardimg`).Joins("inner join bb b on b.Id=ub.bid").Where("ub.uid =?", UserId).Order("ub.Level desc, ub.nowexp desc").Rows()
	//if err != nil {
	//	fmt.Println("error1 : ", err)
	//	return nil
	//}
	//defer results.Close()
	//for results.Next() {
	//	up := models.UPet{}
	//
	//	var id, bid, muchang, tgflag, level, name, wx, cardimg string
	//
	//	_ = results.Scan(&id, &bid, &muchang, &tgflag, &level, &name, &wx, &cardimg)
	//	up.ID = com.StrTo(id).MustInt()
	//	up.Bid = com.StrTo(bid).MustInt()
	//	up.Muchang = com.StrTo(muchang).MustInt()
	//	up.TgFlag = com.StrTo(tgflag).MustInt()
	//	up.Level = com.StrTo(level).MustInt()
	//	up.MModel = &models.MPet{
	//		Name:    name,
	//		Wx:      com.StrTo(wx).MustInt(),
	//		ImgCard: cardimg,
	//	}
	//	_ = up.MModel.AfterFind()
	//	carryPets = append(carryPets, up)
	//}
	//return carryPets
}

func (ps *PetService) IncreaseExp2Pet(pet *models.UPet, exp int) bool {
	pet.GetM()
	NowExp := pet.NowExp + exp
	Level := pet.Level
	maxLevel := ps.GetMaxLevel(pet.MModel)

	for {
		nextExp, ok := ps.GetExpByLv(Level)
		if !ok {
			return false
		}
		if NowExp < nextExp {
			break
		}
		pet.LExp = nextExp
		if Level >= maxLevel {
			if Level == pet.Level {
				return false
			} else {
				NowExp %= nextExp
				break
			}
		} else {
			Level++
			NowExp -= nextExp
		}
	}
	pet.NowExp = NowExp
	if lvDe := Level - pet.Level; lvDe > 0 {
		// 升级加属性
		grouth := common2.GetGrowth(pet.MModel.Wx)

		czl := pet.CC
		pet.Hp += int(float64(grouth.Hp) * czl * float64(lvDe))
		pet.Mp += int(float64(grouth.Mp) * czl * float64(lvDe))
		pet.Ac += int(float64(grouth.Ac) * czl * float64(lvDe))
		pet.Mc += int(float64(grouth.Mc) * czl * float64(lvDe))
		pet.Hits += int(float64(grouth.Hits) * czl * float64(lvDe))
		pet.Miss += int(float64(grouth.Miss) * czl * float64(lvDe))
		pet.Speed += int(float64(grouth.Speed) * czl * float64(lvDe))
		pet.SrcHp = pet.Hp
		pet.SrcMp = pet.Mp
		pet.Level = Level
		ps.GetDb().Model(pet).Update(gin.H{
			"nowexp": pet.NowExp,
			"lexp":   pet.LExp,
			"level":  pet.Level,
			"hp":     pet.Hp,
			"srchp":  pet.Hp,
			"mp":     pet.Mp,
			"srcmp":  pet.Mp,
			"ac":     pet.Ac,
			"mc":     pet.Mc,
			"hits":   pet.Hits,
			"miss":   pet.Miss,
			"speed":  pet.Speed,
		})
		ps.OptSvc.FightSrv.DelZbAttr(pet.ID)
	} else {
		ps.GetDb().Model(pet).Update(gin.H{
			"nowexp": pet.NowExp,
		})
	}
	return true
}

func (ps *PetService) Str2Attribute(str string) (*map[string]int, error) {
	var attributeHM map[string]int
	err := json.Unmarshal([]byte(str), &attributeHM)
	if err != nil {
		return nil, err
	}
	return &attributeHM, err
}

func (ps *PetService) IncreaseExp2MainPet(userId, exp int) {
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	if user != nil {
		pet := ps.GetPet(user.ID, user.Mbid)
		if pet != nil {
			ok := ps.IncreaseExp2Pet(pet, exp)
			if ok {

			}
		}
	}
}

func (ps *PetService) GetPetCnt(uid int) int {
	cnt := 0
	ps.GetDb().Model(&models.UPet{}).Where("uid = ?", uid).Count(&cnt)
	return cnt
}

func (ps *PetService) GetCarryPetCnt(uid int) int {
	cnt := 0
	ps.GetDb().Model(&models.UPet{}).Where("uid = ? and muchang=0", uid).Count(&cnt)
	return cnt
}

func (ps *PetService) GetMuchangPetCnt(uid int) int {
	cnt := 0
	ps.GetDb().Model(&models.UPet{}).Where("uid = ? and muchang>0", uid).Count(&cnt)
	return cnt
}

func (ps *PetService) GetMcPetCnt(uid int) int {
	cnt := 0
	ps.GetDb().Model(&models.UPet{}).Where("uid = ? and muchang=1", uid).Count(&cnt)
	return cnt
}

func (ps *PetService) GetPetKx(kx string) []int {
	items := strings.Split(kx, ",")
	kxs := make([]int, 5)
	for i, v := range items {
		kxs[i] = com.StrTo(v).MustInt()
	}
	return kxs
}

func (ps *PetService) GetPetSkill(petId int) []*models.Uskill {
	skills := []*models.Uskill{}
	ps.GetDb().Where("bid = ?", petId).Find(&skills)
	return skills
}

func (ps *PetService) GetSkill(skillId int) *models.Uskill {
	skill := &models.Uskill{}
	ps.GetDb().Where("id = ?", skillId).First(skill)
	return skill
}

func (ps *PetService) GetSkillBySid(petId, sid int) *models.Uskill {
	skill := &models.Uskill{}
	ps.GetDb().Where("sid = ? and bid=?", sid, petId).First(skill)
	return skill
}

func (ps *PetService) GetMskillByPid(propId int) *models.MSkill {
	return common2.GetMskillByPid(propId)
}

func (ps *PetService) AddPetAttribute(petId int, attrName string, attrValue int) bool {
	// 只能加属性，不能加成长和等级经验
	return ps.GetDb().Model(&models.UPet{ID: petId}).Update(UpMap{attrName: gorm.Expr("? + ?", attrName, attrValue)}).RowsAffected > 0
}

func (ps *PetService) SetPetCzl(petId int, czl string) bool {
	// 设置成长，成长为字符串类型
	return ps.GetDb().Model(&models.UPet{ID: petId}).Update(UpMap{"Czl": czl}).RowsAffected > 0
}

func (ps *PetService) LearnSkill(userId, skillPropId int) (bool, string) {
	prop := ps.OptSvc.PropSrv.GetProp(userId, skillPropId, false)
	if prop == nil || prop.Sums == 0 {
		return false, "您没有该技能书！"
	}
	aimSkill := &models.MSkill{}
	ps.GetDb().Where("pid=?", prop.Pid).Find(aimSkill)
	if aimSkill.ID == 0 {
		return false, "没有技能书所对应的技能！"
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	petSkills := ps.GetPetSkill(user.Mbid)
	for _, skill := range petSkills {
		if skill.Sid == aimSkill.ID {
			return false, "宠物已拥有该技能！"
		}
	}
	mainPet := ps.GetPetById(user.Mbid)
	ps.CreatSkill(mainPet, aimSkill.ID, 1)
	ps.GetDb().Where("id=?", mainPet).Update(gin.H{
		"skillist": mainPet.SkillList + fmt.Sprintf("%d:1", aimSkill.ID),
		"hp":       mainPet.Hp,
		"srchp":    mainPet.SrcHp,
		"mp":       mainPet.Mp,
		"srcmp":    mainPet.SrcMp,
		"ac":       mainPet.Ac,
		"mc":       mainPet.Mc,
		"hits":     mainPet.Hits,
		"miss":     mainPet.Miss,
		"speed":    mainPet.Speed,
	})
	ps.OptSvc.FightSrv.DelZbAttr(mainPet.ID)
	return true, "学习技能成功！"
}

func (ps *PetService) UpdateSkill(userId, skillId int) (bool, string) {
	skill := &models.Uskill{ID: skillId}
	ps.GetDb().Find(skill)
	if skill.Sid == 0 {
		return false, "该技能不存在！"
	}
	pet := ps.GetPetById(skill.Bid)
	if pet == nil || pet.Uid != userId {
		return false, "该技能不存在！"
	}
	if skill.Level == 10 {
		return false, "该技能不能再升级！"
	}
	skill.GetM()
	skillRequires := strings.Split(skill.MModel.Requires, ",")
	if len(skillRequires) > skill.Level {
		if pet.Level < com.StrTo(skillRequires[skill.Level]).MustInt() {
			return false, fmt.Sprintf("宠物需要等级：%s！", skillRequires[skill.Level])
		}
	}

	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	upFlag := false
	lv := skill.Level
	if skill.MModel.Vary == "4" {
		strs := strings.Split(skill.MModel.ImgEft, ",")
		var str string
		if len(strs) <= lv {
			return false, "该技能不能再升级！"
		}
		str = strs[lv]
		items := strings.Split(str, ":")
		if len(items) > 1 && com.IsSliceContainsStr([]string{"addhits", "addmc", "addac", "addhp", "addmp"}, items[0]) {
			numStr := strings.ReplaceAll(items[1], "%", "")
			rate := com.StrTo(numStr).MustFloat64() / 100.0
			upFlag = true
			switch strings.ReplaceAll(items[0], "add", "") {
			case "hits":
				pet.Hits = int(float64(pet.Hits) * (1 + rate))
				ps.GetDb().Model(pet).Update(gin.H{"hits": pet.Hits})
				break
			case "mc":
				pet.Mc = int(float64(pet.Mc) * (1 + rate))
				ps.GetDb().Model(pet).Update(gin.H{"mc": pet.Mc})
				break
			case "ac":
				pet.Ac = int(float64(pet.Ac) * (1 + rate))
				ps.GetDb().Model(pet).Update(gin.H{"ac": pet.Ac})
				break
			case "hp":
				pet.Hp = int(float64(pet.Hp) * (1 + rate))
				pet.SrcHp = pet.Hp
				ps.GetDb().Model(pet).Update(gin.H{"hp": pet.Hp, "srchp": pet.SrcHp})
				break
			case "mp":
				pet.Mp = int(float64(pet.Mp) * (1 + rate))
				pet.SrcMp = pet.Mp
				ps.GetDb().Model(pet).Update(gin.H{"mp": pet.Mp, "srcmp": pet.SrcMp})
				break
			default:
				upFlag = false
			}
		}
		if upFlag {
			ps.OptSvc.FightSrv.DelZbAttr(pet.ID)
		}
		if !ps.OptSvc.PropSrv.DecrPropByPid(userId, 1666, 1) {
			return false, "您的背包没有buff技能升级卷轴！"
		}
		ps.GetDb().Model(skill).Update(gin.H{"img": str, "level": skill.Level + 1})
	} else if skill.MModel.Vary == "3" {
		newuhp := skill.Uhp
		newump := skill.Ump
		if uhpItems := strings.Split(skill.MModel.Uhp, ","); len(uhpItems) > lv {
			newuhp = com.StrTo(uhpItems[lv]).MustInt()
		}
		if umpItems := strings.Split(skill.MModel.Ump, ","); len(umpItems) > lv {
			newump = com.StrTo(umpItems[lv]).MustInt()
		}
		if !ps.OptSvc.PropSrv.DecrPropByPid(userId, 733, 1) {
			return false, "您的背包没有技能升级卷轴！"
		}
		ps.GetDb().Model(skill).Update(gin.H{"uhp": newuhp, "ump": newump, "level": skill.Level + 1})
	} else {

		newValue := skill.Value
		if valueItems := strings.Split(skill.MModel.AckValue, ","); len(valueItems) > lv {
			newValue = valueItems[lv]
		}
		newPlus := skill.Plus
		if plusItems := strings.Split(skill.MModel.Plus, ","); len(plusItems) > lv {
			newPlus = plusItems[lv]
		}
		newImg := skill.Img
		if imgItems := strings.Split(skill.MModel.ImgEft, ","); len(imgItems) > lv {
			newImg = imgItems[lv]
		}
		newUmp := skill.Ump
		if umpItems := strings.Split(skill.MModel.Ump, ","); len(umpItems) > lv {
			newUmp = com.StrTo(umpItems[lv]).MustInt()
		}
		if !ps.OptSvc.PropSrv.DecrPropByPid(userId, 733, 1) {
			return false, "您的背包没有技能升级卷轴！"
		}
		ps.GetDb().Model(skill).Update(gin.H{"ump": newUmp, "value": newValue, "plus": newPlus, "img": newImg, "level": skill.Level + 1})
	}
	ps.OptSvc.Commit()
	return true, "技能升级成功！"
}

func (ps *PetService) Evolution(userId, petId int, pathA bool, fzid int) (bool, string) {
	//rcache.NewRdbHandler("Evolution")
	now := utils.NowUnix()
	if rcache.EvolutionTimer.InCoolTime(userId, now) {
		return false, "进化正在冷却中"
	}
	pet := ps.GetPet(userId, petId)
	if pet == nil || pet.Muchang > 0 {
		return false, "宠物不存在！"
	}
	pet.GetM()
	if pet.MModel.Wx >= 7 {
		return false, "神圣宠物不可在此进化！"
	}
	if pet.ReMakeTimes >= 10 {
		return false, "宠物已进化满10次了！"
	}
	propIds := strings.Split(pet.MModel.ReMakePid, ",")
	petIds := strings.Split(pet.MModel.ReMakeId, ",")
	levels := strings.Split(pet.MModel.ReMakeLevel, ",")
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	var topet *models.MPet
	var needLevel int
	var useProp *models.MProp
	var usePropIds []string
	if pathA {
		apropIds := strings.Split(propIds[0], "|")
		usePropIds = apropIds
		apet := common2.GetMpet(com.StrTo(petIds[0]).MustInt())
		if apet == nil {
			return false, "该宠物无法进化！"
		}
		needLevel = com.StrTo(levels[0]).MustInt()
		if pet.Level < needLevel {
			return false, "该宠物等级未达到进化要求！"
		}
		topet = apet
	} else {
		if len(propIds) > 1 {
			bpropIds := strings.Split(propIds[1], "|")
			usePropIds = bpropIds
		} else {
			return false, "该宠物无法进化！"
		}
		if len(petIds) > 1 {
			bpet := common2.GetMpet(com.StrTo(petIds[1]).MustInt())
			if bpet == nil {
				return false, "该宠物无法进化！"
			}
			topet = bpet
		} else {
			return false, "该宠物无法进化！"
		}
		if len(levels) > 1 {
			needLevel = com.StrTo(levels[1]).MustInt()
		} else {
			return false, "该宠物无法进化！"
		}
	}
	finded := false
	for i, _ := range usePropIds {
		prop := common2.GetMProp(com.StrTo(usePropIds[i]).MustInt())
		if prop != nil && ps.OptSvc.PropSrv.DecrPropByPid(userId, prop.ID, 1) {
			finded = true
			useProp = prop
			break
		}
	}
	if !finded {
		return false, "进化道具不足！"
	}
	if !ps.OptSvc.UserSrv.DecreaseJb(userId, 1000) {
		return false, "金币不足！"
	}
	newCzl := pet.CC
	fmt.Printf("进化方向为A:%s", pathA)
	if pet.MModel.Wx < 6 {
		if pathA {
			if newCzl < 50 {
				newCzl += float64(rand.Intn(5)+1)/10 + float64(pet.Level-needLevel)/200
			} else if newCzl >= 50 && newCzl < 80 {
				newCzl += float64(rand.Intn(3)+1) / 10
			} else {
				newCzl += 0.1
			}
		} else {
			if newCzl < 50 {
				newCzl += float64(rand.Intn(6)+5)/10 + float64(pet.Level-needLevel)/200
			} else if newCzl >= 50 && newCzl < 70 {
				newCzl += float64(rand.Intn(4)+4) / 10
			} else if newCzl >= 70 && newCzl < 80 {
				newCzl += float64(rand.Intn(3)+3) / 10
			} else if newCzl >= 80 && newCzl < 90 {
				newCzl += float64(rand.Intn(2)+2) / 10
			} else {
				newCzl += float64(rand.Intn(3)+1) / 10
			}
		}
		if newCzl > 150 {
			if fzid != 0 {
				fzprop := ps.OptSvc.PropSrv.GetProp(userId, fzid, false)
				if fzprop == nil || !ps.OptSvc.PropSrv.DecrPropById(fzid, 1) {
					return false, "所添加进化成长保护石数量不足！"
				}
				fzprop.GetM()
				pczl := com.StrTo(strings.ReplaceAll(fzprop.MModel.Effect, "keepczl:", "")).MustInt()
				if pczl >= 150 {
					newCzl = float64(pczl)
				}
			} else {
				newCzl = 150
			}
		}

	} else {
		if pathA {
			newCzl += float64(rand.Intn(3)+1) / 10
		} else {
			newCzl += float64(rand.Intn(4)+3) / 10
		}
	}

	ps.GetDb().Model(pet).Update(gin.H{
		"bid":         topet.ID,
		"name":        topet.Name,
		"imgstand":    topet.ImgStand,
		"imgack":      topet.ImgAck,
		"imgdie":      topet.ImgDie,
		"cardimg":     topet.ImgCard,
		"effectimg":   topet.ImgEffect,
		"remaketimes": gorm.Expr("remaketimes+1"),
		"czl":         utils.CzlStr(newCzl)})
	ps.OptSvc.Commit()
	pathName := "A"
	if !pathA {
		pathName = "B"
	}
	SelfGameLog(userId, fmt.Sprintf("宠物进化：%s %s进化得到 %s，进化道具：%s", pet.MModel.Name, pathName, topet.Name, useProp.Name), 99)
	rcache.EvolutionTimer.Set(userId, now)
	return true, "进化成功！"
}

func (ps *PetService) Emerge(userId, apetId, bpetId, apropId, bpropId int, noCheckZb, noCheckProtect bool) (result gin.H, msg string) {
	// 注意宠物要本人的，道具也要本人的
	result = gin.H{"result": false, "have_zb": false, "no_protect": false, "in_cool": true}
	now := utils.NowUnix()
	if rcache.MergeTimer.InCoolTime(userId, now) {
		msg = "合成冷却中！"
		return
	}
	apet := ps.GetPet(userId, apetId)
	if apet == nil {
		return result, "主宠不存在！"
	}
	bpet := ps.GetPet(userId, bpetId)
	if bpet == nil {
		return result, "副宠不存在！"
	}
	if apet.Muchang != 0 || bpet.Muchang != 0 {
		return result, "牧场中的宠物无法参与合成！"
	}
	if apet.CqFlag == 1 || bpet.CqFlag == 1 {
		return result, "抽取过成长的宠物无法参与合成！"
	}
	rule := &models.MergeRule{}
	ps.GetDb().Where("aid=? and bid=?", apet.Bid, bpet.Bid).First(rule)
	if rule.Id == 0 {
		return result, "两只宠物无法合成！"
	}
	if apet.Level < 40 || bpet.Level < 40 {
		return result, "等级不足无法合成！"
	}
	var maxCzl float64 = 0
	if rule.Limits != "" && rule.Limits != "0" {
		ruleItems := strings.Split(rule.Limits, "|")
		if len(ruleItems) > 1 {
			if apet.CC < com.StrTo(ruleItems[0]).MustFloat64() {
				return result, "主宠成长未达到要求！"
			}
			if bpet.CC < com.StrTo(ruleItems[1]).MustFloat64() {
				return result, "副宠成长未达到要求！"
			}
			if len(ruleItems) > 2 {
				maxCzl = com.StrTo(ruleItems[2]).MustFloat64()
			}
		}
	}
	failedProtect := false
	addAttr := make(map[string]int)
	aSuccesRate := 0
	bSuccesRate := 0
	propNotes := []string{}
	if apropId != 0 {
		aprop := ps.OptSvc.PropSrv.GetProp(userId, apropId, false)
		if aprop == nil || aprop.Sums == 0 {
			return result, "添加物一数量不足！"
		}
		aprop.GetM()
		propNotes = append(propNotes, "使用物品1："+aprop.MModel.Name)
		effectItems := strings.Split(aprop.MModel.Effect, "|")
		for _, effcetItem := range effectItems {
			items := strings.SplitN(effcetItem, ":", 2)
			if len(items) > 1 {
				switch items[0] {
				case "hecheng":
					successRateItems := strings.Split(items[1], ",")
					for _, s := range successRateItems {
						sitems := strings.Split(s, ":")
						if sitems[0] == "A" {
							aSuccesRate += com.StrTo(strings.ReplaceAll(sitems[1], "%", "")).MustInt()
						} else if sitems[0] == "B" {
							bSuccesRate += com.StrTo(strings.ReplaceAll(sitems[1], "%", "")).MustInt()
						}
					}
					break
				default:
					if strings.Index(items[0], "add") > -1 {
						addAttr[strings.ReplaceAll(items[0], "add", "")] += com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustInt()
					}
					break
				}
			} else {
				// 是否保底
				if items[0] == "1" {
					failedProtect = true
				}
			}
		}
	}
	if bpropId != 0 {
		bprop := ps.OptSvc.PropSrv.GetProp(userId, bpropId, false)
		if bprop == nil || bprop.Sums == 0 {
			return result, "添加物二数量不足！"
		}
		bprop.GetM()
		propNotes = append(propNotes, "使用物品2："+bprop.MModel.Name)
		effectItems := strings.Split(bprop.MModel.Effect, "|")
		for _, effcetItem := range effectItems {
			items := strings.SplitN(effcetItem, ":", 2)
			if len(items) > 1 {
				switch items[0] {
				case "hecheng":
					successRateItems := strings.Split(items[1], ",")
					for _, s := range successRateItems {
						sitems := strings.Split(s, ":")
						if sitems[0] == "A" {
							aSuccesRate += com.StrTo(strings.ReplaceAll(sitems[1], "%", "")).MustInt()
						} else if sitems[0] == "B" {
							bSuccesRate += com.StrTo(strings.ReplaceAll(sitems[1], "%", "")).MustInt()
						}
					}
					break
				default:
					if strings.Index(items[0], "add") > -1 {
						addAttr[strings.ReplaceAll(items[0], "add", "")] += com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustInt()
					}
					break
				}
			} else {
				// 是否保底
				if items[0] == "1" {
					failedProtect = true
				}
			}
		}
	}
	if !failedProtect && !noCheckProtect {
		result["result"] = true
		result["no_protect"] = true
		return result, "您没有添加失败保护道具，是否继续合成？"
	}
	if (apet.Zb != "" || bpet.Zb != "") && !noCheckZb {
		result["result"] = true
		result["have_zb"] = true
		return result, "宠物身上带有装备，是否继续合成？"
	}

	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	if !ps.OptSvc.UserSrv.DecreaseJb(userId, 50000) {
		return result, "金币不足！"
	}
	if apropId != 0 && !ps.OptSvc.PropSrv.DecrPropById(apropId, 1) {
		return result, "添加物一数量不足！"
	}
	if bpropId != 0 && !ps.OptSvc.PropSrv.DecrPropById(bpropId, 1) {
		return result, "添加物二数量不足！"
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ps.OptSvc.UserSrv.GetUserInfoById(userId)
	apet.GetM()
	bpet.GetM()

	// 计算合成成功率
	successLuckeyNum := float64(rand.Intn(101)+1) / 100
	var successRate float64
	if userInfo.HechengNums == 10 || apet.CC < 5 {
		successRate = 1
	} else {
		successRate = float64(userInfo.HechengNums)/(apet.CC*2) + float64((apet.Level+bpet.Level)/15)*0.01 + float64(aSuccesRate)*0.01 + float64(rand.Intn(5)+1)*0.01
	}

	if successLuckeyNum <= successRate {
		// 成功
		toPetId := rule.Maid
		pathLuckeyNum := float64(rand.Intn(101)+1) / 100
		if pathLuckeyNum <= 0.05+float64(bSuccesRate)*0.01 {
			toPetId = rule.Mbid
		}

		ps.GetDb().Update(userInfo).Update(gin.H{"hecheng_nums": 0})
		if !ps.DropPet(apet) {
			return result, "合成失败！找不到主宠"
		}

		if !ps.DropPet(bpet) {
			return result, "合成失败！找不到副宠"
		}
		toPet := common2.GetMpet(toPetId)
		newPet := &models.UPet{
			Bid:       toPet.ID,
			Uid:       userId,
			Level:     1,
			Ac:        toPet.Ac + int(float64(apet.Ac*apet.Level)/400+float64(bpet.Ac*bpet.Level)/800),
			Mc:        toPet.Mc + int(float64(apet.Mc*apet.Level)/400+float64(bpet.Mc*bpet.Level)/800),
			SrcHp:     toPet.Hp + int(float64(apet.SrcHp*apet.Level)/400+float64(bpet.SrcHp*bpet.Level)/800),
			Hp:        toPet.Hp + int(float64(apet.Hp*apet.Level)/400+float64(bpet.Hp*bpet.Level)/800),
			Mp:        toPet.Mp + int(float64(apet.Mp*apet.Level)/400+float64(bpet.Mp*bpet.Level)/800),
			SrcMp:     toPet.Mp + int(float64(apet.SrcMp*apet.Level)/400+float64(bpet.SrcMp*bpet.Level)/800),
			Stime:     utils.NowUnix(),
			NowExp:    0,
			LExp:      100,
			Hits:      toPet.Hits + int(float64(apet.Hits*apet.Level)/400+float64(bpet.Hits*bpet.Level)/800),
			Miss:      toPet.Miss + int(float64(apet.Miss*apet.Level)/400+float64(bpet.Miss*bpet.Level)/800),
			Speed:     toPet.Speed + int(float64(apet.Speed*apet.Level)/400+float64(bpet.Speed*bpet.Level)/800),
			Kx:        "",
			Fatting:   0,
			Czl:       "",
			Muchang:   0,
			AddSx:     "",
			SkillList: toPet.SkillList,
		}
		addCzl := 0.0
		if apet.CC < 51 {
			addCzl = float64(apet.Level)/float64(apet.CC+10) + float64(bpet.Level)*bpet.CC/200
		} else if apet.CC < 70 {
			addCzl = float64(apet.Level)/apet.CC + float64(bpet.Level)*bpet.CC/350
		} else if apet.CC < 90 {
			addCzl = float64(apet.Level)/apet.CC + float64(bpet.Level)*bpet.CC/500
		} else if apet.CC < 100 {
			addCzl = float64(apet.Level)/apet.CC + float64(bpet.Level)*bpet.CC/700
		} else {
			addCzl = float64(apet.Level)/apet.CC + float64(bpet.Level)*bpet.CC/900
		}
		for k, v := range addAttr {
			if k == "czl" {
				addCzl += addCzl * float64(v) * 0.01
			} else {
				switch k {
				case "hp":
					newPet.Hp += int(float64(newPet.Hp) * float64(v) * 0.01)
					break
				case "mp":
					newPet.Mp += int(float64(newPet.Mp) * float64(v) * 0.01)
					break
				case "ac":
					newPet.Ac += int(float64(newPet.Ac) * float64(v) * 0.01)
					break
				case "mc":
					newPet.Mc += int(float64(newPet.Mc) * float64(v) * 0.01)
					break
				case "hits":
					newPet.Hits += int(float64(newPet.Hits) * float64(v) * 0.01)
					break
				case "miss":
					newPet.Miss += int(float64(newPet.Miss) * float64(v) * 0.01)
					break
				case "speed":
					newPet.Speed += int(float64(newPet.Speed) * float64(v) * 0.01)
					break
				}
			}
		}
		if maxCzl == 0 {
			if toPet.Wx == 6 {
				maxCzl = 60
			}
		}
		if maxCzl != 0 && apet.CC+addCzl > maxCzl {
			newPet.Czl = utils.CzlStr(maxCzl)
		} else {
			newPet.Czl = utils.CzlStr(apet.CC + addCzl)
		}

		if ok := ps.CreatPet(newPet); !ok {
			return result, "合成失败！宠物数据出错！"
		}
		ps.OptSvc.UserSrv.SetMBid(userId, newPet.ID)
		ps.GetDb().Model(userInfo).Update(gin.H{"hecheng_nums": 0})
		ps.OptSvc.Commit()
		SelfGameLog(userId,
			fmt.Sprintf("合成结果：成功！新宠物，名字：%s  czl:%s  ac:%d  hits:%d,\n%s,\n主宠物：%s level:%d  czl:%s  ac:%d  hits:%d\n副宠物：%s  level:%d  czl:%s  ac:%d  hits:%d",
				toPet.Name, newPet.Czl, newPet.Ac, newPet.Hits, strings.Join(propNotes, ",\n"), apet.MModel.Name, apet.Level, apet.Czl, apet.Ac, apet.Hits, bpet.MModel.Name, bpet.Level, bpet.Czl, bpet.Ac, bpet.Hits), 4)
		AnnounceAll(user.Nickname, fmt.Sprintf("成功的合成了一只%s,真是太幸运了!", toPet.Name))
		result["result"] = true
		rcache.MergeTimer.Set(userId, now)
		return result, "成功合成：" + toPet.Name
	} else {
		ps.GetDb().Model(userInfo).Update(gin.H{"hecheng_nums": gorm.Expr("hecheng_nums+1")})
		rcache.MergeTimer.Set(userId, now)
		result["result"] = true
		if !failedProtect {
			ps.DropPet(bpet)
			if user.Mbid == bpet.ID {
				ps.OptSvc.UserSrv.SetMBid(userId, apetId)
			}
			SelfGameLog(userId, fmt.Sprintf("合成结果：失败！\n%s,\n消失副宠物：%s  level:%d  czl:%s  ac:%d  hits:%d",
				strings.Join(propNotes, ",\n"), bpet.MModel.Name, bpet.Level, bpet.Czl, bpet.Ac, bpet.Hits), 4)
			ps.OptSvc.Commit()
			return result, "合成失败！副宠消失！"
		} else {
			SelfGameLog(userId, fmt.Sprintf("合成结果：失败！\n%s", strings.Join(propNotes, ",\n")), 4)
			ps.OptSvc.Commit()
			return result, "合成失败！"
		}

	}
}

func (ps *PetService) PetOffzb(userId, zbId int) (bool, string) {
	prop := ps.OptSvc.PropSrv.GetProp(userId, zbId, false)
	if prop == nil {
		return false, "道具不存在！"
	}
	if prop.Zbpets == 0 {
		return false, "道具并未装备到宠物！"
	}
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	if pet := ps.OptSvc.PetSrv.GetPet(userId, prop.Zbpets); pet != nil {
		newzbs := []string{}
		for _, v := range strings.Split(pet.Zb, ",") {
			if items := strings.Split(v, ":"); len(items) > 1 {
				if items[1] != strconv.Itoa(prop.ID) {
					newzbs = append(newzbs, v)
				}
			}
		}
		if ps.GetDb().Model(prop).Update(UpMap{"zbing": 0, "zbpets": 0}).RowsAffected > 0 {
			if ps.GetDb().Model(pet).Update("zb", strings.Join(newzbs, ",")).RowsAffected > 0 {
				ps.OptSvc.Commit()
				ps.OptSvc.FightSrv.DelZbAttr(pet.ID)
				return true, "卸载装备成功！"
			}
		}
	}

	return false, "操作失败！"
}

func (ps *PetService) ZhuanSheng(userId, apetId, bpetId, cpetId, apropId, bpropId int) (bool, string) {
	now := utils.NowUnix()
	if rcache.ZhuangshengTimer.InCoolTime(userId, now) {
		return false, "转生冷却中！"
	}
	apet := ps.GetPet(userId, apetId)
	if apet == nil {
		return false, "主宠不存在！"
	}
	if apet.Muchang > 0 {
		return false, "无法将牧场中的宠物进行涅槃！"
	}
	if apet.Level < 60 {
		return false, "主宠等级不足！"
	}

	bpet := ps.GetPet(userId, bpetId)
	if bpet == nil {
		return false, "副宠不存在！"
	}
	if bpet.Muchang > 0 {
		return false, "无法将牧场中的宠物进行涅槃！"
	}
	if bpet.Level < 60 {
		return false, "副宠等级不足！"
	}

	cpet := ps.GetPet(userId, cpetId)
	if cpet == nil {
		return false, "涅槃兽不存在！"
	}
	if cpet.Muchang > 0 {
		return false, "无法将牧场中的宠物进行涅槃！"
	}
	if cpet.Level < 60 {
		return false, "涅槃兽等级不足！"
	}
	apet.GetM()
	bpet.GetM()
	cpet.GetM()
	if strings.Index(cpet.MModel.Name, "涅") < 0 {
		return false, "请添加正确的涅槃兽！"
	}
	rule := &models.ZsRule{}
	ps.GetDb().Where("aid=? and bid=?", apet.MModel.ID, bpet.MModel.ID).First(rule)
	if rule.Id == 0 {
		return false, "两只神宠无法进行涅槃！"
	}

	failedProtect := false
	var addSuccessRate float64 = 0
	var addCzlRate float64 = 0
	var propNotes []string
	if apropId != 0 {
		aprop := ps.OptSvc.PropSrv.GetProp(userId, apropId, false)
		if aprop == nil || aprop.Sums == 0 {
			return false, "添加道具一不存在！"
		}
		aprop.GetM()
		propNotes = append(propNotes, "添加道具一："+aprop.MModel.Name)
		for _, v := range strings.Split(aprop.MModel.Effect, ",") {
			items := strings.Split(v, ":")
			if len(items) > 1 {
				if items[0] == "npbb" {
					failedProtect = true
				} else if items[0] == "npcg" {
					addSuccessRate += com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustFloat64() * 0.01
				} else if items[0] == "npcz" {
					addCzlRate += com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustFloat64() * 0.01
				}
			}
		}
	}
	if bpropId != 0 {
		bprop := ps.OptSvc.PropSrv.GetProp(userId, bpropId, false)
		if bprop == nil || bprop.Sums == 0 {
			return false, "添加道具二不存在！"
		}
		bprop.GetM()
		propNotes = append(propNotes, "添加道具二："+bprop.MModel.Name)
		for _, v := range strings.Split(bprop.MModel.Effect, ",") {
			items := strings.Split(v, ":")
			if len(items) > 1 {
				if items[0] == "npbb" {
					failedProtect = true
				} else if items[0] == "npcg" {
					addSuccessRate += com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustFloat64() * 0.01
				} else if items[0] == "npcz" {
					addCzlRate += com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustFloat64() * 0.01
				}
			}
		}
	}

	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()

	if !ps.OptSvc.UserSrv.DecreaseJb(userId, 500000) {
		return false, "金币不足！"
	}
	if apropId != 0 && !ps.OptSvc.PropSrv.DecrPropById(apropId, 1) {
		return false, "添加物一数量不足！"
	}
	if bpropId != 0 && !ps.OptSvc.PropSrv.DecrPropById(bpropId, 1) {
		return false, "添加物二数量不足！"
	}
	luckeyNum := rand.Intn(10000) + 1
	if luckeyNum <= int(float64(apet.Level+bpet.Level)/30*(1+addSuccessRate)*100) {
		// 成功
		toPetId := rule.Mid
		if cpet.Name == "涅磐兽（亥）" {
			addCzlRate += 0.05
		} else if cpet.Name == "涅磐兽（午）" {
			addCzlRate += 0.15
		} else if cpet.Name == "涅磐兽（卯）" {
			addCzlRate += 0.3
		}
		toPet := common2.GetMpet(toPetId)
		newPet := &models.UPet{
			Bid:       toPet.ID,
			Uid:       userId,
			Level:     1,
			Ac:        toPet.Ac + int(float64(apet.Ac*apet.Level)/400+float64(bpet.Ac*bpet.Level)/800),
			Mc:        toPet.Mc + int(float64(apet.Mc*apet.Level)/400+float64(bpet.Mc*bpet.Level)/800),
			SrcHp:     toPet.Hp + int(float64(apet.SrcHp*apet.Level)/400+float64(bpet.SrcHp*bpet.Level)/800),
			Hp:        toPet.Hp + int(float64(apet.Hp*apet.Level)/400+float64(bpet.Hp*bpet.Level)/800),
			Mp:        toPet.Mp + int(float64(apet.Mp*apet.Level)/400+float64(bpet.Mp*bpet.Level)/800),
			SrcMp:     toPet.Mp + int(float64(apet.SrcMp*apet.Level)/400+float64(bpet.SrcMp*bpet.Level)/800),
			Stime:     utils.NowUnix(),
			NowExp:    0,
			LExp:      100,
			Hits:      toPet.Hits + int(float64(apet.Hits*apet.Level)/400+float64(bpet.Hits*bpet.Level)/800),
			Miss:      toPet.Miss + int(float64(apet.Miss*apet.Level)/400+float64(bpet.Miss*bpet.Level)/800),
			Speed:     toPet.Speed + int(float64(apet.Speed*apet.Level)/400+float64(bpet.Speed*bpet.Level)/800),
			Kx:        "",
			Fatting:   0,
			Czl:       "",
			Muchang:   0,
			AddSx:     "",
			SkillList: toPet.SkillList,
		}

		addCzlPetSettings1 := common2.GetWelcome("zs1")
		addCzlPetSettings2 := common2.GetWelcome("zs2")
		addCzlPetSettings3 := common2.GetWelcome("zs3")
		var num1, num2 float64 = 0, 0
		if addCzlPetSettings1 != nil && com.IsSliceContainsStr(strings.Split(addCzlPetSettings1.Content, ","), apet.MModel.Name) {
			// 没有加成的
			if apet.CC <= 10.9 {
				num1, num2 = 1, 2200
			} else if apet.CC <= 30.9 {
				num1, num2 = 1, 250
			} else if apet.CC <= 49.9 {
				num1, num2 = 1, 350
			} else if apet.CC <= 60.9 {
				num1, num2 = 1, 480
			} else if apet.CC <= 70.9 {
				num1, num2 = 1, 600
			} else if apet.CC <= 80.9 {
				num1, num2 = 1, 800
			} else if apet.CC <= 90.9 {
				num1, num2 = 2, 1200
			} else {
				num1, num2 = 2, 2200
			}
		} else if addCzlPetSettings2 != nil && com.IsSliceContainsStr(strings.Split(addCzlPetSettings2.Content, ","), apet.MModel.Name) {
			// 有小加成的
			if apet.CC <= 10.9 {
				num1, num2 = 1, 2190
			} else if apet.CC <= 30.9 {
				num1, num2 = 1, 240
			} else if apet.CC <= 49.9 {
				num1, num2 = 1, 340
			} else if apet.CC <= 60.9 {
				num1, num2 = 1, 470
			} else if apet.CC <= 70.9 {
				num1, num2 = 1, 590
			} else if apet.CC <= 80.9 {
				num1, num2 = 1, 780
			} else if apet.CC <= 90.9 {
				num1, num2 = 2, 1100
			} else {
				num1, num2 = 2, 1800
			}
		} else if addCzlPetSettings3 != nil && com.IsSliceContainsStr(strings.Split(addCzlPetSettings3.Content, ","), apet.MModel.Name) {
			// 有大加成的
			fmt.Printf("%s\n", addCzlPetSettings3.Content)
			if apet.CC <= 10.9 {
				num1, num2 = 1, 2180
			} else if apet.CC <= 30.9 {
				num1, num2 = 1, 230
			} else if apet.CC <= 49.9 {
				num1, num2 = 1, 330
			} else if apet.CC <= 60.9 {
				num1, num2 = 1, 450
			} else if apet.CC <= 70.9 {
				num1, num2 = 1, 570
			} else if apet.CC <= 80.9 {
				num1, num2 = 1, 760
			} else if apet.CC <= 90.9 {
				num1, num2 = 2, 1000
			} else {
				num1, num2 = 2, 1500
			}
		} else {
			return false, "涅槃失败！找不到主宠加成数据！"
		}
		addCzl := (float64(apet.Level)/apet.CC/num1 + float64(bpet.Level)*bpet.CC/num2) * (1 + addCzlRate)
		newPet.Czl = utils.CzlStr(apet.CC + addCzl)

		if !ps.DropPet(apet) {
			return false, "涅槃失败！找不到主宠"
		}

		if !ps.DropPet(bpet) {
			return false, "涅槃失败！找不到副宠"
		}

		if !ps.DropPet(cpet) {
			return false, "涅槃失败！找不到涅槃兽！"
		}

		if ok := ps.CreatPet(newPet); !ok {
			return false, "涅槃失败！宠物数据出错！"
		}
		ps.OptSvc.UserSrv.SetMBid(userId, newPet.ID)
		ps.OptSvc.Commit()
		SelfGameLog(userId,
			fmt.Sprintf("涅槃结果：成功！新宠物，名字：%s  czl:%s  ac:%d  hits:%d,\n%s,\n主宠物：%s level:%d  czl:%s  ac:%d  hits:%d\n副宠物：%s  level:%d  czl:%s  ac:%d  hits:%d\n涅槃兽：%s  level:%d  czl:%s  ac:%d  hits:%d",
				toPet.Name, newPet.Czl, newPet.Ac, newPet.Hits, strings.Join(propNotes, ",\n"), apet.MModel.Name, apet.Level, apet.Czl, apet.Ac, apet.Hits, bpet.MModel.Name, bpet.Level, bpet.Czl, bpet.Ac, bpet.Hits, cpet.MModel.Name, cpet.Level, cpet.Czl, cpet.Ac, cpet.Hits), 4)
		user := ps.OptSvc.UserSrv.GetUserById(userId)
		AnnounceAll(user.Nickname, fmt.Sprintf("成功的转生获得了一只%s,真是太幸运了!", toPet.Name))
		rcache.ZhuangshengTimer.Set(userId, now)
		return true, "成功的转生：" + toPet.Name
	} else {
		// 失败
		if !failedProtect {
			ps.DropPet(cpet)

		}
		ps.OptSvc.Commit()
		SelfGameLog(userId,
			fmt.Sprintf("涅槃结果：失败！\n主宠物：%s level:%d  czl:%s  ac:%d  hits:%d\n副宠物：%s  level:%d  czl:%s  ac:%d  hits:%d\n涅槃兽：%s  level:%d  czl:%s  ac:%d  hits:%d",
				strings.Join(propNotes, ",\n"), apet.MModel.Name, apet.Level, apet.Czl, apet.Ac, apet.Hits, bpet.MModel.Name, bpet.Level, bpet.Czl, bpet.Ac, bpet.Hits, cpet.MModel.Name, cpet.Level, cpet.Czl, cpet.Ac, cpet.Hits), 4)
		rcache.ZhuangshengTimer.Set(userId, now)
		return true, "转生失败！"
	}
}

func (ps *PetService) Chouqu(userId, petId, apid, bpid int) (bool, string) {
	pet := ps.GetPet(userId, petId)
	if pet == nil {
		return false, "宠物不存在！"
	}
	if pet.CqFlag == 1 {
		return false, "宠物已抽取过成长！"
	}
	pet.GetM()
	if pet.MModel.Wx > 6 {
		return false, "神圣宠物无法抽取成长！"
	}
	if pet.CC < 30 {
		return false, "成长过低无法进行抽取！"
	}
	var cqRate int = 0
	var wxProtect = false
	propNote := []string{}
	if apid != 0 {
		aprop := ps.OptSvc.PropSrv.GetProp(userId, apid, false)
		if aprop == nil || aprop.Sums == 0 {
			return false, "添加物一数量不足！"
		}
		aprop.GetM()
		if aprop.Pid == 3383 {
			if pet.MModel.Wx == 6 {
				return false, "非五系宠无法使用成长保护石！"
			}
			wxProtect = true
		} else {
			if pet.MModel.Wx < 6 {
				return false, "五系宠物无法使用抽取成长比例道具！"
			}
			effectItems := strings.Split(aprop.MModel.Effect, ":")
			if len(effectItems) > 1 && effectItems[0] == "inczhl" {
				cqRate += com.StrTo(effectItems[1]).MustInt()
			}
		}
		propNote = append(propNote, aprop.MModel.Name)
	}
	if bpid != 0 {
		bprop := ps.OptSvc.PropSrv.GetProp(userId, bpid, false)
		if bprop == nil || bprop.Sums == 0 {
			return false, "添加物二数量不足！"
		}
		bprop.GetM()
		if bprop.Pid == 3383 {
			if pet.MModel.Wx == 6 {
				return false, "非五系宠无法使用成长保护石！"
			}
			if wxProtect {
				return false, "不可同时添加两个成保护石！"
			}
			wxProtect = true
		} else {
			if pet.MModel.Wx < 6 {
				return false, "五系宠物无法使用抽取成长比例道具！"
			}
			effectItems := strings.Split(bprop.MModel.Effect, ":")
			if len(effectItems) > 1 && effectItems[0] == "inczhl" {
				cqRate += com.StrTo(effectItems[1]).MustInt()
			}
		}
		propNote = append(propNote, bprop.MModel.Name)
	}

	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()

	if apid != 0 && !ps.OptSvc.PropSrv.DecrProp(userId, apid, 1) {
		return false, "添加物一数量不足！"
	}
	if bpid != 0 && !ps.OptSvc.PropSrv.DecrProp(userId, bpid, 1) {
		return false, "添加物二数量不足！"
	}

	if pet.MModel.Wx < 6 {
		cqRate = utils.RandInt(5, 15)
	} else if pet.CC < 65 {
		if pet.Bid == 156 {
			cqRate += utils.RandInt(8, 12)
		} else {
			cqRate += utils.RandInt(10, 20)
		}
	} else if pet.CC < 85 {
		cqRate += utils.RandInt(30, 50)
	} else if pet.CC < 100 {
		cqRate += utils.RandInt(50, 65)
	} else if pet.CC < 110 {
		cqRate += 65
	} else if pet.CC < 115 {
		cqRate += 70
	} else if pet.CC < 120 {
		cqRate += 75
	} else {
		cqRate += 80
	}
	cqCzl := int(pet.CC * float64(cqRate) / 100)
	if cqCzl > 600 {
		cqCzl = 600
	}
	cqJb := int(pet.CC * 10000)
	if pet.CC > 600 {
		cqJb = 600 * 10000
	}
	if !ps.OptSvc.UserSrv.DecreaseJb(userId, cqJb) {
		return false, "金币不足！"
	}
	ps.GetDb().Model(&models.UserInfo{Uid: userId}).Update(gin.H{"czl_ss": gorm.Expr("czl_ss+?", cqCzl)})
	ps.GetDb().Model(pet).Update(gin.H{"czl": "1", "cqflag": 1})
	SelfGameLog(userId, fmt.Sprintf("抽取成长：被抽取的宠物 %s (id=%d 成长:%s),抽取了:%d成长,使用物品:%s",
		pet.MModel.Name, pet.ID, pet.Czl, cqCzl, strings.Join(propNote, "，")), 103)
	ps.OptSvc.Commit()
	return true, "成功抽取：" + strconv.Itoa(cqCzl)
}

func (ps *PetService) Zhuanhua(userId, petId, czl int) (bool, string) {
	pet := ps.GetPet(userId, petId)
	if pet == nil {
		return false, "宠物不存在！"
	}
	pet.GetM()
	if pet.MModel.Wx < 7 {
		return false, "非神圣宠物不可转化成长！"
	}
	ssRule := common2.GetSSJhRule(pet.MModel.ID)
	if ssRule == nil {
		return false, "该宠物不可转化成长！"
	}
	if ssRule.MaxCzl <= int(pet.CC) {
		return false, "该宠物成长已达最大值！"
	}
	userInfo := ps.OptSvc.UserSrv.GetUserInfoById(userId)
	if userInfo.CzlSS < czl || ps.GetDb().Model(userInfo).Where("czl_ss>=?", czl).Update(gin.H{"czl_ss": gorm.Expr("czl_ss-?", czl)}).RowsAffected == 0 {
		return false, "可用成长不足！"
	}
	newCzl := pet.CC + float64(czl)
	if ssRule.MaxCzl <= int(newCzl) {
		ps.GetDb().Model(pet).Update(gin.H{"czl": strconv.Itoa(ssRule.MaxCzl)})
	} else {
		ps.GetDb().Model(pet).Update(gin.H{"czl": utils.CzlStr(newCzl)})
	}
	SelfGameLog(userId, fmt.Sprintf("神圣成长转化：转化%d成长给%s", czl, pet.MModel.Name), 103)
	return true, "转化成长成功！"
}

func (ps *PetService) SSEvolution(userId, petId, fzid int) (bool, string) {
	pet := ps.GetPet(userId, petId)
	if pet == nil {
		return false, "宠物不存在！"
	}
	if pet.ReMakeTimes >= 10 {
		return false, "宠物进化次数达上限！"
	}

	pet.GetM()
	if pet.MModel.Wx != 7 {
		return false, "该宠物非神圣宠物！"
	}
	ssjhRule := common2.GetSSJhRule(pet.MModel.ID)
	if ssjhRule == nil {
		return false, "该宠物无法进化！"
	}
	logNote := "神圣进化："

	levelItems := strings.Split(ssjhRule.NeedLevels, ",")
	propItems := strings.Split(ssjhRule.NeedProps, ",")
	if len(levelItems) <= pet.ReMakeTimes || len(propItems) <= pet.ReMakeTimes {
		return false, ""
	}
	items := strings.Split(propItems[pet.ReMakeTimes], ":")
	propId := com.StrTo(items[0]).MustInt()
	num := com.StrTo(items[1]).MustInt()
	level := com.StrTo(levelItems[pet.ReMakeTimes]).MustInt()
	jb := (ssjhRule.ZsProgress + pet.ReMakeTimes) * 10000
	if pet.Level < level {
		return false, "宠物等级不足！"
	}
	mprop := common2.GetMProp(propId)
	logNote += fmt.Sprintf("消耗物品 %s*%d, ", mprop.Name, num)
	addCzl := 0.0
	czlItems := strings.Split(strings.ReplaceAll(mprop.Effect, "ssjh:", ""), ":")
	if len(czlItems) > 1 {
		start := com.StrTo(czlItems[0]).MustFloat64()
		end := com.StrTo(czlItems[1]).MustFloat64()
		addCzl = rand.Float64()*(end+0.1-start) + start
	}

	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	if !ps.OptSvc.UserSrv.DecreaseJb(userId, jb) {
		return false, "金币不足！"
	}
	if !ps.OptSvc.PropSrv.DecrPropByPid(userId, propId, num) {
		return false, "进化道具不足！"
	}
	attrType := ""
	attrNum := 0
	if fzid != 0 {
		fzprop := ps.OptSvc.PropSrv.GetProp(userId, fzid, false)
		if fzprop == nil || fzprop.Sums == 0 || !ps.OptSvc.PropSrv.DecrProp(userId, fzprop.ID, 1) {
			return false, "所添加辅助道具数量不足！"
		}
		fzprop.GetM()
		if strings.Index(fzprop.MModel.Effect, "zjsxdj_") == -1 {
			return false, "所添加辅助道具种类不对！"
		}
		effectItems := strings.Split(fzprop.MModel.Effect, ":")
		if len(effectItems) < 2 {
			return false, "所添加辅助道具种类不对！"
		}
		attrType = strings.ReplaceAll(effectItems[0], "zjsxdj_", "")
		attrNum = com.StrTo(effectItems[1]).MustInt()
		logNote += fmt.Sprintf("属性加成物品：%s, ", fzprop.MModel.Name)
	}

	growth := common2.GetGrowth(pet.MModel.Wx)
	newHp := int(float64(pet.Hp) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Hp))))
	newMp := int(float64(pet.Mp) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Mp))))
	newAc := int(float64(pet.Ac) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Ac))))
	newMc := int(float64(pet.Mc) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Mc))))
	newHits := int(float64(pet.Hits) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Hits))))
	newMiss := int(float64(pet.Miss) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Miss))))
	newSpeed := int(float64(pet.Speed) * (0.3 + float64(pet.ReMakeTimes+1)/30 + float64(pet.ReMakeTimes+1)*float64(ssjhRule.ZsProgress)/(pet.CC*float64(growth.Speed))))
	newCzl := pet.CC + addCzl
	if newCzl > float64(ssjhRule.MaxCzl) {
		newCzl = float64(ssjhRule.MaxCzl)
	}
	logNote += fmt.Sprintf("成长：%s->%s, \n", pet.Czl, utils.CzlStr(newCzl))
	logNote += fmt.Sprintf("进化前属性：hp=%d, mp=%d, ac=%d, mc=%d, hits=%d, miss=%d, speed=%d\n", pet.Hp, pet.Mp, pet.Ac, pet.Mc, pet.Hits, pet.Miss, pet.Speed)
	logNote += fmt.Sprintf("进化后属性：hp=%d, mp=%d, ac=%d, mc=%d, hits=%d, miss=%d, speed=%d\n", newHp, newMp, newAc, newMc, newHits, newMiss, newSpeed)
	switch attrType {
	case "hp":
		newHp = int((1 + float64(attrNum)*0.01) * float64(newHp))
		break
	case "mp":
		newMp = int((1 + float64(attrNum)*0.01) * float64(newHp))
		break
	case "ac":
		newAc = int((1 + float64(attrNum)*0.01) * float64(newAc))
		break
	case "mc":
		newMc = int((1 + float64(attrNum)*0.01) * float64(newMc))
		break
	case "hits":
		newHits = int((1 + float64(attrNum)*0.01) * float64(newHits))
		break
	case "miss":
		newMiss = int((1 + float64(attrNum)*0.01) * float64(newMiss))
		break
	case "speed":
		newSpeed = int((1 + float64(attrNum)*0.01) * float64(newSpeed))
		break
	}
	if ps.GetDb().Model(pet).Where("level>1").Update(gin.H{
		"remaketimes": gorm.Expr("remaketimes+1"),
		"level":       1,
		"lexp":        100,
		"nowexp":      0,
		"hp":          newHp,
		"srchp":       newHp,
		"mp":          newMp,
		"srcmp":       newMp,
		"ac":          newAc,
		"mc":          newMc,
		"hits":        newHits,
		"miss":        newMiss,
		"speed":       newSpeed,
		"czl":         utils.CzlStr(newCzl),
	}).RowsAffected > 0 {
		ps.OptSvc.Commit()
		ps.OptSvc.FightSrv.DelZbAttr(pet.ID)
		SelfGameLog(userId, logNote, 103)
		return true, "进化成功！"
	} else {
		return false, "进化失败！不可进化！"
	}
}

func (ps *PetService) SSZhuanshengInfo(userId, petId int) []gin.H {
	zsData := []gin.H{}
	pet := ps.GetPet(userId, petId)
	if pet == nil {
		return zsData
	}
	pet.GetM()
	sszsRule := []*models.SSzsRule{}
	ps.GetDb().Where("cur_pet_id=?", pet.MModel.ID).Find(&sszsRule)
	for _, rule := range sszsRule {
		topet := common2.GetMpet(rule.NextPetId)
		if topet == nil {
			continue
		}
		needProps := []string{}
		for _, str := range strings.Split(rule.NeedProps, ",") {
			items := strings.Split(str, "|")
			if len(items) > 1 {
				_prop := common2.GetMProp(com.StrTo(items[0]).MustInt())
				needProps = append(needProps, fmt.Sprintf("%s x %s", _prop.Name, items[1]))
			}
		}
		jhRule := common2.GetSSJhRule(topet.ID)
		needJb := jhRule.ZsProgress * 100000
		zsData = append(zsData, gin.H{
			"id":         topet.ID,
			"name":       topet.Name,
			"img":        topet.ImgCard,
			"need_level": rule.NeedLevel,
			"need_czl":   rule.NeedCzl,
			"need_jb":    needJb,
			"need_props": strings.Join(needProps, ","),
		})
	}
	return zsData
}

func (ps *PetService) SSZhuanSheng(userId, petId, topetId, apropId, bpropId int) (bool, string) {
	pet := ps.GetPet(userId, petId)
	if pet == nil {
		return false, "宠物不存在！"
	}
	pet.GetM()
	if pet.MModel.Wx != 7 {
		return false, "非神圣宠物不可在此转生！"
	}
	sszsRule := &models.SSzsRule{}
	ps.GetDb().Where("cur_pet_id=? and next_pet_id=?", pet.MModel.ID, topetId).First(sszsRule)
	if sszsRule.ID == 0 {
		return false, "所选转生路径不存在！"
	}
	if pet.CC < float64(sszsRule.NeedCzl) {
		return false, "成长需求不足！"
	}
	if pet.Level < sszsRule.NeedLevel {
		return false, "等级需求不足！"
	}

	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()

	jhRule := common2.GetSSJhRule(topetId)
	if jhRule == nil {
		return false, "该神圣宠物无法转生！"
	}
	needJb := jhRule.ZsProgress * 100000
	if !ps.OptSvc.UserSrv.DecreaseJb(userId, needJb) {
		return false, "金币需求不足！"
	}
	needProps := []string{}
	for _, str := range strings.Split(sszsRule.NeedProps, ",") {
		items := strings.Split(str, "|")
		if len(items) > 1 {
			pid := com.StrTo(items[0]).MustInt()
			num := com.StrTo(items[1]).MustInt()
			_prop := common2.GetMProp(pid)
			needProps = append(needProps, fmt.Sprintf("%s x %d", _prop.Name, num))
			if !ps.OptSvc.PropSrv.DecrPropByPid(userId, pid, num) {
				return false, "转生所需道具不足！"
			}
		}
	}
	czlSaveRate := 10
	addSuccessRate := 0
	attrType := ""
	attrNum := 0.0

	addProps := []string{}
	if apropId != 0 {
		prop := ps.OptSvc.PropSrv.GetProp(userId, apropId, false)
		if prop == nil || prop.Sums == 0 || !ps.OptSvc.PropSrv.DecrProp(userId, apropId, 1) {
			return false, "添加道具一数量不足！"
		}
		prop.GetM()
		if prop.MModel.VaryName != 23 {
			return false, "添加道具一不可在神圣转生中使用！"
		}
		effectItems := strings.Split(prop.MModel.Effect, ":")
		if len(effectItems) < 2 {
			return false, "添加道具一不可在神圣转生中使用！"
		}
		switch effectItems[0] {
		case "sszs":
			addSuccessRate += com.StrTo(effectItems[1]).MustInt()
			break
		case "sszsczlbh":
			czlSaveRate += com.StrTo(effectItems[1]).MustInt()
			break
		default:
			if strings.Index(effectItems[0], "add") > -1 {
				attrType = strings.ReplaceAll(effectItems[0], "add", "")
				attrNum = com.StrTo(effectItems[1]).MustFloat64()
			}
		}
		addProps = append(addProps, "，添加物一："+prop.MModel.Name)
	}
	if bpropId != 0 {
		prop := ps.OptSvc.PropSrv.GetProp(userId, bpropId, false)
		if prop == nil || prop.Sums == 0 || !ps.OptSvc.PropSrv.DecrProp(userId, bpropId, 1) {
			return false, "添加道具二数量不足！"
		}
		prop.GetM()
		if prop.MModel.VaryName != 23 {
			return false, "添加道具二不可在神圣转生中使用！"
		}
		effectItems := strings.Split(prop.MModel.Effect, ":")
		if len(effectItems) < 2 {
			return false, "添加道具二不可在神圣转生中使用！"
		}
		switch effectItems[0] {
		case "sszs":
			addSuccessRate += com.StrTo(effectItems[1]).MustInt()
			break
		case "sszsczlbh":
			czlSaveRate += com.StrTo(effectItems[1]).MustInt()
			break
		default:
			if strings.Index(effectItems[0], "add") > -1 {
				if attrType != "" {
					return false, "不可同时添加两个属性增幅道具！"
				}
				attrType = strings.ReplaceAll(effectItems[0], "add", "")
				attrNum = com.StrTo(effectItems[1]).MustFloat64()
			}
		}
		addProps = append(addProps, "，添加物二："+prop.MModel.Name)
	}
	successNum := int(float64(pet.Level) / 30 * float64(1+addSuccessRate) * 100)
	if utils.RandInt(1, 10000) <= successNum {
		// 成功
		topet := ps.GetMpet(topetId)
		newHp := int(float64(topet.Hp*jhRule.ZsProgress) + float64(pet.Hp*pet.Level)/6000 + float64(pet.Hp)*pet.CC/9000)
		newMp := int(float64(topet.Mp*jhRule.ZsProgress) + float64(pet.Mp*pet.Level)/6000 + float64(pet.Mp)*pet.CC/9000)
		newAc := int(float64(topet.Ac*jhRule.ZsProgress) + float64(pet.Ac*pet.Level)/6000 + float64(pet.Ac)*pet.CC/9000)
		newMc := int(float64(topet.Mc*jhRule.ZsProgress) + float64(pet.Mc*pet.Level)/6000 + float64(pet.Mc)*pet.CC/9000)
		newHits := int(float64(topet.Hits*jhRule.ZsProgress) + float64(pet.Hits*pet.Level)/6000 + float64(pet.Hits)*pet.CC/9000)
		newMiss := int(float64(topet.Miss*jhRule.ZsProgress) + float64(pet.Miss*pet.Level)/6000 + float64(pet.Miss)*pet.CC/9000)
		newSpeed := int(float64(topet.Speed*jhRule.ZsProgress) + float64(pet.Speed*pet.Level)/6000 + float64(pet.Speed)*pet.CC/9000)
		newCzl := pet.CC * float64(czlSaveRate) * 0.01
		switch attrType {
		case "hp":
			newHp = int((1 + attrNum) * float64(newHp))
			break
		case "mp":
			newMp = int((1 + attrNum) * float64(newHp))
			break
		case "ac":
			newAc = int((1 + attrNum) * float64(newAc))
			break
		case "mc":
			newMc = int((1 + attrNum) * float64(newMc))
			break
		case "hits":
			newHits = int((1 + attrNum) * float64(newHits))
			break
		case "miss":
			newMiss = int((1 + attrNum) * float64(newMiss))
			break
		case "speed":
			newSpeed = int((1 + attrNum) * float64(newSpeed))
			break
		}
		newPet := &models.UPet{
			Bid:         topet.ID,
			Name:        topet.Name,
			Uid:         userId,
			Level:       1,
			Wx:          topet.Wx,
			Ac:          newAc,
			Mc:          newMc,
			SrcHp:       newHp,
			Hp:          newHp,
			Mp:          newMp,
			SrcMp:       newMp,
			SkillList:   topet.SkillList,
			Stime:       utils.NowUnix(),
			NowExp:      0,
			LExp:        100,
			Hits:        newHits,
			Miss:        newMiss,
			Speed:       newSpeed,
			Kx:          topet.Kx,
			Fatting:     0,
			Czl:         utils.CzlStr(newCzl),
			Muchang:     0,
			ReMakeTimes: 0,
		}
		if ps.CreatPet(newPet) {
			ps.DropPet(pet)
			newPet.GetM()
			AnnounceAll(ps.OptSvc.UserSrv.GetUserById(userId).Nickname, "获得神圣宠物 "+newPet.MModel.Name)
			ps.OptSvc.UserSrv.SetMBid(userId, newPet.ID)
			ps.OptSvc.Commit()

			SelfGameLog(userId, fmt.Sprintf("神圣转生：成功！获得 %s, 消耗物品：%s, 辅助道具：%s, \n原宠物：%s, czl: %s, ac:%d, hits:%d\n, 新宠物：%s, czl: %s, ac:%d, hits:%d\n",
				newPet.MModel.Name, strings.Join(needProps, ","), strings.Join(addProps, ","), pet.MModel.Name, pet.Czl, pet.Ac, pet.Hits, newPet.MModel.Name, newPet.Czl, newPet.Ac, newPet.Hits), 104)
			return true, "转生成功！"
		} else {
			return false, "转生错误！"
		}
	} else {
		ps.OptSvc.Commit()
		SelfGameLog(userId, fmt.Sprintf("神圣转生：失败！消耗物品：%s, 辅助道具：%s",
			strings.Join(needProps, ","), strings.Join(addProps, ",")), 104)
		return true, "转生失败！"
	}

}
