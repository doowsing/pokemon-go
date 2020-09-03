package services

import (
	"pokemon/pkg/models"
	"pokemon/pkg/persistence"
	"pokemon/pkg/services/datastore"
	"pokemon/pkg/utils"
)

var db = persistence.GetOrm()

func InitDataStore() {
	now := utils.NowUnix()
	go datastore.Mpet.Update(now)
	go datastore.Mprop.Update(now)
	go datastore.Mskill.Update(now)
	go datastore.Welcome.Update(now)
	go datastore.ExpList.Update(now)
	go datastore.Growth.Update(now)
	go datastore.SSjh.Update(now)
	go datastore.SSzs.Update(now)
	go datastore.TimeConfig.Update(now)
	go datastore.InitMap(now)
	go datastore.Task.Update(now)
}

func GetMpet(petId int) *models.MPet {
	return datastore.Mpet.Get(petId)
}

func GetMpetByName(Name string) *models.MPet {
	return datastore.Mpet.GetByName(Name)
}

func GetMskill(skillId int) *models.MSkill {
	return datastore.Mskill.Get(skillId)
}

func GetMskillByPid(pid int) *models.MSkill {
	return datastore.Mskill.GetByPid(pid)
}

func GetMProp(MPropId int) *models.MProp {
	return datastore.Mprop.Get(MPropId)
}

func GetWelcome(code string) *models.Welcome {
	return datastore.Welcome.Get(code)
}

func GetNextExp(nowLv int) int {
	return datastore.ExpList.Get(nowLv)
}

func GetGrowth(wx int) *models.Growth {
	return datastore.Growth.Get(wx)
}

func GetSSJhRule(petid int) *models.SSjhRule {
	return datastore.SSjh.Get(petid)
}

func GetSSZsRule(petid int) *models.SSzsRule {
	return datastore.SSzs.Get(petid)
}

func GetCardTitle(codeName string) *models.CardTitle {
	return datastore.CardTitle.Get(codeName)
}

func GetMMap(mapId int) *models.Map {
	return datastore.Mmap.Get(mapId)
}

func GetGpc(gpcId int) *models.Gpc {
	return datastore.Mgpc.Get(gpcId)
}
func GetFbSetting(mapId int) *datastore.Fuben {
	return datastore.GetFbSetting(mapId)
}

func GetTimeConfig(title string) *models.TimeConfig {
	return datastore.TimeConfig.Get(title)
}

func GetTask(id int) *models.Task {
	return datastore.Task.Get(id)
}
