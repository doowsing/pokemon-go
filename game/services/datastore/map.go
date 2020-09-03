package datastore

import (
	"pokemon/game/models"
	"strconv"
	"sync"
)

var Mmap = &mmapStore{store: newStore(), dataMap: make([]*models.Map, 500), dataGpc: make([]*models.Gpc, 10000)}

type mmapStore struct {
	*store
	dataMap    []*models.Map
	_dataMap   []*models.Map
	dataGpc    []*models.Gpc
	_dataGpc   []*models.Gpc
	cGroupData sync.Map
}

func (s *mmapStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s.CopyBackup()
	s.updating = true

	s.dataMap = make([]*models.Map, 500)
	s.dataGpc = make([]*models.Gpc, 10000)
	var allMap []*models.Map
	var allGpc []*models.Gpc
	GetDb().Find(&allMap)
	GetDb().Find(&allGpc)
	lv2gpcIds := make(map[int][]int)
	lv2gpcIdsNoBoss := make(map[int][]int)
	for _, v := range allGpc {
		s.dataGpc[v.ID] = v
		if v.Boss < 4 {
			lv2gpcIds[v.Level] = append(lv2gpcIds[v.Level], v.ID)
			if v.Boss < 3 {
				lv2gpcIdsNoBoss[v.Level] = append(lv2gpcIdsNoBoss[v.Level], v.ID)
			}
		}
	}
	for _, v := range allMap {
		if len(v.Levels) > 1 {
			v.Lv2GpcIds = make(map[int][]int)
			for i := v.Levels[0]; i <= v.Levels[1]; i++ {
				if ids := lv2gpcIds[i]; len(ids) > 0 {
					v.Lv2GpcIds[i] = ids
					if ids = lv2gpcIdsNoBoss[i]; len(ids) > 0 {
						v.GpcIds = append(v.GpcIds, ids...)
					}
				}
			}
		}
		s.dataMap[v.ID] = v
	}
	//fmt.Printf("lv2gpcIds:%v\n", lv2gpcIds)

	var allGroup []*models.GpcGroup
	GetDb().Find(&allGroup)
	lv2GroupId := make([][]int, 1000)
	for _, g := range allGroup {
		s.cGroupData.Store(g.ID, g)
		lv2GroupId[g.Level] = append(lv2GroupId[g.Level], g.ID)
	}
	for lv, ids := range lv2GroupId {
		if len(ids) > 0 {
			s.cGroupData.Store("level"+strconv.Itoa(lv), ids)
		}
	}
	s.updating = false
}

func (s *mmapStore) CopyBackup() {

	s._dataMap = make([]*models.Map, 500)
	s._dataGpc = make([]*models.Gpc, 10000)
	copy(s._dataGpc, s.dataGpc)
	for _, v := range s.dataMap {
		if v == nil {
			continue
		}
		if len(v.Levels) > 1 {
			for _, gpc := range s._dataGpc {
				if gpc == nil {
					continue
				}
				if gpc.Boss < 4 && (gpc.Level >= v.Levels[0] || gpc.Level <= v.Levels[1]) {
					v.GpcIds = append(v.GpcIds, gpc.ID)
				}
			}
		}
		s._dataMap[v.ID] = v
	}
}

func (s *mmapStore) getDataMap() []*models.Map {
	if s.updating {
		return s._dataMap
	}
	return s.dataMap
}

func (s *mmapStore) getDataGpc() []*models.Gpc {
	if s.updating {
		return s._dataGpc
	}
	return s.dataGpc
}

func (s *mmapStore) GetMap(id int) *models.Map {
	if id < len(s.getDataMap()) {
		return s.getDataMap()[id]
	}
	return nil
}

func (s *mmapStore) GetGroup(id int) *models.GpcGroup {
	if data, ok := s.cGroupData.Load(id); ok {
		return data.(*models.GpcGroup)
	}
	return nil
}

func (s *mmapStore) GetGroupIds(level int) []int {
	if data, ok := s.cGroupData.Load("level" + strconv.Itoa(level)); ok {
		return data.([]int)
	}
	return nil
}

func (s *mmapStore) GetGpc(id int) *models.Gpc {
	if id < len(s.getDataGpc()) {
		return s.getDataGpc()[id]
	}
	return nil
}

var allYiwangCards map[int]*models.TarotCard
var yiwangCards map[int][]*models.TarotCard
var yiwangSjCards map[int][]*models.TarotCard
var yiwangBossCards map[int][]*models.TarotCard

func InitYiwangCard() {
	cards := []*models.TarotCard{}
	GetDb().Find(&cards)

	allYiwangCards = make(map[int]*models.TarotCard)
	yiwangCards = make(map[int][]*models.TarotCard)
	yiwangSjCards = make(map[int][]*models.TarotCard)
	yiwangBossCards = make(map[int][]*models.TarotCard)
	for _, card := range cards {
		allYiwangCards[card.Id] = card
		if card.Flag == 0 {
			if card.Sj > 0 {
				yiwangSjCards[card.Multiple] = append(yiwangSjCards[card.Multiple], card)
			} else {
				yiwangCards[card.Multiple] = append(yiwangCards[card.Multiple], card)
			}

		} else {
			yiwangBossCards[card.Multiple] = append(yiwangBossCards[card.Multiple], card)
		}
	}
}

func GetYiwangCards(multiple int, needSj bool) []*models.TarotCard {
	cards := []*models.TarotCard{}
	if needSj {
		copy(cards, yiwangSjCards[multiple])
	} else {
		copy(cards, yiwangCards[multiple])
	}

	return cards
}

func GetYiwangCard(id int) *models.TarotCard {
	return allYiwangCards[id]
}

func GetYiwangBossCards(multiple int) []*models.TarotCard {
	cards := []*models.TarotCard{}
	copy(cards, yiwangBossCards[multiple])
	return cards
}
