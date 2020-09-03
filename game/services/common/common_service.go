package common

import (
	"log"
	"pokemon/game/models"
	"pokemon/game/services/datastore"
	"pokemon/game/utils"
	"sync"
)

func InitDataStore() {
	now := utils.NowUnix()
	var waitGroup = &sync.WaitGroup{}
	funcs := []func(i int){
		datastore.Mpet.Update,
		datastore.Mprop.Update,
		datastore.Mskill.Update,
		datastore.Welcome.Update,
		datastore.ExpList.Update,
		datastore.Growth.Update,
		datastore.SSjh.Update,
		datastore.SSzs.Update,
		datastore.TimeConfig.Update,
		datastore.Mmap.Update,
		datastore.Task.Update,
		datastore.CardTitle.Update,
		datastore.CardPrize.Update,
		datastore.CardSerie.Update,
	}
	funcs1 := []func(){
		datastore.InitYiwangCard,
	}
	waitGroup.Add(len(funcs))
	for _, f := range funcs {
		fff := f
		go func() {
			fff(now)
			waitGroup.Done()
		}()
	}
	waitGroup.Add(len(funcs1))
	for _, f := range funcs1 {
		go func() {
			f()
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	log.Printf("data load over!")
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

func GetMPropIdsByName(names []string) []int {
	return datastore.Mprop.GetIdsByName(names)
}

func GetMPropIdByName(name string) int {
	return datastore.Mprop.GetIdByName(name)
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

func GetCardTitleById(id int) *models.CardTitle {
	return datastore.CardTitle.GetById(id)
}

func GetAllCardTitle() []*models.CardTitle {
	return datastore.CardTitle.GetAll()
}

func GetMMap(mapId int) *models.Map {
	return datastore.Mmap.GetMap(mapId)
}

func GetGpc(gpcId int) *models.Gpc {
	return datastore.Mmap.GetGpc(gpcId)
}

func GetGpcGroup(id int) *models.GpcGroup {
	return datastore.Mmap.GetGroup(id)
}

func GetGpcGroupByLevel(level int) []int {
	return datastore.Mmap.GetGroupIds(level)
}

//func GetFbSetting(mapId int) *datastore.Fuben {
//	return datastore.GetFbSetting(mapId)
//}

func GetTimeConfig(title string) *models.TimeConfig {
	return datastore.TimeConfig.Get(title)
}

func GetTask(id int) *models.Task {
	return datastore.Task.Get(id)
}

func GetCardPrize(id int) *models.CardPrize {
	return datastore.CardPrize.Get(id)
}

func GetAllCardPrize() []*models.CardPrize {
	return datastore.CardPrize.GetAll()
}

func GetCardSeries(id int) *models.CardSeries {
	return datastore.CardSerie.Get(id)
}

func GetAllCardSeries() []*models.CardSeries {
	return datastore.CardSerie.GetAll()
}

func GetYiwangCard(id int) *models.TarotCard {
	return datastore.GetYiwangCard(id)
}

func GetYiwangCards(multiple int, needSj bool) []*models.TarotCard {
	return datastore.GetYiwangCards(multiple, needSj)
}

func GetYiwangBossCards(multiple int) []*models.TarotCard {
	return datastore.GetYiwangBossCards(multiple)
}
